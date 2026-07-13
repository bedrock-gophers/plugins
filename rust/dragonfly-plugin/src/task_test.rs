use std::{
    sync::{
        Arc, Barrier,
        atomic::{AtomicBool, AtomicUsize, Ordering},
        mpsc,
    },
    thread,
    time::Duration,
};

use super::{SpawnError, TaskGroup, TaskGroupState};

const SHORT_WAIT: Duration = Duration::from_millis(25);
const TEST_TIMEOUT: Duration = Duration::from_secs(2);

#[test]
fn cancellation_token_waits_and_remains_cancelled() {
    let group = Arc::new(TaskGroup::new());
    let token = group.cancellation_token();
    assert!(!token.is_cancelled());
    assert!(!token.wait_timeout(SHORT_WAIT));

    let shutdown = Arc::clone(&group);
    let caller = thread::spawn(move || shutdown.shutdown());
    token.wait();
    caller.join().unwrap();

    assert!(token.is_cancelled());
    assert!(token.wait_timeout(Duration::ZERO));
    assert_eq!(group.state(), TaskGroupState::Closed);
}

#[test]
fn shutdown_cancels_and_joins_every_accepted_worker() {
    let group = TaskGroup::new();
    let finished = Arc::new(AtomicUsize::new(0));
    for _ in 0..8 {
        let finished = Arc::clone(&finished);
        group
            .spawn(move |cancellation| {
                cancellation.wait();
                finished.fetch_add(1, Ordering::SeqCst);
            })
            .unwrap();
    }

    group.shutdown();

    assert_eq!(finished.load(Ordering::SeqCst), 8);
    assert_eq!(group.state(), TaskGroupState::Closed);
    assert_eq!(group.spawn(|_| {}), Err(SpawnError::ShuttingDown));
}

#[test]
fn shutdown_rejects_a_concurrent_spawn_or_joins_it() {
    for _ in 0..64 {
        let group = Arc::new(TaskGroup::new());
        let start = Arc::new(Barrier::new(2));
        let completed = Arc::new(AtomicBool::new(false));

        let spawning_group = Arc::clone(&group);
        let spawning_start = Arc::clone(&start);
        let spawning_completed = Arc::clone(&completed);
        let spawn = thread::spawn(move || {
            spawning_start.wait();
            spawning_group.spawn(move |cancellation| {
                cancellation.wait();
                spawning_completed.store(true, Ordering::SeqCst);
            })
        });

        start.wait();
        group.shutdown();
        match spawn.join().unwrap() {
            Ok(()) => assert!(completed.load(Ordering::SeqCst)),
            Err(error) => assert_eq!(error, SpawnError::ShuttingDown),
        }
        assert_eq!(group.state(), TaskGroupState::Closed);
    }
}

#[test]
fn concurrent_shutdown_callers_wait_for_worker_join() {
    let group = Arc::new(TaskGroup::new());
    let (cancelled_tx, cancelled_rx) = mpsc::channel();
    let (release_tx, release_rx) = mpsc::channel();
    group
        .spawn(move |cancellation| {
            cancellation.wait();
            cancelled_tx.send(()).unwrap();
            release_rx.recv().unwrap();
        })
        .unwrap();

    let (first_done_tx, first_done_rx) = mpsc::channel();
    let first_group = Arc::clone(&group);
    let first = thread::spawn(move || {
        first_group.shutdown();
        first_done_tx.send(()).unwrap();
    });
    cancelled_rx.recv_timeout(TEST_TIMEOUT).unwrap();
    assert_eq!(group.state(), TaskGroupState::Closing);
    assert_eq!(group.spawn(|_| {}), Err(SpawnError::ShuttingDown));

    let (second_done_tx, second_done_rx) = mpsc::channel();
    let second_group = Arc::clone(&group);
    let second = thread::spawn(move || {
        second_group.shutdown();
        second_done_tx.send(()).unwrap();
    });
    assert!(first_done_rx.recv_timeout(SHORT_WAIT).is_err());
    assert!(second_done_rx.recv_timeout(SHORT_WAIT).is_err());

    release_tx.send(()).unwrap();
    first_done_rx.recv_timeout(TEST_TIMEOUT).unwrap();
    second_done_rx.recv_timeout(TEST_TIMEOUT).unwrap();
    first.join().unwrap();
    second.join().unwrap();
    assert_eq!(group.state(), TaskGroupState::Closed);
}

#[test]
fn worker_panic_is_contained() {
    let group = TaskGroup::new();
    let completed = Arc::new(AtomicBool::new(false));
    group.spawn(|_| panic!("worker panic")).unwrap();
    let completed_worker = Arc::clone(&completed);
    group
        .spawn(move |cancellation| {
            cancellation.wait();
            completed_worker.store(true, Ordering::SeqCst);
        })
        .unwrap();

    group.shutdown();

    assert!(completed.load(Ordering::SeqCst));
    assert_eq!(group.state(), TaskGroupState::Closed);
}

#[test]
fn drop_cancels_and_joins_workers() {
    let (finished_tx, finished_rx) = mpsc::channel();
    {
        let group = TaskGroup::new();
        group
            .spawn(move |cancellation| {
                cancellation.wait();
                finished_tx.send(()).unwrap();
            })
            .unwrap();
    }

    finished_rx.recv_timeout(TEST_TIMEOUT).unwrap();
}
