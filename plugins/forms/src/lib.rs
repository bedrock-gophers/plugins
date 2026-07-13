use dragonfly::{Event, Plugin, form, plugin};

#[derive(Default)]
struct FormsPlugin;

#[plugin]
impl Plugin for FormsPlugin {
    fn on_join(&self, event: &mut Event::PlayerJoin<'_>) {
        let mut menu = form::Menu::new("Rust forms").body("Choose a form type.");
        let custom = menu.button(form::Button::new("Custom form"));
        let modal = menu.button(form::Button::new("Modal form"));
        event.player().send_form(menu, move |player, response| {
            let Some(response) = response else { return };
            if response.selected() == custom {
                send_custom(player);
            } else if response.selected() == modal {
                send_modal(player);
            }
        });
    }
}

fn send_modal(player: dragonfly::Player) {
    let modal = form::Modal::yes_no("Modal form").body("Open the custom form?");
    player.send_form(modal, |player, response| {
        if response.is_some_and(|response| response.accepted()) {
            send_custom(player);
        }
    });
}

fn send_custom(player: dragonfly::Player) {
    let mut custom = form::Custom::new("Custom form");
    let name = custom.input(form::Input::new("Name").placeholder("Steve"));
    let enabled = custom.toggle(form::Toggle::new("Enable greeting", true));
    player.send_form(custom, move |player, response| {
        let Some(response) = response else { return };
        let Some(name) = response.get(name) else {
            return;
        };
        if response.get(enabled) == Some(true) {
            player.message(&format!("Hello {name}"));
        }
    });
}
