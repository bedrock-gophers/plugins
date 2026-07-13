use dragonfly::{Event, Plugin, plugin};

#[derive(Default)]
struct MovementGuard;

#[plugin]
impl Plugin for MovementGuard {
    fn on_move(&self, event: &mut Event::PlayerMove<'_>) {
        if event.new_position().y < -64.0 {
            event.cancel();
        }
    }
}
