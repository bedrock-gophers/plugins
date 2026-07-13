use std::{collections::BTreeMap, time::Duration};

use crate::{ItemStack, Potion, Rotation, Vec3, World, block};

const MAX_ENTITY_TYPE_BYTES: usize = 256;
const MAX_ENTITY_TAG_BYTES: usize = 4 << 10;

/// A stable identity for an entity managed by the host.
#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct EntityId {
    uuid: [u8; 16],
    generation: u64,
}

impl EntityId {
    pub const fn uuid_bytes(self) -> [u8; 16] {
        self.uuid
    }

    pub const fn generation(self) -> u64 {
        self.generation
    }

    pub(crate) const fn from_parts(uuid: [u8; 16], generation: u64) -> Self {
        Self { uuid, generation }
    }

    pub(crate) const fn raw(self) -> dragonfly_plugin_sys::DfEntityId {
        dragonfly_plugin_sys::DfEntityId {
            bytes: self.uuid,
            generation: self.generation,
        }
    }
}

impl From<dragonfly_plugin_sys::DfEntityId> for EntityId {
    fn from(value: dragonfly_plugin_sys::DfEntityId) -> Self {
        Self::from_parts(value.bytes, value.generation)
    }
}

/// A persistent handle to an entity.
///
/// The handle is safe to retain. Entity state is resolved by the host when an
/// operation is performed, so a handle may become stale after its entity is
/// removed.
#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct Entity {
    id: EntityId,
}

impl Entity {
    pub const fn id(self) -> EntityId {
        self.id
    }

    pub(crate) const fn from_id(id: EntityId) -> Self {
        Self { id }
    }

    pub(crate) const fn raw_id(self) -> dragonfly_plugin_sys::DfEntityId {
        self.id.raw()
    }

    pub fn state(&self) -> Option<EntityState> {
        let host = crate::host_api()?;
        let read = host.entity_state?;
        let mut entity_type = [0; MAX_ENTITY_TYPE_BYTES];
        let mut name_tag = vec![0; MAX_ENTITY_TAG_BYTES];
        let mut raw = dragonfly_plugin_sys::DfEntityState {
            position: dragonfly_plugin_sys::DfVec3::default(),
            rotation: dragonfly_plugin_sys::DfRotation::default(),
            velocity: dragonfly_plugin_sys::DfVec3::default(),
            capabilities: 0,
            world: dragonfly_plugin_sys::DfWorldId::default(),
            entity_type: crate::bytes_buffer(&mut entity_type),
            name_tag: crate::bytes_buffer(&mut name_tag),
        };
        let status = unsafe {
            read(
                host.context,
                crate::current_invocation(),
                self.raw_id(),
                &mut raw,
            )
        };
        if status != dragonfly_plugin_sys::DF_STATUS_OK {
            return None;
        }
        let type_length = usize::try_from(raw.entity_type.len).ok()?;
        let name_length = usize::try_from(raw.name_tag.len).ok()?;
        Some(EntityState {
            entity_type: String::from_utf8(entity_type.get(..type_length)?.to_vec()).ok()?,
            position: raw.position.into(),
            rotation: raw.rotation.into(),
            world: World::from_raw(raw.world.value),
            velocity: (raw.capabilities & dragonfly_plugin_sys::DF_ENTITY_HAS_VELOCITY != 0)
                .then(|| raw.velocity.into()),
            name_tag: if raw.capabilities & dragonfly_plugin_sys::DF_ENTITY_HAS_NAME_TAG != 0 {
                Some(String::from_utf8(name_tag.get(..name_length)?.to_vec()).ok()?)
            } else {
                None
            },
            can_teleport: raw.capabilities & dragonfly_plugin_sys::DF_ENTITY_CAN_TELEPORT != 0,
        })
    }

    pub fn entity_type(&self) -> Option<String> {
        self.state().map(|state| state.entity_type)
    }
    pub fn position(&self) -> Option<Vec3> {
        self.state().map(|state| state.position)
    }
    pub fn rotation(&self) -> Option<Rotation> {
        self.state().map(|state| state.rotation)
    }
    pub fn world(&self) -> Option<World> {
        self.state()?.world
    }
    pub fn velocity(&self) -> Option<Vec3> {
        self.state()?.velocity
    }
    pub fn name_tag(&self) -> Option<String> {
        self.state()?.name_tag
    }

    pub fn teleport(&self, position: Vec3) {
        let Some(host) = crate::host_api() else {
            return;
        };
        let Some(call) = host.entity_teleport else {
            return;
        };
        let _ = unsafe {
            call(
                host.context,
                crate::current_invocation(),
                self.raw_id(),
                position.into(),
            )
        };
    }

    pub fn set_velocity(&self, velocity: Vec3) {
        let Some(host) = crate::host_api() else {
            return;
        };
        let Some(call) = host.entity_velocity_set else {
            return;
        };
        let _ = unsafe {
            call(
                host.context,
                crate::current_invocation(),
                self.raw_id(),
                velocity.into(),
            )
        };
    }

    pub fn set_name_tag(&self, name_tag: &str) {
        let Some(host) = crate::host_api() else {
            return;
        };
        let Some(call) = host.entity_name_tag_set else {
            return;
        };
        let _ = unsafe {
            call(
                host.context,
                crate::current_invocation(),
                self.raw_id(),
                crate::string_view_from_str(name_tag),
            )
        };
    }

    pub fn despawn(self) {
        let Some(host) = crate::host_api() else {
            return;
        };
        let Some(call) = host.entity_despawn else {
            return;
        };
        let _ = unsafe { call(host.context, crate::current_invocation(), self.raw_id()) };
    }
}

impl From<dragonfly_plugin_sys::DfEntityId> for Entity {
    fn from(value: dragonfly_plugin_sys::DfEntityId) -> Self {
        Self::from_id(value.into())
    }
}

/// A snapshot of the common state exposed by an entity handle.
#[derive(Clone, Debug, Default, PartialEq)]
pub struct EntityState {
    pub entity_type: String,
    pub position: Vec3,
    pub rotation: Rotation,
    pub world: Option<World>,
    pub velocity: Option<Vec3>,
    pub name_tag: Option<String>,
    pub can_teleport: bool,
}

/// Physics applied by Dragonfly after a plugin tick callback.
#[derive(Clone, Copy, Debug, PartialEq)]
pub struct Physics {
    pub gravity: f64,
    pub drag: f64,
    pub drag_before_gravity: bool,
}

impl Physics {
    pub const fn passive(gravity: f64, drag: f64) -> Self {
        Self {
            gravity,
            drag,
            drag_before_gravity: true,
        }
    }

    pub const fn drag_after_gravity(mut self) -> Self {
        self.drag_before_gravity = false;
        self
    }
}

/// Failure to decode or encode plugin-owned persistent entity state.
#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub struct StateError;

/// Read-only plugin-owned entity state loaded from disk.
pub struct SavedState<'a> {
    version: u32,
    values: BTreeMap<String, crate::Value>,
    _lifetime: core::marker::PhantomData<&'a [u8]>,
}

