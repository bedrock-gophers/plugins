//! Bedrock forms with owned, one-shot response callbacks.

use crate::{Player, bytes_view, host_api};
use core::{ffi::c_void, marker::PhantomData};
pub use serde_json::{Value, json};

mod private {
    pub trait Sealed {}
}

pub trait Form: private::Sealed + Send + 'static {
    type Response: Send + 'static;
    #[doc(hidden)]
    fn request_json(&self) -> Vec<u8>;
    #[doc(hidden)]
    fn parse_response(&self, data: &[u8]) -> Option<Self::Response>;
}

/// An explicit escape hatch for forms not represented by the typed API.
///
/// Prefer [`Menu`], [`Modal`], or [`Custom`]. `Raw` exposes the Bedrock JSON
/// request and response directly for experimental or custom form shapes.
#[derive(Clone, Debug, PartialEq)]
pub struct Raw {
    request: Value,
}

impl Raw {
    pub fn new(request: Value) -> Self {
        Self { request }
    }
}

impl private::Sealed for Raw {}

impl Form for Raw {
    type Response = Value;

    fn request_json(&self) -> Vec<u8> {
        serde_json::to_vec(&self.request).unwrap_or_default()
    }

    fn parse_response(&self, data: &[u8]) -> Option<Self::Response> {
        serde_json::from_slice(data).ok()
    }
}

#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Button {
    text: String,
    image: String,
}
impl Button {
    pub fn new(text: impl Into<String>) -> Self {
        Self {
            text: text.into(),
            image: String::new(),
        }
    }
    pub fn image(mut self, image: impl Into<String>) -> Self {
        self.image = image.into();
        self
    }
    fn json(&self) -> Value {
        let mut value = json!({"type":"button", "text":self.text});
        if !self.image.is_empty() {
            let kind = if self.image.starts_with("http:") || self.image.starts_with("https:") {
                "url"
            } else {
                "path"
            };
            value["image"] = json!({"type":kind, "data":self.image});
        }
        value
    }
}

#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Divider;
#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Header(pub String);
#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Label(pub String);
impl Header {
    pub fn new(text: impl Into<String>) -> Self {
        Self(text.into())
    }
}
impl Label {
    pub fn new(text: impl Into<String>) -> Self {
        Self(text.into())
    }
}

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct ButtonId(usize);
#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub struct MenuResponse {
    selected: ButtonId,
}
impl MenuResponse {
    pub const fn selected(&self) -> ButtonId {
        self.selected
    }
}

enum MenuElement {
    Button(Button),
    Divider,
    Header(Header),
    Label(Label),
}
pub struct Menu {
    title: String,
    body: String,
    elements: Vec<MenuElement>,
    buttons: usize,
}
impl Menu {
    pub fn new(title: impl Into<String>) -> Self {
        Self {
            title: title.into(),
            body: String::new(),
            elements: Vec::new(),
            buttons: 0,
        }
    }
    pub fn body(mut self, body: impl Into<String>) -> Self {
        self.body = body.into();
        self
    }
    pub fn button(&mut self, button: Button) -> ButtonId {
        let id = ButtonId(self.buttons);
        self.buttons += 1;
        self.elements.push(MenuElement::Button(button));
        id
    }
    pub fn divider(&mut self) -> &mut Self {
        self.elements.push(MenuElement::Divider);
        self
    }
    pub fn header(&mut self, text: impl Into<String>) -> &mut Self {
        self.elements.push(MenuElement::Header(Header::new(text)));
        self
    }
    pub fn label(&mut self, text: impl Into<String>) -> &mut Self {
        self.elements.push(MenuElement::Label(Label::new(text)));
        self
    }
}
impl private::Sealed for Menu {}
impl Form for Menu {
    type Response = MenuResponse;
    fn request_json(&self) -> Vec<u8> {
        let elements: Vec<Value> = self
            .elements
            .iter()
            .map(|element| match element {
                MenuElement::Button(v) => v.json(),
                MenuElement::Divider => json!({"type":"divider","text":""}),
                MenuElement::Header(v) => json!({"type":"header","text":v.0}),
                MenuElement::Label(v) => json!({"type":"label","text":v.0}),
            })
            .collect();
        serde_json::to_vec(
            &json!({"type":"form","title":self.title,"content":self.body,"elements":elements}),
        )
        .unwrap_or_default()
    }
    fn parse_response(&self, data: &[u8]) -> Option<Self::Response> {
        let index: usize = serde_json::from_slice(data).ok()?;
        (index < self.buttons).then_some(MenuResponse {
            selected: ButtonId(index),
        })
    }
}

