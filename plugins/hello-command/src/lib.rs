use dragonfly_plugin::{
    Command, CommandEvent, CommandOverload, CommandParameter, CommandValue, Plugin, plugin,
};

static STYLES: &[CommandValue] = &[
    CommandValue::new("plain"),
    CommandValue::new("excited"),
];
static COMMANDS: &[Command] = &[Command::new("hello", "Greets the command source").with_overloads(
    &[
        CommandOverload::new(&[
            CommandParameter::subcommand("say"),
            CommandParameter::enumeration("style", STYLES),
        ]),
        CommandOverload::new(&[CommandParameter::subcommand("about")]),
    ],
)];

#[derive(Default)]
struct HelloCommand;

#[plugin]
impl Plugin for HelloCommand {
    fn commands(&self) -> &'static [Command] {
        COMMANDS
    }

    fn on_command(&self, _command: usize, event: &mut CommandEvent<'_>) {
        event
            .reply(&format!(
                "Hello, {}! You passed: {}",
                event.source(),
                event.arguments()
            ))
            .expect("command reply fits the runtime output buffer");
    }
}