impl<'a> SavedState<'a> {
    #[doc(hidden)]
    pub fn decode(version: u32, bytes: &'a [u8]) -> Result<Self, StateError> {
        let values = if bytes.is_empty() {
            BTreeMap::new()
        } else {
            crate::item_nbt::decode_values(bytes).map_err(|_| StateError)?
        };
        Ok(Self {
            version,
            values,
            _lifetime: core::marker::PhantomData,
        })
    }

    pub const fn version(&self) -> u32 {
        self.version
    }

    pub fn i64(&self, key: &str) -> Option<i64> {
        match self.values.get(key)? {
            crate::Value::Long(value) => Some(*value),
            crate::Value::Int(value) => Some(i64::from(*value)),
            _ => None,
        }
    }

    pub fn f64(&self, key: &str) -> Option<f64> {
        match self.values.get(key)? {
            crate::Value::Double(value) => Some(*value),
            crate::Value::Float(value) => Some(f64::from(*value)),
            _ => None,
        }
    }

    pub fn bool(&self, key: &str) -> Option<bool> {
        match self.values.get(key)? {
            crate::Value::Byte(value) => Some(*value != 0),
            _ => None,
        }
    }

    pub fn str(&self, key: &str) -> Option<&str> {
        match self.values.get(key)? {
            crate::Value::String(value) => Some(value),
            _ => None,
        }
    }

    pub fn value(&self, key: &str) -> Option<&crate::Value> {
        self.values.get(key)
    }
}

/// Builder for plugin-owned persistent entity state.
#[derive(Default)]
pub struct SavedStateMut {
    values: BTreeMap<String, crate::Value>,
}

impl SavedStateMut {
    pub const fn new() -> Self {
        Self {
            values: BTreeMap::new(),
        }
    }

    pub fn set(&mut self, key: impl Into<String>, value: impl Into<crate::Value>) {
        self.values.insert(key.into(), value.into());
    }

    pub fn remove(&mut self, key: &str) {
        self.values.remove(key);
    }

    #[doc(hidden)]
    pub fn encode(&self) -> Result<Vec<u8>, StateError> {
        crate::item_nbt::encode_values(&self.values).map_err(|_| StateError)
    }
}

/// Context passed to a plugin-owned entity tick.
pub struct TickContext<'a> {
    input: &'a dragonfly_plugin_sys::DfEntityTickInput,
    state: &'a mut dragonfly_plugin_sys::DfEntityTickState,
}

impl<'a> TickContext<'a> {
    #[doc(hidden)]
    pub unsafe fn from_raw(
        input: &'a dragonfly_plugin_sys::DfEntityTickInput,
        state: &'a mut dragonfly_plugin_sys::DfEntityTickState,
    ) -> Self {
        Self { input, state }
    }

    pub fn entity(&self) -> Entity {
        Entity::from_id(self.input.entity.into())
    }

    pub const fn current(&self) -> i64 {
        self.input.current
    }

    pub fn age(&self) -> Duration {
        Duration::from_millis(self.input.age_milliseconds)
    }

    pub fn despawn(&mut self) {
        self.state.despawn = 1;
    }
}

/// Mutable pre-damage callback for a plugin-owned living entity.
pub struct Hurt<'a> {
    input: &'a dragonfly_plugin_sys::DfEntityHurtInput,
    state: &'a mut dragonfly_plugin_sys::DfEntityHurtState,
}

impl<'a> Hurt<'a> {
    #[doc(hidden)]
    pub unsafe fn from_raw(
        input: &'a dragonfly_plugin_sys::DfEntityHurtInput,
        state: &'a mut dragonfly_plugin_sys::DfEntityHurtState,
    ) -> Self {
        Self { input, state }
    }

    pub fn entity(&self) -> Entity {
        Entity::from_id(self.input.entity.into())
    }

    pub fn damage_source(&self) -> crate::damage::Source<'_> {
        unsafe { crate::damage::Source::from_raw(&self.input.source) }.unwrap_or_else(|| {
            crate::damage::Custom::new(
                "unknown",
                crate::damage::Traits::default(),
                crate::damage::AffectedProtections::NONE,
            )
            .into()
        })
    }

    pub const fn health(&self) -> f64 {
        self.input.health
    }

    pub const fn max_health(&self) -> f64 {
        self.input.max_health
    }

    pub const fn damage(&self) -> f64 {
        self.state.damage
    }

    pub fn set_damage(&mut self, damage: f64) {
        self.state.damage = if damage.is_finite() {
            damage.max(0.0)
        } else {
            0.0
        };
    }

    pub const fn cancelled(&self) -> bool {
        self.state.cancelled != 0
    }

    pub fn cancel(&mut self) {
        self.state.cancelled = 1;
    }
}

/// Mutable pre-heal callback for a plugin-owned living entity.
pub struct Heal<'a> {
    input: &'a dragonfly_plugin_sys::DfEntityHealInput,
    state: &'a mut dragonfly_plugin_sys::DfEntityHealState,
}

impl<'a> Heal<'a> {
    #[doc(hidden)]
    pub unsafe fn from_raw(
        input: &'a dragonfly_plugin_sys::DfEntityHealInput,
        state: &'a mut dragonfly_plugin_sys::DfEntityHealState,
    ) -> Self {
        Self { input, state }
    }

    pub fn entity(&self) -> Entity {
        Entity::from_id(self.input.entity.into())
    }

    pub fn healing_source(&self) -> crate::healing::Source<'_> {
        unsafe { crate::healing::Source::from_raw(&self.input.source) }
            .unwrap_or_else(|| crate::healing::Custom::new("unknown").into())
    }

    pub const fn health_before(&self) -> f64 {
        self.input.health
    }

    pub const fn max_health(&self) -> f64 {
        self.input.max_health
    }

    pub const fn amount(&self) -> f64 {
        self.state.health
    }

    pub fn set_amount(&mut self, amount: f64) {
        self.state.health = if amount.is_finite() {
            amount.max(0.0)
        } else {
            0.0
        };
    }

    pub const fn cancelled(&self) -> bool {
        self.state.cancelled != 0
    }

    pub fn cancel(&mut self) {
        self.state.cancelled = 1;
    }
}

/// Mutable lethal-damage callback for a plugin-owned living entity.
pub struct Death<'a> {
    input: &'a dragonfly_plugin_sys::DfEntityDeathInput,
    state: &'a mut dragonfly_plugin_sys::DfEntityDeathState,
}

impl<'a> Death<'a> {
    #[doc(hidden)]
    pub unsafe fn from_raw(
        input: &'a dragonfly_plugin_sys::DfEntityDeathInput,
        state: &'a mut dragonfly_plugin_sys::DfEntityDeathState,
    ) -> Self {
        Self { input, state }
    }

    pub fn entity(&self) -> Entity {
        Entity::from_id(self.input.entity.into())
    }