#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub enum ModalChoice {
    First,
    Second,
}
#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub struct ModalResponse {
    choice: ModalChoice,
}
impl ModalResponse {
    pub const fn choice(&self) -> ModalChoice {
        self.choice
    }
    pub const fn accepted(&self) -> bool {
        matches!(self.choice, ModalChoice::First)
    }
}
pub struct Modal {
    title: String,
    body: String,
    first: Button,
    second: Button,
}
impl Modal {
    pub fn new(title: impl Into<String>, first: Button, second: Button) -> Self {
        Self {
            title: title.into(),
            body: String::new(),
            first,
            second,
        }
    }
    pub fn yes_no(title: impl Into<String>) -> Self {
        Self::new(title, Button::new("gui.yes"), Button::new("gui.no"))
    }
    pub fn body(mut self, body: impl Into<String>) -> Self {
        self.body = body.into();
        self
    }
}
impl private::Sealed for Modal {}
impl Form for Modal {
    type Response = ModalResponse;
    fn request_json(&self) -> Vec<u8> {
        serde_json::to_vec(&json!({"type":"modal","title":self.title,"content":self.body,"button1":self.first.text,"button2":self.second.text})).unwrap_or_default()
    }
    fn parse_response(&self, data: &[u8]) -> Option<Self::Response> {
        Some(ModalResponse {
            choice: if serde_json::from_slice(data).ok()? {
                ModalChoice::First
            } else {
                ModalChoice::Second
            },
        })
    }
}

