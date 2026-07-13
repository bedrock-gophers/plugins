use super::*;

#[derive(Default)]
struct FalliblePlugin;

impl Plugin for FalliblePlugin {
    fn on_enable(&self) -> PluginResult {
        Err(std::io::Error::other("configuration rejected").into())
    }
}

#[test]
fn plugin_enable_uses_ordinary_rust_errors() {
    assert_eq!(
        FalliblePlugin.on_enable().unwrap_err().to_string(),
        "configuration rejected"
    );
}

#[test]
fn lifecycle_error_truncation_preserves_utf8() {
    let mut bytes = [0u8; 3];
    let mut buffer = dragonfly_plugin_sys::DfStringBuffer {
        data: bytes.as_mut_ptr(),
        len: u64::MAX,
        capacity: bytes.len() as u64,
    };
    unsafe { __write_plugin_error(&mut buffer, "éé") };
    assert_eq!(buffer.len, 2);
    assert_eq!(&bytes[..2], "é".as_bytes());
}