    pub fn damage_source(&self) -> crate::damage::Source<'_> {
        unsafe { crate::damage::Source::from_raw(&self.input.source) }.unwrap_or_else(|| {
            crate::damage::Custom::new(
                "unknown",
                crate::damage::Traits::default(),
                crate::damage::AffectedProtections::NONE,
            )
            .into()
        })
    }

    pub const fn health(&self) -> f64 {
        self.input.health
    }

    pub const fn damage(&self) -> f64 {
        self.input.damage
    }

    pub const fn cancelled(&self) -> bool {
        self.state.cancelled != 0
    }

    pub fn cancel(&mut self) {
        self.state.cancelled = 1;
    }
}

/// Behaviour of a plugin-owned non-living ticking entity.
pub trait Ticking: Default + Send + 'static {
    const PHYSICS: Option<Physics> = None;
    const STATE_VERSION: u32 = 0;

    fn load(_state: &SavedState<'_>) -> Self {
        Self::default()
    }
    fn save(&self, _state: &mut SavedStateMut) {}
    fn tick(&mut self, _context: &mut TickContext<'_>) {}
}

/// Behaviour of a plugin-owned living entity.
pub trait Living: Default + Send + 'static {
    const INITIAL_HEALTH: f64 = 20.0;
    const MAX_HEALTH: f64 = 20.0;
    const SPEED: f64 = 0.1;
    const PHYSICS: Option<Physics> = Some(Physics::passive(0.08, 0.02));
    const STATE_VERSION: u32 = 0;

    fn load(_state: &SavedState<'_>) -> Self {
        Self::default()
    }
    fn save(&self, _state: &mut SavedStateMut) {}
    fn tick(&mut self, _context: &mut TickContext<'_>) {}
    fn hurt(&mut self, _event: &mut Hurt<'_>) {}
    fn heal(&mut self, _event: &mut Heal<'_>) {}
    fn death(&mut self, _event: &mut Death<'_>) {}
}

trait AdvancedAdapter {
    type State: Default + Send + 'static;
    const STATE_VERSION: u32;

    fn load(state: &SavedState<'_>) -> Self::State;
    fn save(value: &Self::State, state: &mut SavedStateMut);
    fn tick(value: &mut Self::State, context: &mut TickContext<'_>);
    fn hurt(_value: &mut Self::State, _event: &mut Hurt<'_>) -> bool {
        false
    }
    fn heal(_value: &mut Self::State, _event: &mut Heal<'_>) -> bool {
        false
    }
    fn death(_value: &mut Self::State, _event: &mut Death<'_>) -> bool {
        false
    }
}

struct TickingAdapter<T>(core::marker::PhantomData<T>);

impl<T: Ticking> AdvancedAdapter for TickingAdapter<T> {
    type State = T;
    const STATE_VERSION: u32 = T::STATE_VERSION;

    fn load(state: &SavedState<'_>) -> T {
        T::load(state)
    }
    fn save(value: &T, state: &mut SavedStateMut) {
        value.save(state);
    }
    fn tick(value: &mut T, context: &mut TickContext<'_>) {
        value.tick(context);
    }
}

struct LivingAdapter<T>(core::marker::PhantomData<T>);

impl<T: Living> AdvancedAdapter for LivingAdapter<T> {
    type State = T;
    const STATE_VERSION: u32 = T::STATE_VERSION;

    fn load(state: &SavedState<'_>) -> T {
        T::load(state)
    }
    fn save(value: &T, state: &mut SavedStateMut) {
        value.save(state);
    }
    fn tick(value: &mut T, context: &mut TickContext<'_>) {
        value.tick(context);
    }
    fn hurt(value: &mut T, event: &mut Hurt<'_>) -> bool {
        value.hurt(event);
        true
    }
    fn heal(value: &mut T, event: &mut Heal<'_>) -> bool {
        value.heal(event);
        true
    }
    fn death(value: &mut T, event: &mut Death<'_>) -> bool {
        value.death(event);
        true
    }
}

#[doc(hidden)]
pub unsafe fn handle_ticking<T: Ticking>(
    operation: u32,
    instance: u64,
    input: *const core::ffi::c_void,
    state: *mut core::ffi::c_void,
) -> dragonfly_plugin_sys::DfStatus {
    unsafe { handle_advanced::<TickingAdapter<T>>(operation, instance, input, state) }
}

#[doc(hidden)]
pub unsafe fn handle_living<T: Living>(
    operation: u32,
    instance: u64,
    input: *const core::ffi::c_void,
    state: *mut core::ffi::c_void,
) -> dragonfly_plugin_sys::DfStatus {
    unsafe { handle_advanced::<LivingAdapter<T>>(operation, instance, input, state) }
}

unsafe fn handle_advanced<A: AdvancedAdapter>(
    operation: u32,
    instance: u64,
    input: *const core::ffi::c_void,
    state: *mut core::ffi::c_void,
) -> dragonfly_plugin_sys::DfStatus {
    use dragonfly_plugin_sys as sys;
    match operation {
        sys::DF_ENTITY_OPERATION_ADOPT => {
            if instance == 0 {
                sys::DF_STATUS_ERROR
            } else {
                sys::DF_STATUS_OK
            }
        }
        sys::DF_ENTITY_OPERATION_LOAD => {
            let (Some(input), Some(output)) = (
                unsafe { input.cast::<sys::DfEntityLoadInput>().as_ref() },
                unsafe { state.cast::<sys::DfEntityLoadState>().as_mut() },
            ) else {
                return sys::DF_STATUS_ERROR;
            };
            let Ok(length) = usize::try_from(input.data.len) else {
                return sys::DF_STATUS_ERROR;
            };
            if length > 16 << 20 || (length != 0 && input.data.data.is_null()) {
                return sys::DF_STATUS_ERROR;
            }
            let bytes = if length == 0 {
                &[]
            } else {
                unsafe { core::slice::from_raw_parts(input.data.data, length) }
            };
            let Ok(saved) = SavedState::decode(input.version, bytes) else {
                return sys::DF_STATUS_ERROR;
            };
            output.instance = Box::into_raw(Box::new(A::load(&saved))) as usize as u64;
            sys::DF_STATUS_OK
        }
        sys::DF_ENTITY_OPERATION_SAVE => {
            let (Some(value), Some(output)) = (
                unsafe { (instance as usize as *const A::State).as_ref() },
                unsafe { state.cast::<sys::DfEntitySaveState>().as_mut() },
            ) else {
                return sys::DF_STATUS_ERROR;
            };
            let mut saved = SavedStateMut::new();
            A::save(value, &mut saved);
            let Ok(bytes) = saved.encode() else {
                return sys::DF_STATUS_ERROR;
            };
            output.version = A::STATE_VERSION;
            output.data.len = bytes.len() as u64;
            let Ok(capacity) = usize::try_from(output.data.capacity) else {
                return sys::DF_STATUS_ERROR;
            };
            if bytes.len() > capacity || (bytes.len() != 0 && output.data.data.is_null()) {
                return sys::DF_STATUS_ERROR;
            }
            if !bytes.is_empty() {
                unsafe {
                    core::ptr::copy_nonoverlapping(bytes.as_ptr(), output.data.data, bytes.len())
                };
            }
            sys::DF_STATUS_OK
        }
        sys::DF_ENTITY_OPERATION_TICK => {
            let (Some(value), Some(input), Some(output)) = (
                unsafe { (instance as usize as *mut A::State).as_mut() },
                unsafe { input.cast::<sys::DfEntityTickInput>().as_ref() },
                unsafe { state.cast::<sys::DfEntityTickState>().as_mut() },
            ) else {
                return sys::DF_STATUS_ERROR;
            };
            let mut context = unsafe { TickContext::from_raw(input, output) };
            crate::with_invocation(input.invocation, || A::tick(value, &mut context));
            sys::DF_STATUS_OK
        }
        sys::DF_ENTITY_OPERATION_HURT => {
            let (Some(value), Some(input), Some(output)) = (
                unsafe { (instance as usize as *mut A::State).as_mut() },
                unsafe { input.cast::<sys::DfEntityHurtInput>().as_ref() },
                unsafe { state.cast::<sys::DfEntityHurtState>().as_mut() },
            ) else {
                return sys::DF_STATUS_ERROR;
            };
            let mut event = unsafe { Hurt::from_raw(input, output) };
            if crate::with_invocation(input.invocation, || A::hurt(value, &mut event)) {
                sys::DF_STATUS_OK
            } else {
                sys::DF_STATUS_ERROR
            }
        }
        sys::DF_ENTITY_OPERATION_HEAL => {
            let (Some(value), Some(input), Some(output)) = (
                unsafe { (instance as usize as *mut A::State).as_mut() },
                unsafe { input.cast::<sys::DfEntityHealInput>().as_ref() },
                unsafe { state.cast::<sys::DfEntityHealState>().as_mut() },
            ) else {
                return sys::DF_STATUS_ERROR;
            };
            let mut event = unsafe { Heal::from_raw(input, output) };
            if crate::with_invocation(input.invocation, || A::heal(value, &mut event)) {
                sys::DF_STATUS_OK
            } else {
                sys::DF_STATUS_ERROR
            }
        }
        sys::DF_ENTITY_OPERATION_DEATH => {
            let (Some(value), Some(input), Some(output)) = (
                unsafe { (instance as usize as *mut A::State).as_mut() },
                unsafe { input.cast::<sys::DfEntityDeathInput>().as_ref() },
                unsafe { state.cast::<sys::DfEntityDeathState>().as_mut() },
            ) else {
                return sys::DF_STATUS_ERROR;
            };
            let mut event = unsafe { Death::from_raw(input, output) };
            if crate::with_invocation(input.invocation, || A::death(value, &mut event)) {
                sys::DF_STATUS_OK
            } else {
                sys::DF_STATUS_ERROR
            }
        }
        sys::DF_ENTITY_OPERATION_DESTROY => {
            if instance == 0 {
                return sys::DF_STATUS_ERROR;
            }
            drop(unsafe { Box::from_raw(instance as usize as *mut A::State) });
            sys::DF_STATUS_OK
        }
        _ => sys::DF_STATUS_ERROR,
    }
}

/// Common options applied when an entity is spawned.
#[derive(Clone, Debug, PartialEq)]
pub struct SpawnOptions {
    pub position: Vec3,
    pub rotation: Option<Rotation>,
    pub velocity: Option<Vec3>,
    pub name_tag: Option<String>,
}

impl SpawnOptions {
    pub const fn new(position: Vec3) -> Self {
        Self {
            position,
            rotation: None,
            velocity: None,
            name_tag: None,
        }
    }

    pub const fn rotation(mut self, rotation: Rotation) -> Self {
        self.rotation = Some(rotation);
        self
    }

    pub const fn velocity(mut self, velocity: Vec3) -> Self {
        self.velocity = Some(velocity);
        self
    }

    pub fn name_tag(mut self, name_tag: impl Into<String>) -> Self {
        self.name_tag = Some(name_tag.into());
        self
    }
}

/// A stationary text label.
#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Text {
    text: String,
}

impl Text {
    pub fn new(text: impl Into<String>) -> Self {
        Self { text: text.into() }
    }

    pub fn text(&self) -> &str {
        &self.text
    }
}

/// A lightning bolt.
#[derive(Clone, Copy, Debug, PartialEq)]
pub struct Lightning {
    damage: f64,
    block_fire: bool,
    entity_fire_duration: Duration,
}

impl Lightning {
    pub const fn new() -> Self {
        Self {
            damage: 5.0,
            block_fire: true,
            entity_fire_duration: Duration::from_secs(8),
        }
    }

    pub const fn damage(mut self, damage: f64) -> Self {
        self.damage = damage;
        self
    }

    pub const fn block_fire(mut self, block_fire: bool) -> Self {
        self.block_fire = block_fire;
        self
    }

    pub const fn entity_fire_duration(mut self, duration: Duration) -> Self {
        self.entity_fire_duration = duration;
        self
    }

    pub const fn damage_value(self) -> f64 {
        self.damage
    }

    pub const fn creates_block_fire(self) -> bool {
        self.block_fire
    }

    pub const fn fire_duration(self) -> Duration {
        self.entity_fire_duration
    }
}

impl Default for Lightning {
    fn default() -> Self {
        Self::new()
    }
}

/// Primed TNT that explodes after its fuse expires.
#[allow(clippy::upper_case_acronyms)]
#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub struct TNT {
    fuse: Duration,
}

impl TNT {
    pub const fn new(fuse: Duration) -> Self {
        Self { fuse }
    }

    pub const fn fuse(self) -> Duration {
        self.fuse
    }
}

/// A single experience orb.
#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub struct ExperienceOrb {
    experience: u32,
}

impl ExperienceOrb {
    pub const fn new(experience: u32) -> Self {
        Self { experience }
    }

    pub const fn experience(self) -> u32 {
        self.experience
    }
}

/// An item stack dropped into a world.
#[derive(Clone, Debug, PartialEq)]
pub struct DroppedItem {
    stack: ItemStack,
    pickup_delay: Option<Duration>,
}

impl DroppedItem {
    pub const fn new(stack: ItemStack) -> Self {
        Self {
            stack,
            pickup_delay: None,
        }
    }

    pub const fn pickup_delay(mut self, delay: Duration) -> Self {
        self.pickup_delay = Some(delay);
        self
    }

    pub const fn stack(&self) -> &ItemStack {
        &self.stack
    }

    pub const fn delay(&self) -> Option<Duration> {
        self.pickup_delay
    }
}

/// A block state falling as an entity.
#[derive(Clone, Debug, Eq, PartialEq)]
pub struct FallingBlock {
    block: block::Block,
}

impl FallingBlock {
    pub fn new(block: impl Into<block::Block>) -> Self {
        Self {
            block: block.into(),
        }
    }

    pub const fn block(&self) -> &block::Block {
        &self.block
    }
}

#[derive(Clone, Copy, Debug, PartialEq)]
pub struct Arrow {
    owner: Entity,
    damage: f64,
    tip: Potion,
    critical: bool,
    disable_pickup: bool,
    obtain_on_pickup: bool,
    punch_level: u32,
    piercing_level: u32,
}

impl Arrow {
    pub const fn new(owner: Entity) -> Self {
        Self {
            owner,
            damage: 2.0,
            tip: Potion::Water,
            critical: false,
            disable_pickup: false,
            obtain_on_pickup: false,
            punch_level: 0,
            piercing_level: 0,
        }
    }

    pub const fn damage(mut self, damage: f64) -> Self {
        self.damage = damage;
        self
    }
    pub const fn tip(mut self, tip: Potion) -> Self {
        self.tip = tip;
        self
    }
    pub const fn critical(mut self, critical: bool) -> Self {
        self.critical = critical;
        self
    }
    pub const fn disable_pickup(mut self, disabled: bool) -> Self {
        self.disable_pickup = disabled;
        self
    }
    pub const fn obtain_on_pickup(mut self, obtain: bool) -> Self {
        self.obtain_on_pickup = obtain;
        self
    }
    pub const fn punch_level(mut self, level: u32) -> Self {
        self.punch_level = level;
        self
    }
    pub const fn piercing_level(mut self, level: u32) -> Self {
        self.piercing_level = level;
        self
    }
}

macro_rules! owner_projectile {
    ($name:ident) => {
        #[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
        pub struct $name {
            owner: Entity,
        }
        impl $name {
            pub const fn new(owner: Entity) -> Self {
                Self { owner }
            }
            pub const fn owner(self) -> Entity {
                self.owner
            }
        }
    };
}

owner_projectile!(Egg);
owner_projectile!(Snowball);
owner_projectile!(EnderPearl);
owner_projectile!(BottleOfEnchanting);

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct SplashPotion {
    owner: Entity,
    potion: Potion,
}

impl SplashPotion {
    pub const fn new(owner: Entity, potion: Potion) -> Self {
        Self { owner, potion }
    }
}

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct LingeringPotion {
    owner: Entity,
    potion: Potion,
}

impl LingeringPotion {
    pub const fn new(owner: Entity, potion: Potion) -> Self {
        Self { owner, potion }
    }
}

/// Static definition of a plugin-owned base entity type.
#[derive(Clone, Copy, Debug, PartialEq)]
pub struct Definition {
    identifier: &'static str,
    network_identifier: &'static str,
    bounds: [f64; 6],
}

impl Definition {
    #[doc(hidden)]
    pub const fn new(
        identifier: &'static str,
        network_identifier: &'static str,
        bounds: [f64; 6],
    ) -> Self {
        Self {
            identifier,
            network_identifier,
            bounds,
        }
    }

    pub const fn identifier(self) -> &'static str {
        self.identifier
    }

    pub const fn network_identifier(self) -> &'static str {
        self.network_identifier
    }

    /// Returns `[min_x, min_y, min_z, max_x, max_y, max_z]`.
    pub const fn bounds(self) -> [f64; 6] {
        self.bounds
    }
}

