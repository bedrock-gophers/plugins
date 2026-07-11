use dragonfly_plugin::{Plugin, plugin};

#[derive(Default)]
struct LifecycleLogger;

#[plugin]
impl Plugin for LifecycleLogger {
    fn on_enable(&self) {
        eprintln!("enabled");
    }

    fn on_disable(&self) {
        eprintln!("disabled");
    }
}