#[derive(Clone, Debug)]
pub struct Input {
    text: String,
    default: String,
    placeholder: String,
    tooltip: String,
}
impl Input {
    pub fn new(text: impl Into<String>) -> Self {
        Self {
            text: text.into(),
            default: String::new(),
            placeholder: String::new(),
            tooltip: String::new(),
        }
    }
    pub fn default(mut self, v: impl Into<String>) -> Self {
        self.default = v.into();
        self
    }
    pub fn placeholder(mut self, v: impl Into<String>) -> Self {
        self.placeholder = v.into();
        self
    }
    pub fn tooltip(mut self, v: impl Into<String>) -> Self {
        self.tooltip = v.into();
        self
    }
}
#[derive(Clone, Debug)]
pub struct Toggle {
    text: String,
    default: bool,
    tooltip: String,
}
impl Toggle {
    pub fn new(text: impl Into<String>, default: bool) -> Self {
        Self {
            text: text.into(),
            default,
            tooltip: String::new(),
        }
    }
    pub fn tooltip(mut self, v: impl Into<String>) -> Self {
        self.tooltip = v.into();
        self
    }
}
#[derive(Clone, Debug)]
pub struct Slider {
    text: String,
    min: f64,
    max: f64,
    step: f64,
    default: f64,
    tooltip: String,
}
impl Slider {
    pub fn new(text: impl Into<String>, min: f64, max: f64, step: f64, default: f64) -> Self {
        Self {
            text: text.into(),
            min,
            max,
            step,
            default,
            tooltip: String::new(),
        }
    }
    pub fn tooltip(mut self, v: impl Into<String>) -> Self {
        self.tooltip = v.into();
        self
    }
}
#[derive(Clone, Debug)]
pub struct Dropdown {
    text: String,
    options: Vec<String>,
    default: usize,
    tooltip: String,
}
impl Dropdown {
    pub fn new(
        text: impl Into<String>,
        options: impl IntoIterator<Item = impl Into<String>>,
        default: usize,
    ) -> Self {
        Self {
            text: text.into(),
            options: options.into_iter().map(Into::into).collect(),
            default,
            tooltip: String::new(),
        }
    }
    pub fn tooltip(mut self, v: impl Into<String>) -> Self {
        self.tooltip = v.into();
        self
    }
}
#[derive(Clone, Debug)]
pub struct StepSlider(Dropdown);
impl StepSlider {
    pub fn new(
        text: impl Into<String>,
        options: impl IntoIterator<Item = impl Into<String>>,
        default: usize,
    ) -> Self {
        Self(Dropdown::new(text, options, default))
    }
    pub fn tooltip(mut self, v: impl Into<String>) -> Self {
        self.0.tooltip = v.into();
        self
    }
}

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct Field<T> {
    index: usize,
    marker: PhantomData<fn() -> T>,
}
impl<T> Field<T> {
    fn new(index: usize) -> Self {
        Self {
            index,
            marker: PhantomData,
        }
    }
}
enum CustomElement {
    Divider,
    Header(String),
    Label(String),
    Input(Input),
    Toggle(Toggle),
    Slider(Slider),
    Dropdown(Dropdown),
    StepSlider(StepSlider),
}
pub struct Custom {
    title: String,
    elements: Vec<CustomElement>,
}
impl Custom {
    pub fn new(title: impl Into<String>) -> Self {
        Self {
            title: title.into(),
            elements: Vec::new(),
        }
    }
    pub fn divider(&mut self) -> &mut Self {
        self.elements.push(CustomElement::Divider);
        self
    }
    pub fn header(&mut self, text: impl Into<String>) -> &mut Self {
        self.elements.push(CustomElement::Header(text.into()));
        self
    }
    pub fn label(&mut self, text: impl Into<String>) -> &mut Self {
        self.elements.push(CustomElement::Label(text.into()));
        self
    }
    pub fn input(&mut self, v: Input) -> Field<String> {
        let f = Field::new(self.elements.len());
        self.elements.push(CustomElement::Input(v));
        f
    }
    pub fn toggle(&mut self, v: Toggle) -> Field<bool> {
        let f = Field::new(self.elements.len());
        self.elements.push(CustomElement::Toggle(v));
        f
    }
    pub fn slider(&mut self, v: Slider) -> Field<f64> {
        let f = Field::new(self.elements.len());
        self.elements.push(CustomElement::Slider(v));
        f
    }
    pub fn dropdown(&mut self, v: Dropdown) -> Field<usize> {
        let f = Field::new(self.elements.len());
        self.elements.push(CustomElement::Dropdown(v));
        f
    }
    pub fn step_slider(&mut self, v: StepSlider) -> Field<usize> {
        let f = Field::new(self.elements.len());
        self.elements.push(CustomElement::StepSlider(v));
        f
    }
}
pub struct CustomResponse {
    values: Vec<Option<Value>>,
}
pub trait ResponseValue: private::Sealed + Sized {
    fn read(value: &Value) -> Option<Self>;
}
impl private::Sealed for String {}
impl ResponseValue for String {
    fn read(v: &Value) -> Option<Self> {
        v.as_str().map(str::to_owned)
    }
}
impl private::Sealed for bool {}
impl ResponseValue for bool {
    fn read(v: &Value) -> Option<Self> {
        v.as_bool()
    }
}
impl private::Sealed for f64 {}
impl ResponseValue for f64 {
    fn read(v: &Value) -> Option<Self> {
        v.as_f64()
    }
}
impl private::Sealed for usize {}
impl ResponseValue for usize {
    fn read(v: &Value) -> Option<Self> {
        usize::try_from(v.as_u64()?).ok()
    }
}
impl CustomResponse {
    pub fn get<T: ResponseValue>(&self, field: Field<T>) -> Option<T> {
        T::read(self.values.get(field.index)?.as_ref()?)
    }
}
impl private::Sealed for Custom {}
impl Form for Custom {
    type Response = CustomResponse;
    fn request_json(&self) -> Vec<u8> {
        fn tooltip(mut v: Value, t: &str) -> Value {
            if !t.is_empty() {
                v["tooltip"] = json!(t)
            }
            v
        }
        let content:Vec<Value>=self.elements.iter().map(|e|match e{
            CustomElement::Divider=>json!({"type":"divider","text":""}),CustomElement::Header(v)=>json!({"type":"header","text":v}),CustomElement::Label(v)=>json!({"type":"label","text":v}),
            CustomElement::Input(v)=>tooltip(json!({"type":"input","text":v.text,"default":v.default,"placeholder":v.placeholder}),&v.tooltip),
            CustomElement::Toggle(v)=>tooltip(json!({"type":"toggle","text":v.text,"default":v.default}),&v.tooltip),
            CustomElement::Slider(v)=>tooltip(json!({"type":"slider","text":v.text,"min":v.min,"max":v.max,"step":v.step,"default":v.default}),&v.tooltip),
            CustomElement::Dropdown(v)=>tooltip(json!({"type":"dropdown","text":v.text,"options":v.options,"default":v.default}),&v.tooltip),
            CustomElement::StepSlider(v)=>tooltip(json!({"type":"step_slider","text":v.0.text,"steps":v.0.options,"default":v.0.default}),&v.0.tooltip),
        }).collect();
        serde_json::to_vec(&json!({"type":"custom_form","title":self.title,"content":content}))
            .unwrap_or_default()
    }
    fn parse_response(&self, data: &[u8]) -> Option<Self::Response> {
        let raw: Vec<Value> = serde_json::from_slice(data).ok()?;
        if raw.len() < self.elements.len() {
            return None;
        }
        let mut values = Vec::with_capacity(self.elements.len());
        for (e, v) in self.elements.iter().zip(raw) {
            values.push(match e {
                CustomElement::Divider | CustomElement::Header(_) | CustomElement::Label(_) => None,
                CustomElement::Input(_) => v.is_string().then_some(v),
                CustomElement::Toggle(_) => v.is_boolean().then_some(v),
                CustomElement::Slider(s) => v
                    .as_f64()
                    .filter(|n| *n >= s.min && *n <= s.max)
                    .map(Value::from),
                CustomElement::Dropdown(d) => v
                    .as_u64()
                    .filter(|i| (*i as usize) < d.options.len())
                    .map(Value::from),
                CustomElement::StepSlider(s) => v
                    .as_u64()
                    .filter(|i| (*i as usize) < s.0.options.len())
                    .map(Value::from),
            });
            if !matches!(
                e,
                CustomElement::Divider | CustomElement::Header(_) | CustomElement::Label(_)
            ) && values.last().is_some_and(Option::is_none)
            {
                return None;
            }
        }
        Some(CustomResponse { values })
    }
}