/// Implemented by unit structs annotated with `#[dragonfly::entity]`.
pub trait CustomType {
    const DEFINITION: Definition;

    #[doc(hidden)]
    fn encode_spawn(self) -> EncodedSpawnable
    where
        Self: Sized,
    {
        EncodedSpawnable {
            kind: SpawnKind::Custom,
            payload: SpawnPayload::Custom {
                identifier: Self::DEFINITION.identifier(),
                instance: None,
            },
        }
    }
}

#[doc(hidden)]
pub type EntityHandler = unsafe fn(
    operation: u32,
    instance: u64,
    input: *const core::ffi::c_void,
    state: *mut core::ffi::c_void,
) -> dragonfly_plugin_sys::DfStatus;

pub struct RegisteredType {
    descriptor: dragonfly_plugin_sys::DfEntityTypeDescriptorV2,
    handler: Option<EntityHandler>,
}

impl RegisteredType {
    pub const fn new(definition: Definition) -> Self {
        let bounds = definition.bounds();
        Self {
            descriptor: dragonfly_plugin_sys::DfEntityTypeDescriptorV2 {
                save_id: dragonfly_plugin_sys::DfStringView {
                    data: definition.identifier().as_ptr(),
                    len: definition.identifier().len() as u64,
                },
                network_id: dragonfly_plugin_sys::DfStringView {
                    data: definition.network_identifier().as_ptr(),
                    len: definition.network_identifier().len() as u64,
                },
                min: dragonfly_plugin_sys::DfVec3 {
                    x: bounds[0],
                    y: bounds[1],
                    z: bounds[2],
                },
                max: dragonfly_plugin_sys::DfVec3 {
                    x: bounds[3],
                    y: bounds[4],
                    z: bounds[5],
                },
                type_key: 0,
                family: dragonfly_plugin_sys::DF_ENTITY_FAMILY_BASE,
                callback_flags: 0,
                initial_health: 0.0,
                max_health: 0.0,
                speed: 0.0,
                state_version: 0,
                physics_flags: 0,
                gravity: 0.0,
                drag: 0.0,
            },
            handler: None,
        }
    }

