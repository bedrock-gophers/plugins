use super::*;
use dragonfly_plugin_sys::{DfAbiHeader, DfPluginDisableFn, DfPluginEnableFn};
use std::sync::mpsc;
use std::thread;
use std::time::Duration;

static CALLS: Mutex<Vec<&'static str>> = Mutex::new(Vec::new());
static TEST_LOCK: Mutex<()> = Mutex::new(());

fn record(value: &'static str) {
    CALLS.lock().unwrap().push(value);
}

unsafe extern "C" fn enable_first(_instance: *mut c_void, _error: *mut DfStringBuffer) -> DfStatus {
    record("first enable");
    DF_STATUS_OK
}

unsafe extern "C" fn disable_first(_instance: *mut c_void) -> DfStatus {
    record("first disable");
    DF_STATUS_OK
}

unsafe extern "C" fn enable_second(
    _instance: *mut c_void,
    _error: *mut DfStringBuffer,
) -> DfStatus {
    record("second enable");
    DF_STATUS_OK
}

unsafe extern "C" fn disable_second(_instance: *mut c_void) -> DfStatus {
    record("second disable");
    DF_STATUS_OK
}

unsafe extern "C" fn enable_failure(
    _instance: *mut c_void,
    error: *mut DfStringBuffer,
) -> DfStatus {
    record("failure enable");
    let error = unsafe { &mut *error };
    let message = b"invalid practice configuration";
    unsafe { ptr::copy_nonoverlapping(message.as_ptr(), error.data, message.len()) };
    error.len = message.len() as u64;
    DF_STATUS_ERROR
}

unsafe extern "C" fn disable_failure(_instance: *mut c_void) -> DfStatus {
    record("failure disable");
    DF_STATUS_OK
}

const fn plugin_api(
    id: &'static [u8],
    enable: DfPluginEnableFn,
    disable: DfPluginDisableFn,
) -> DfPluginApiV4 {
    DfPluginApiV4 {
        header: DfAbiHeader {
            abi_version: DF_ABI_VERSION,
            struct_size: size_of::<DfPluginApiV4>() as u32,
            subscriptions: 0,
        },
        plugin_id: DfStringView {
            data: id.as_ptr(),
            len: id.len() as u64,
        },
        create: None,
        enable: Some(enable),
        disable: Some(disable),
        commands: None,
        entity_type_count: None,
        entity_type_at: None,
        handle_entity: None,
        handle_command: None,
        command_enum_options: None,
        set_host: None,
        destroy: None,
        handle_event: None,
    }
}

static FIRST_API: DfPluginApiV4 = plugin_api(b"first", enable_first, disable_first);
static SECOND_API: DfPluginApiV4 = plugin_api(b"second", enable_second, disable_second);
static FAILURE_API: DfPluginApiV4 = plugin_api(b"failure", enable_failure, disable_failure);

#[cfg(unix)]
fn current_library() -> Library {
    libloading::os::unix::Library::this().into()
}

#[cfg(windows)]
fn current_library() -> Library {
    libloading::os::windows::Library::this().unwrap().into()
}

fn plugin(api: &'static DfPluginApiV4, id: &str) -> LoadedPlugin {
    LoadedPlugin {
        api,
        instance: ptr::null_mut(),
        id: id.to_owned(),
        enabled: AtomicBool::new(false),
        _library: current_library(),
    }
}

fn runtime(plugins: Vec<LoadedPlugin>) -> DfRuntime {
    DfRuntime {
        plugins,
        entity_types: Vec::new(),
        entity_instances: RwLock::new(HashMap::new()),
        next_entity_instance: AtomicU64::new(1),
        commands: RwLock::new(Vec::new()),
        subscriptions: 0,
        lifecycle: Mutex::new(()),
        admission: Admission::new(),
    }
}

#[test]
fn failed_enable_cleans_failing_plugin_then_prior_plugins_in_reverse() {
    let _test = TEST_LOCK.lock().unwrap();
    CALLS.lock().unwrap().clear();
    let runtime = runtime(vec![
        plugin(&FIRST_API, "first"),
        plugin(&SECOND_API, "second"),
        plugin(&FAILURE_API, "failure"),
    ]);

    assert_eq!(
        runtime.enable(),
        Err("plugin \"failure\" failed to enable: invalid practice configuration".to_owned())
    );
    assert_eq!(
        CALLS.lock().unwrap().as_slice(),
        [
            "first enable",
            "second enable",
            "failure enable",
            "failure disable",
            "second disable",
            "first disable",
        ]
    );
    assert!(
        runtime
            .plugins
            .iter()
            .all(|plugin| !plugin.enabled.load(Ordering::Acquire))
    );
    assert!(runtime.admission.enter(AdmissionKind::Ordinary).is_none());
    assert!(runtime.admission.enter(AdmissionKind::Entity).is_some());

    runtime.finish_disable();
    assert_eq!(CALLS.lock().unwrap().len(), 6);
    assert!(runtime.admission.enter(AdmissionKind::Entity).is_none());
}

#[test]
fn disable_is_idempotent() {
    let _test = TEST_LOCK.lock().unwrap();
    CALLS.lock().unwrap().clear();
    let runtime = runtime(vec![plugin(&FIRST_API, "first")]);

    assert_eq!(runtime.enable(), Ok(()));
    runtime.disable();
    runtime.disable();
    assert_eq!(
        CALLS.lock().unwrap().as_slice(),
        ["first enable", "first disable"]
    );
}

#[test]
fn begin_waits_for_ordinary_calls_and_keeps_entities_admitted() {
    let admission = Arc::new(Admission::new());
    admission.start_enable().unwrap();
    admission.finish_enable();
    let ordinary = admission.enter(AdmissionKind::Ordinary).unwrap();
    let (done_tx, done_rx) = mpsc::channel();
    let worker = Arc::clone(&admission);
    let join = thread::spawn(move || {
        worker.begin_stopping();
        done_tx.send(()).unwrap();
    });

    assert!(done_rx.recv_timeout(Duration::from_millis(25)).is_err());
    drop(ordinary);
    done_rx.recv_timeout(Duration::from_secs(1)).unwrap();
    join.join().unwrap();

    assert!(admission.enter(AdmissionKind::Ordinary).is_none());
    assert!(admission.enter(AdmissionKind::Entity).is_some());
}

#[test]
fn finish_waits_for_entity_calls_then_rejects_them() {
    let admission = Arc::new(Admission::new());
    admission.start_enable().unwrap();
    admission.finish_enable();
    admission.begin_stopping();
    let entity = admission.enter(AdmissionKind::Entity).unwrap();
    let (done_tx, done_rx) = mpsc::channel();
    let worker = Arc::clone(&admission);
    let join = thread::spawn(move || {
        worker.begin_finishing();
        worker.finish();
        done_tx.send(()).unwrap();
    });

    assert!(done_rx.recv_timeout(Duration::from_millis(25)).is_err());
    drop(entity);
    done_rx.recv_timeout(Duration::from_secs(1)).unwrap();
    join.join().unwrap();

    assert!(admission.enter(AdmissionKind::Entity).is_none());
}
