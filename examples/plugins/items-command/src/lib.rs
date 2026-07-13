use dragonfly::{Context, Enchantment, Player, Plugin, Value, item, plugin};

#[derive(Default)]
struct ItemsCommand;

#[plugin]
#[command("items")]
impl Plugin for ItemsCommand {
    #[command]
    fn root(&self, context: &mut Context<'_, Player>) {
        context
            .source()
            .message("Use /items give-sword or /items copy <from> <to>.");
    }

    #[subcommand("give-sword")]
    fn give_sword(&self, context: &mut Context<'_, Player>) {
        let sword = item::new(item::Sword::new(item::ToolTier::Diamond), 1)
            .with_custom_name("Rust Sword")
            .with_lore(["Created by a native plugin", "Still a Dragonfly item"])
            .with_value("plugin", Value::from("items-command"))
            .with_enchantment(Enchantment::Sharpness, 5);
        let added = context.source().inventory().add_item(&sword);
        context.source().message(&format!("Added {added} item."));
    }

    #[subcommand("copy")]
    fn copy(&self, context: &mut Context<'_, Player>, from: i64, to: i64) {
        let (Ok(from), Ok(to)) = (usize::try_from(from), usize::try_from(to)) else {
            context.source().message("Slots must be non-negative.");
            return;
        };
        let inventory = context.source().inventory();
        let Some(stack) = inventory.item(from) else {
            context.source().message("Could not read that slot.");
            return;
        };
        inventory.set_item(to, &stack);
        context.source().message("Item copied.");
    }
}