    #[doc(hidden)]
    pub const fn advanced(
        definition: Definition,
        family: u32,
        callback_flags: u32,
        initial_health: f64,
        max_health: f64,
        speed: f64,
        state_version: u32,
        physics: Option<Physics>,
        handler: EntityHandler,
    ) -> Self {
        let mut registered = Self::new(definition);
        registered.descriptor.family = family;
        registered.descriptor.callback_flags = callback_flags;
        registered.descriptor.initial_health = initial_health;
        registered.descriptor.max_health = max_health;
        registered.descriptor.speed = speed;
        registered.descriptor.state_version = state_version;
        if let Some(physics) = physics {
            registered.descriptor.physics_flags = dragonfly_plugin_sys::DF_ENTITY_PHYSICS_ENABLED
                | if physics.drag_before_gravity {
                    dragonfly_plugin_sys::DF_ENTITY_PHYSICS_DRAG_BEFORE_GRAVITY
                } else {
                    0
                };
            registered.descriptor.gravity = physics.gravity;
            registered.descriptor.drag = physics.drag;
        }
        registered.handler = Some(handler);
        registered
    }

    #[doc(hidden)]
    pub const fn descriptor(&self) -> dragonfly_plugin_sys::DfEntityTypeDescriptorV2 {
        self.descriptor
    }

    #[doc(hidden)]
    pub const fn handler(&self) -> Option<EntityHandler> {
        self.handler
    }
}

// The descriptor only contains pointers to static UTF-8 strings.
unsafe impl Sync for RegisteredType {}

#[linkme::distributed_slice]
#[doc(hidden)]
pub static REGISTERED_TYPES: [RegisteredType] = [..];

mod private {
    pub trait Sealed {}
}

/// A typed descriptor that the host can turn into a Dragonfly entity.
///
/// This trait is sealed so every descriptor has a stable host encoding.
pub trait Spawnable: private::Sealed {
    #[doc(hidden)]
    fn encode(self) -> EncodedSpawnable;
}

#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub(crate) enum SpawnKind {
    Text,
    Lightning,
    Tnt,
    ExperienceOrb,
    DroppedItem,
    FallingBlock,
    Arrow,
    Egg,
    Snowball,
    EnderPearl,
    BottleOfEnchanting,
    SplashPotion,
    LingeringPotion,
    Custom,
}

