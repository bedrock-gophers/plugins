use dragonfly::{Event, Plugin, Scoreboard, plugin};

#[derive(Default)]
struct ScoreboardPlugin;

#[plugin]
impl Plugin for ScoreboardPlugin {
    fn on_join(&self, event: &mut Event::PlayerJoin<'_>) {
        let mut board = Scoreboard::new("Rust scoreboard");
        if board
            .push_line(format!("Welcome, {}", event.name()))
            .and_then(|board| board.push_line("Native plugins"))
            .is_err()
        {
            return;
        }
        board.remove_padding().set_descending(true);
        event.player().send_scoreboard(&board);
    }

    fn on_quit(&self, event: &Event::PlayerQuit<'_>) {
        event.player().remove_scoreboard();
    }
}
