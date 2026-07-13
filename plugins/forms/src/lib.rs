use dragonfly::{Event, Plugin, form, plugin};

#[derive(Default)]
struct FormsPlugin;

#[plugin]
impl Plugin for FormsPlugin {
    fn on_join(&self, event: &mut Event::PlayerJoin<'_>) {
        let mut menu = form::Menu::new("Rust forms").body("Choose a Dragonfly form type.");
        menu.header("Form examples");
        menu.label("Responses are handled by owned Rust callbacks.");
        menu.divider();
        let custom = menu.button(form::Button::new("Custom form"));
        let modal =
            menu.button(form::Button::new("Modal form").image("textures/ui/icon_recipe_item"));
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
    let modal = form::Modal::yes_no("Modal form").body("Would you like to open the custom form?");
    player.send_form(modal, |player, response| {
        if response.is_some_and(|response| response.accepted()) {
            send_custom(player);
        }
    });
}

fn send_custom(player: dragonfly::Player) {
    let mut custom = form::Custom::new("Custom form");
    custom.header("Profile");
    custom.label("Every Dragonfly v0.11 custom element is available.");
    custom.divider();
    let name = custom.input(
        form::Input::new("Name")
            .placeholder("Steve")
            .tooltip("Your display name"),
    );
    let enabled = custom.toggle(
        form::Toggle::new("Enable greetings", true).tooltip("Controls the response message"),
    );
    let volume = custom
        .slider(form::Slider::new("Volume", 0.0, 10.0, 1.0, 5.0).tooltip("Choose from 0 to 10"));
    let colour = custom.dropdown(form::Dropdown::new("Colour", ["Red", "Green", "Blue"], 1));
    let speed = custom.step_slider(form::StepSlider::new(
        "Speed",
        ["Slow", "Normal", "Fast"],
        1,
    ));
    player.send_form(custom, move |player, response| {
        let Some(response) = response else { return };
        let Some(name) = response.get(name) else {
            return;
        };
        let Some(enabled) = response.get(enabled) else {
            return;
        };
        let Some(volume) = response.get(volume) else {
            return;
        };
        let Some(colour) = response.get(colour) else {
            return;
        };
        let Some(speed) = response.get(speed) else {
            return;
        };
        if enabled {
            player.message(&format!(
                "Hello {name}: volume {volume}, colour #{colour}, speed #{speed}"
            ));
        }
    });
}