#[derive(Debug, PartialEq)]
pub(crate) enum SpawnPayload {
    Text(String),
    Lightning {
        damage: f64,
        block_fire: bool,
        entity_fire_duration: Duration,
    },
    Tnt {
        fuse: Duration,
    },
    ExperienceOrb {
        experience: u32,
    },
    DroppedItem {
        stack: ItemStack,
        pickup_delay: Option<Duration>,
    },
    FallingBlock(block::Block),
    Arrow(Arrow),
    Egg(Entity),
    Snowball(Entity),
    EnderPearl(Entity),
    BottleOfEnchanting(Entity),
    SplashPotion {
        owner: Entity,
        potion: Potion,
    },
    LingeringPotion {
        owner: Entity,
        potion: Potion,
    },
    Custom {
        identifier: &'static str,
        instance: Option<OwnedInstance>,
    },
}

pub(crate) struct OwnedInstance {
    raw: u64,
    destroy: unsafe fn(u64),
}

impl core::fmt::Debug for OwnedInstance {
    fn fmt(&self, formatter: &mut core::fmt::Formatter<'_>) -> core::fmt::Result {
        formatter
            .debug_tuple("OwnedInstance")
            .field(&self.raw)
            .finish()
    }
}

impl PartialEq for OwnedInstance {
    fn eq(&self, other: &Self) -> bool {
        self.raw == other.raw
    }
}

impl Drop for OwnedInstance {
    fn drop(&mut self) {
        if self.raw != 0 {
            unsafe { (self.destroy)(self.raw) };
        }
    }
}

/// Internal, owned representation passed from a typed descriptor to the SDK's
/// host adapter.
#[doc(hidden)]
#[derive(Debug, PartialEq)]
pub struct EncodedSpawnable {
    pub(crate) kind: SpawnKind,
    pub(crate) payload: SpawnPayload,
}

impl EncodedSpawnable {
    #[doc(hidden)]
    pub fn custom_owned<T: Send + 'static>(identifier: &'static str, value: T) -> Self {
        let raw = Box::into_raw(Box::new(value)) as usize as u64;
        Self {
            kind: SpawnKind::Custom,
            payload: SpawnPayload::Custom {
                identifier,
                instance: Some(OwnedInstance {
                    raw,
                    destroy: Self::destroy_owned::<T>,
                }),
            },
        }
    }

    #[doc(hidden)]
    pub unsafe fn destroy_owned<T>(raw: u64) {
        if raw != 0 {
            drop(unsafe { Box::from_raw(raw as usize as *mut T) });
        }
    }

    pub(crate) fn with_raw(
        mut self,
        options: &SpawnOptions,
        callback: impl FnOnce(
            &dragonfly_plugin_sys::DfEntitySpawnViewV3,
        ) -> dragonfly_plugin_sys::DfStatus,
    ) -> Option<dragonfly_plugin_sys::DfStatus> {
        let mut callback = Some(callback);
        let mut view = dragonfly_plugin_sys::DfEntitySpawnViewV3 {
            kind: 0,
            flags: 0,
            options: dragonfly_plugin_sys::DfEntitySpawnOptions {
                position: options.position.into(),
                rotation: options.rotation.unwrap_or_default().into(),
                velocity: options.velocity.unwrap_or_default().into(),
                name_tag: crate::string_view_from_str(options.name_tag.as_deref().unwrap_or("")),
            },
            owner: dragonfly_plugin_sys::DfEntityId::default(),
            damage: 0.0,
            fuse_milliseconds: 0,
            experience: 0,
            potion: 0,
            punch_level: 0,
            piercing_level: 0,
            text: dragonfly_plugin_sys::DfStringView::default(),
            item: core::ptr::null(),
            block: core::ptr::null(),
            custom_type: dragonfly_plugin_sys::DfStringView::default(),
            custom_instance: 0,
        };
        let result = match &mut self.payload {
            SpawnPayload::Text(text) => {
                view.kind = dragonfly_plugin_sys::DF_ENTITY_TEXT;
                view.text = crate::string_view_from_str(text);
                Some(callback.take()?(&view))
            }
            SpawnPayload::Lightning {
                damage,
                block_fire,
                entity_fire_duration,
            } => {
                view.kind = dragonfly_plugin_sys::DF_ENTITY_LIGHTNING;
                view.damage = *damage;
                view.flags =
                    u32::from(*block_fire) * dragonfly_plugin_sys::DF_ENTITY_LIGHTNING_BLOCK_FIRE;
                view.fuse_milliseconds = entity_fire_duration.as_millis().try_into().ok()?;
                Some(callback.take()?(&view))
            }
            SpawnPayload::Tnt { fuse } => {
                view.kind = dragonfly_plugin_sys::DF_ENTITY_TNT;
                view.fuse_milliseconds = fuse.as_millis().try_into().ok()?;
                Some(callback.take()?(&view))
            }
            SpawnPayload::ExperienceOrb { experience } => {
                view.kind = dragonfly_plugin_sys::DF_ENTITY_EXPERIENCE_ORB;
                view.experience = (*experience).try_into().ok()?;
                Some(callback.take()?(&view))
            }
            SpawnPayload::DroppedItem {
                stack,
                pickup_delay,
            } => stack
                .with_raw(|item| {
                    view.kind = dragonfly_plugin_sys::DF_ENTITY_ITEM;
                    view.item = item;
                    if let Some(delay) = pickup_delay {
                        view.flags = dragonfly_plugin_sys::DF_ENTITY_ITEM_HAS_PICKUP_DELAY;
                        view.fuse_milliseconds = delay.as_millis().try_into().unwrap_or(u64::MAX);
                    }
                    callback.take().map(|callback| callback(&view))
                })
                .flatten(),
            SpawnPayload::FallingBlock(block) => {
                let properties = block.properties_nbt()?;
                let raw_block = dragonfly_plugin_sys::DfBlockView {
                    identifier: crate::string_view_from_str(block.identifier()),
                    properties_nbt: crate::bytes_view(&properties),
                };
                view.kind = dragonfly_plugin_sys::DF_ENTITY_FALLING_BLOCK;
                view.block = &raw_block;
                Some(callback.take()?(&view))
            }
            SpawnPayload::Arrow(arrow) => {
                view.kind = dragonfly_plugin_sys::DF_ENTITY_ARROW;
                view.owner = arrow.owner.raw_id();
                view.damage = arrow.damage;
                view.potion = u32::from(arrow.tip.id());
                view.punch_level = arrow.punch_level.try_into().ok()?;
                view.piercing_level = arrow.piercing_level.try_into().ok()?;
                if arrow.critical {
                    view.flags |= dragonfly_plugin_sys::DF_ENTITY_ARROW_CRITICAL;
                }
                if arrow.disable_pickup {
                    view.flags |= dragonfly_plugin_sys::DF_ENTITY_ARROW_DISABLE_PICKUP;
                }
                if arrow.obtain_on_pickup {
                    view.flags |= dragonfly_plugin_sys::DF_ENTITY_ARROW_OBTAIN_ON_PICKUP;
                }
                Some(callback.take()?(&view))
            }
            SpawnPayload::Egg(owner) => encode_owner_projectile(
                &mut view,
                dragonfly_plugin_sys::DF_ENTITY_EGG,
                *owner,
                callback.take()?,
            ),
            SpawnPayload::Snowball(owner) => encode_owner_projectile(
                &mut view,
                dragonfly_plugin_sys::DF_ENTITY_SNOWBALL,
                *owner,
                callback.take()?,
            ),
            SpawnPayload::EnderPearl(owner) => encode_owner_projectile(
                &mut view,
                dragonfly_plugin_sys::DF_ENTITY_ENDER_PEARL,
                *owner,
                callback.take()?,
            ),
            SpawnPayload::BottleOfEnchanting(owner) => encode_owner_projectile(
                &mut view,
                dragonfly_plugin_sys::DF_ENTITY_BOTTLE_OF_ENCHANTING,
                *owner,
                callback.take()?,
            ),
            SpawnPayload::SplashPotion { owner, potion } => {
                view.potion = u32::from(potion.id());
                encode_owner_projectile(
                    &mut view,
                    dragonfly_plugin_sys::DF_ENTITY_SPLASH_POTION,
                    *owner,
                    callback.take()?,
                )
            }
            SpawnPayload::LingeringPotion { owner, potion } => {
                view.potion = u32::from(potion.id());
                encode_owner_projectile(
                    &mut view,
                    dragonfly_plugin_sys::DF_ENTITY_LINGERING_POTION,
                    *owner,
                    callback.take()?,
                )
            }
            SpawnPayload::Custom {
                identifier,
                instance,
            } => {
                view.kind = dragonfly_plugin_sys::DF_ENTITY_CUSTOM;
                view.custom_type = crate::string_view_from_str(identifier);
                view.custom_instance = instance.as_ref().map_or(0, |instance| instance.raw);
                Some(callback.take()?(&view))
            }
        };
        if result == Some(dragonfly_plugin_sys::DF_STATUS_OK)
            && let SpawnPayload::Custom {
                instance: Some(instance),
                ..
            } = &mut self.payload
        {
            instance.raw = 0;
        }
        result
    }
}

