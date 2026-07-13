use dragonfly::{Event, Plugin, plugin};

#[derive(Default)]
struct ChatFilter;

#[plugin]
impl Plugin for ChatFilter {
    fn on_chat(&self, event: &mut Event::PlayerChat<'_>) {
        let filtered = event.message().replace("foo", "bar");
        event.replace_message(&filtered).expect("message too long");
    }
}