struct Completion<F: Form, C> {
    form: F,
    callback: C,
}
pub(crate) fn send<F, C>(player: &Player, form: F, callback: C)
where
    F: Form,
    C: FnOnce(Player, Option<F::Response>) + Send + 'static,
{
    let Some(host) = host_api() else { return };
    let Some(send) = host.player_form_send else {
        return;
    };
    let request = form.request_json();
    if request.is_empty() {
        return;
    }
    let completion = Box::new(Completion { form, callback });
    let context = Box::into_raw(completion).cast::<c_void>();
    let view = crate::__private::sys::DfFormView {
        request_json: bytes_view(&request),
        callback_context: context,
        response: Some(respond::<F, C>),
        drop: Some(drop_completion::<F, C>),
    };
    if unsafe {
        send(
            host.context,
            crate::current_invocation(),
            player.raw_id(),
            &view,
        )
    } != crate::__private::sys::DF_STATUS_OK
    {
        unsafe { drop(Box::from_raw(context.cast::<Completion<F, C>>())) };
    }
}
unsafe extern "C" fn respond<F, C>(
    context: *mut c_void,
    invocation: crate::__private::sys::DfInvocationId,
    submitter: crate::__private::sys::DfPlayerId,
    outcome: u32,
    response: crate::__private::sys::DfStringView,
) -> crate::__private::sys::DfStatus
where
    F: Form,
    C: FnOnce(Player, Option<F::Response>) + Send + 'static,
{
    if context.is_null() {
        return crate::__private::sys::DF_STATUS_ERROR;
    }
    let completion = unsafe { Box::from_raw(context.cast::<Completion<F, C>>()) };
    let result = std::panic::catch_unwind(std::panic::AssertUnwindSafe(|| {
        crate::with_invocation(invocation, || {
            let value = if outcome == crate::__private::sys::DF_FORM_RESPONSE_CLOSED {
                None
            } else if outcome == crate::__private::sys::DF_FORM_RESPONSE_SUBMITTED {
                let bytes = if response.len == 0 {
                    &[][..]
                } else if response.data.is_null() {
                    return false;
                } else {
                    unsafe { core::slice::from_raw_parts(response.data, response.len as usize) }
                };
                let Some(parsed) = completion.form.parse_response(bytes) else {
                    return false;
                };
                Some(parsed)
            } else {
                return false;
            };
            (completion.callback)(Player::from_id(submitter), value);
            true
        })
    }));
    if matches!(result, Ok(true)) {
        crate::__private::sys::DF_STATUS_OK
    } else {
        crate::__private::sys::DF_STATUS_ERROR
    }
}
unsafe extern "C" fn drop_completion<F, C>(context: *mut c_void)
where
    F: Form,
    C: FnOnce(Player, Option<F::Response>) + Send + 'static,
{
    if !context.is_null() {
        let _ = std::panic::catch_unwind(std::panic::AssertUnwindSafe(|| unsafe {
            drop(Box::from_raw(context.cast::<Completion<F, C>>()))
        }));
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn raw_form_round_trips_json_values() {
        let form = Raw::new(json!({
            "type": "modal",
            "title": "Raw",
            "content": "Escape hatch",
            "button1": "Yes",
            "button2": "No"
        }));
        assert_eq!(
            serde_json::from_slice::<Value>(&form.request_json()).ok(),
            Some(json!({
                "type": "modal",
                "title": "Raw",
                "content": "Escape hatch",
                "button1": "Yes",
                "button2": "No"
            }))
        );
        assert_eq!(
            form.parse_response(br#"{"accepted":true}"#),
            Some(json!({"accepted": true}))
        );
        assert_eq!(form.parse_response(b"not-json"), None);
    }

    #[test]
    fn menu_response_returns_a_typed_button_id() {
        let mut menu = Menu::new("T").body("B");
        menu.header("Header").divider().label("Label");
        let id = menu.button(Button::new("Go").image("https://example.com/a.png"));
        assert_eq!(
            serde_json::from_slice::<Value>(&menu.request_json()).ok(),
            Some(json!({
                "type": "form",
                "title": "T",
                "content": "B",
                "elements": [
                    {"type": "header", "text": "Header"},
                    {"type": "divider", "text": ""},
                    {"type": "label", "text": "Label"},
                    {"type": "button", "text": "Go", "image": {"type": "url", "data": "https://example.com/a.png"}}
                ]
            }))
        );
        assert_eq!(menu.parse_response(b"0").map(|v| v.selected()), Some(id));
        assert_eq!(menu.parse_response(b"1"), None);
    }

    #[test]
    fn modal_response_returns_a_typed_choice() {
        let modal = Modal::yes_no("T").body("B");
        assert_eq!(
            serde_json::from_slice::<Value>(&modal.request_json()).ok(),
            Some(json!({
                "type": "modal", "title": "T", "content": "B",
                "button1": "gui.yes", "button2": "gui.no"
            }))
        );
        assert_eq!(
            modal.parse_response(b"true").map(|v| v.choice()),
            Some(ModalChoice::First)
        );
        assert_eq!(
            modal.parse_response(b"false").map(|v| v.choice()),
            Some(ModalChoice::Second)
        );
        assert_eq!(modal.parse_response(b"0"), None);
    }

    #[test]
    fn custom_response_returns_each_field_type() {
        let mut custom = Custom::new("C");
        custom.header("Profile");
        custom.label("Typed");
        custom.divider();
        let name = custom.input(Input::new("Name"));
        let enabled = custom.toggle(Toggle::new("Enabled", false));
        let amount = custom.slider(Slider::new("Amount", 0.0, 10.0, 0.5, 0.0));
        let choice = custom.dropdown(Dropdown::new("Pick", ["A", "B"], 0));
        let step = custom.step_slider(StepSlider::new("Step", ["A", "B"], 0));
        assert_eq!(
            serde_json::from_slice::<Value>(&custom.request_json()).ok(),
            Some(json!({
                "type": "custom_form", "title": "C", "content": [
                    {"type": "header", "text": "Profile"},
                    {"type": "label", "text": "Typed"},
                    {"type": "divider", "text": ""},
                    {"type": "input", "text": "Name", "default": "", "placeholder": ""},
                    {"type": "toggle", "text": "Enabled", "default": false},
                    {"type": "slider", "text": "Amount", "min": 0.0, "max": 10.0, "step": 0.5, "default": 0.0},
                    {"type": "dropdown", "text": "Pick", "options": ["A", "B"], "default": 0},
                    {"type": "step_slider", "text": "Step", "steps": ["A", "B"], "default": 0}
                ]
            }))
        );
        let response = custom.parse_response(br#"[null,null,null,"alex",true,2.5,1,0]"#);
        assert_eq!(
            response.as_ref().and_then(|r| r.get(name)),
            Some("alex".to_owned())
        );
        assert_eq!(response.as_ref().and_then(|r| r.get(enabled)), Some(true));
        assert_eq!(response.as_ref().and_then(|r| r.get(amount)), Some(2.5));
        assert_eq!(response.as_ref().and_then(|r| r.get(choice)), Some(1));
        assert_eq!(response.and_then(|r| r.get(step)), Some(0));
        assert!(
            custom
                .parse_response(br#"[null,null,null,"alex",true,11,1,0]"#)
                .is_none()
        );
        assert!(
            custom
                .parse_response(br#"[null,null,null,"alex",true,2.5,2,0]"#)
                .is_none()
        );
    }
}