impl<T: CustomType> private::Sealed for T {}

impl<T: CustomType> Spawnable for T {
    fn encode(self) -> EncodedSpawnable {
        T::encode_spawn(self)
    }
}

fn encode_owner_projectile<R>(
    view: &mut dragonfly_plugin_sys::DfEntitySpawnViewV3,
    kind: u32,
    owner: Entity,
    callback: impl FnOnce(&dragonfly_plugin_sys::DfEntitySpawnViewV3) -> R,
) -> Option<R> {
    view.kind = kind;
    view.owner = owner.raw_id();
    Some(callback(view))
}

impl From<Vec3> for dragonfly_plugin_sys::DfVec3 {
    fn from(value: Vec3) -> Self {
        Self {
            x: value.x,
            y: value.y,
            z: value.z,
        }
    }
}

impl From<Rotation> for dragonfly_plugin_sys::DfRotation {
    fn from(value: Rotation) -> Self {
        Self {
            yaw: value.yaw,
            pitch: value.pitch,
        }
    }
}

impl From<dragonfly_plugin_sys::DfRotation> for Rotation {
    fn from(value: dragonfly_plugin_sys::DfRotation) -> Self {
        Self {
            yaw: value.yaw,
            pitch: value.pitch,
        }
    }
}

macro_rules! impl_spawnable {
    ($type:ty, $kind:ident, $encode:expr) => {
        impl private::Sealed for $type {}

        impl Spawnable for $type {
            fn encode(self) -> EncodedSpawnable {
                EncodedSpawnable {
                    kind: SpawnKind::$kind,
                    payload: ($encode)(&self),
                }
            }
        }
    };
}

impl_spawnable!(Text, Text, |value: &Text| SpawnPayload::Text(
    value.text.clone()
));
impl_spawnable!(Lightning, Lightning, |value: &Lightning| {
    SpawnPayload::Lightning {
        damage: value.damage,
        block_fire: value.block_fire,
        entity_fire_duration: value.entity_fire_duration,
    }
});
impl_spawnable!(TNT, Tnt, |value: &TNT| SpawnPayload::Tnt {
    fuse: value.fuse,
});
impl_spawnable!(ExperienceOrb, ExperienceOrb, |value: &ExperienceOrb| {
    SpawnPayload::ExperienceOrb {
        experience: value.experience,
    }
});
impl_spawnable!(DroppedItem, DroppedItem, |value: &DroppedItem| {
    SpawnPayload::DroppedItem {
        stack: value.stack.clone(),
        pickup_delay: value.pickup_delay,
    }
});
impl_spawnable!(FallingBlock, FallingBlock, |value: &FallingBlock| {
    SpawnPayload::FallingBlock(value.block.clone())
});
impl_spawnable!(Arrow, Arrow, |value: &Arrow| SpawnPayload::Arrow(*value));
impl_spawnable!(Egg, Egg, |value: &Egg| SpawnPayload::Egg(value.owner));
impl_spawnable!(Snowball, Snowball, |value: &Snowball| {
    SpawnPayload::Snowball(value.owner)
});
impl_spawnable!(EnderPearl, EnderPearl, |value: &EnderPearl| {
    SpawnPayload::EnderPearl(value.owner)
});
impl_spawnable!(
    BottleOfEnchanting,
    BottleOfEnchanting,
    |value: &BottleOfEnchanting| SpawnPayload::BottleOfEnchanting(value.owner)
);
impl_spawnable!(SplashPotion, SplashPotion, |value: &SplashPotion| {
    SpawnPayload::SplashPotion {
        owner: value.owner,
        potion: value.potion,
    }
});
impl_spawnable!(
    LingeringPotion,
    LingeringPotion,
    |value: &LingeringPotion| SpawnPayload::LingeringPotion {
        owner: value.owner,
        potion: value.potion
    }
);

#[cfg(test)]
mod tests {
    use super::*;

    #[derive(Default)]
    struct MacroLivingEntity;

    #[dragonfly_plugin_macros::entity(width = 0.6, height = 1.8)]
    impl Living for MacroLivingEntity {
        fn tick(&mut self, _context: &mut TickContext<'_>) {}
        fn hurt(&mut self, _event: &mut Hurt<'_>) {}
    }
    use crate::item;

    #[test]
    fn spawn_options_have_dragonfly_style_builders() {
        let options = SpawnOptions::new(Vec3 {
            x: 1.0,
            y: 64.0,
            z: -2.0,
        })
        .rotation(Rotation {
            yaw: 90.0,
            pitch: 10.0,
        })
        .velocity(Vec3 {
            x: 0.1,
            y: 0.2,
            z: 0.3,
        })
        .name_tag("Example");

        assert_eq!(options.position.y, 64.0);
        assert_eq!(options.rotation.map(|rotation| rotation.yaw), Some(90.0));
        assert_eq!(options.velocity.map(|velocity| velocity.y), Some(0.2));
        assert_eq!(options.name_tag.as_deref(), Some("Example"));
    }

    #[test]
    fn descriptors_encode_owned_payloads() {
        assert_eq!(
            Text::new("floating text").encode().payload,
            SpawnPayload::Text("floating text".to_owned())
        );
        assert_eq!(
            Lightning::new()
                .damage(8.5)
                .block_fire(false)
                .entity_fire_duration(Duration::from_secs(3))
                .encode()
                .payload,
            SpawnPayload::Lightning {
                damage: 8.5,
                block_fire: false,
                entity_fire_duration: Duration::from_secs(3),
            }
        );
        assert_eq!(
            TNT::new(Duration::from_secs(4)).encode().payload,
            SpawnPayload::Tnt {
                fuse: Duration::from_secs(4)
            }
        );
        assert_eq!(
            ExperienceOrb::new(37).encode().payload,
            SpawnPayload::ExperienceOrb { experience: 37 }
        );

        let stack = item::new(item::Diamond, 3);
        assert_eq!(
            DroppedItem::new(stack.clone())
                .pickup_delay(Duration::from_millis(500))
                .encode()
                .payload,
            SpawnPayload::DroppedItem {
                stack,
                pickup_delay: Some(Duration::from_millis(500)),
            }
        );

        let block = block::new(block::Sand);
        assert_eq!(
            FallingBlock::new(block.clone()).encode().payload,
            SpawnPayload::FallingBlock(block)
        );
    }

    #[test]
    fn entity_state_keeps_optional_capabilities_explicit() {
        let state = EntityState {
            entity_type: "dragonfly:text".to_owned(),
            position: Vec3 {
                x: 2.0,
                y: 3.0,
                z: 4.0,
            },
            rotation: Rotation {
                yaw: 20.0,
                pitch: -5.0,
            },
            world: None,
            velocity: None,
            name_tag: Some("Label".to_owned()),
            can_teleport: false,
        };

        assert_eq!(state.name_tag.as_deref(), Some("Label"));
        assert!(state.velocity.is_none());
        assert!(!state.can_teleport);
    }

    #[test]
    fn descriptors_are_spawnable_without_adapter_types() {
        fn assert_spawnable(_: &impl Spawnable) {}

        assert_spawnable(&Text::new("text"));
        assert_spawnable(&Lightning::new());
        assert_spawnable(&TNT::new(Duration::from_secs(4)));
        assert_spawnable(&ExperienceOrb::new(7));
        assert_spawnable(&DroppedItem::new(item::new(item::Apple, 1)));
        assert_spawnable(&FallingBlock::new(block::Sand));
        let owner = Entity::from_id(EntityId::from_parts([7; 16], 9));
        assert_spawnable(&Arrow::new(owner).tip(Potion::Poison).critical(true));
        assert_spawnable(&Egg::new(owner));
        assert_spawnable(&Snowball::new(owner));
        assert_spawnable(&EnderPearl::new(owner));
        assert_spawnable(&BottleOfEnchanting::new(owner));
        assert_spawnable(&SplashPotion::new(owner, Potion::Healing));
        assert_spawnable(&LingeringPotion::new(owner, Potion::Poison));
    }

    #[test]
    fn saved_state_round_trips_typed_values() {
        let mut written = SavedStateMut::new();
        written.set("hits", 42_i64);
        written.set("name", "dummy");
        written.set("active", true);
        let bytes = written.encode().unwrap();
        let read = SavedState::decode(7, &bytes).unwrap();
        assert_eq!(read.version(), 7);
        assert_eq!(read.i64("hits"), Some(42));
        assert_eq!(read.str("name"), Some("dummy"));
        assert_eq!(read.bool("active"), Some(true));
        assert_eq!(read.i64("missing"), None);
    }

    #[test]
    fn entity_context_cancellation_is_monotonic() {
        let input = dragonfly_plugin_sys::DfEntityHurtInput {
            entity: dragonfly_plugin_sys::DfEntityId {
                bytes: [4; 16],
                generation: 9,
            },
            health: 20.0,
            max_health: 40.0,
            ..Default::default()
        };
        let mut state = dragonfly_plugin_sys::DfEntityHurtState {
            damage: 6.0,
            cancelled: 0,
        };
        let mut event = unsafe { Hurt::from_raw(&input, &mut state) };
        assert_eq!(event.entity().id().generation(), 9);
        assert_eq!(event.health(), 20.0);
        assert_eq!(event.max_health(), 40.0);
        assert_eq!(event.damage(), 6.0);
        event.set_damage(-1.0);
        assert_eq!(event.damage(), 0.0);
        event.cancel();
        event.set_damage(2.0);
        assert!(event.cancelled());
    }

    #[derive(Default)]
    struct TestLiving;

    impl Living for TestLiving {}

    #[test]
    fn living_defaults_match_dragonfly() {
        assert_eq!(TestLiving::INITIAL_HEALTH, 20.0);
        assert_eq!(TestLiving::MAX_HEALTH, 20.0);
        assert_eq!(TestLiving::SPEED, 0.1);
        assert_eq!(TestLiving::PHYSICS, Some(Physics::passive(0.08, 0.02)));
    }

    #[test]
    fn advanced_entity_macro_registers_owned_living_state() {
        let registered = REGISTERED_TYPES
            .iter()
            .find(|registered| {
                registered.descriptor().family == dragonfly_plugin_sys::DF_ENTITY_FAMILY_LIVING
            })
            .expect("advanced entity registration");
        assert_eq!(
            registered.descriptor().family,
            dragonfly_plugin_sys::DF_ENTITY_FAMILY_LIVING
        );
        assert!(registered.handler().is_some());
    }

    #[test]
    fn custom_spawn_instance_drops_on_failure_and_transfers_on_success() {
        use std::sync::{
            Arc,
            atomic::{AtomicUsize, Ordering},
        };

        #[derive(Clone)]
        struct DropCount(Arc<AtomicUsize>);
        impl Drop for DropCount {
            fn drop(&mut self) {
                self.0.fetch_add(1, Ordering::SeqCst);
            }
        }

        let failed = Arc::new(AtomicUsize::new(0));
        let encoded = EncodedSpawnable::custom_owned("example:failed", DropCount(failed.clone()));
        assert_eq!(
            encoded.with_raw(&SpawnOptions::new(Vec3::default()), |_| {
                dragonfly_plugin_sys::DF_STATUS_ERROR
            }),
            Some(dragonfly_plugin_sys::DF_STATUS_ERROR)
        );
        assert_eq!(failed.load(Ordering::SeqCst), 1);

        let adopted = Arc::new(AtomicUsize::new(0));
        let encoded = EncodedSpawnable::custom_owned("example:ok", DropCount(adopted.clone()));
        let mut raw = 0;
        assert_eq!(
            encoded.with_raw(&SpawnOptions::new(Vec3::default()), |view| {
                raw = view.custom_instance;
                dragonfly_plugin_sys::DF_STATUS_OK
            }),
            Some(dragonfly_plugin_sys::DF_STATUS_OK)
        );
        assert_eq!(adopted.load(Ordering::SeqCst), 0);
        unsafe { EncodedSpawnable::destroy_owned::<DropCount>(raw) };
        assert_eq!(adopted.load(Ordering::SeqCst), 1);
    }
}
