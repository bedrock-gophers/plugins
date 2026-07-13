use std::collections::BTreeMap;

use crate::Value;

#[derive(Clone, Debug, Eq, PartialEq)]
pub enum Property {
    Bool(bool),
    Uint8(u8),
    Int32(i32),
    String(String),
}

impl From<bool> for Property {
    fn from(value: bool) -> Self {
        Self::Bool(value)
    }
}

impl From<u8> for Property {
    fn from(value: u8) -> Self {
        Self::Uint8(value)
    }
}

impl From<i32> for Property {
    fn from(value: i32) -> Self {
        Self::Int32(value)
    }
}

impl From<String> for Property {
    fn from(value: String) -> Self {
        Self::String(value)
    }
}

impl From<&str> for Property {
    fn from(value: &str) -> Self {
        Self::String(value.to_owned())
    }
}

/// An owned block state returned by the host.
///
/// Built-in states are changed through their generated typed values. Raw state mutation is kept
/// on [`Custom`] so a built-in block cannot silently bypass its generated state types.
///
/// ```compile_fail
/// use dragonfly::block;
///
/// let _ = block::new(block::Sand).with_property("invented", 1i32);
/// ```
#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Block {
    identifier: String,
    properties: BTreeMap<String, Property>,
}

impl Block {
    fn new(identifier: impl Into<String>) -> Self {
        Self {
            identifier: identifier.into(),
            properties: BTreeMap::new(),
        }
    }

    fn from_properties<const N: usize>(
        identifier: impl Into<String>,
        properties: [(&str, Property); N],
    ) -> Self {
        Self {
            identifier: identifier.into(),
            properties: properties
                .into_iter()
                .map(|(name, value)| (name.to_owned(), value))
                .collect(),
        }
    }

    pub fn air() -> Self {
        Self::new("minecraft:air")
    }

    pub fn identifier(&self) -> &str {
        &self.identifier
    }

    pub fn properties(&self) -> &BTreeMap<String, Property> {
        &self.properties
    }

    pub fn property(&self, name: &str) -> Option<&Property> {
        self.properties.get(name)
    }

    pub(crate) fn with_property(
        mut self,
        name: impl Into<String>,
        value: impl Into<Property>,
    ) -> Self {
        self.set_property(name, value);
        self
    }

    pub(crate) fn set_property(&mut self, name: impl Into<String>, value: impl Into<Property>) {
        self.properties.insert(name.into(), value.into());
    }

    pub(crate) fn properties_nbt(&self) -> Option<Vec<u8>> {
        let values = self
            .properties
            .iter()
            .map(|(name, value)| {
                let (kind, value) = match value {
                    Property::Bool(value) => (0, Value::Byte(i8::from(*value))),
                    Property::Uint8(value) => (1, Value::Byte(i8::from_le_bytes([*value]))),
                    Property::Int32(value) => (2, Value::Int(*value)),
                    Property::String(value) => (3, Value::String(value.clone())),
                };
                (
                    name.clone(),
                    Value::Compound(BTreeMap::from([
                        ("kind".to_owned(), Value::Int(kind)),
                        ("value".to_owned(), value),
                    ])),
                )
            })
            .collect();
        crate::item_nbt::encode_values(&values).ok()
    }

    pub(crate) fn from_nbt(identifier: String, bytes: &[u8]) -> Option<Self> {
        let values = crate::item_nbt::decode_values(bytes).ok()?;
        let mut properties = BTreeMap::new();
        for (name, value) in values {
            let Value::Compound(mut tagged) = value else {
                return None;
            };
            let Some(Value::Int(kind)) = tagged.remove("kind") else {
                return None;
            };
            let value = tagged.remove("value")?;
            if !tagged.is_empty() {
                return None;
            }
            let property = match (kind, value) {
                (0, Value::Byte(value)) if value == 0 || value == 1 => Property::Bool(value != 0),
                (1, Value::Byte(value)) => Property::Uint8(value.to_le_bytes()[0]),
                (2, Value::Int(value)) => Property::Int32(value),
                (3, Value::String(value)) => Property::String(value),
                _ => return None,
            };
            properties.insert(name, property);
        }
        Some(Self {
            identifier,
            properties,
        })
    }
}

impl Default for Block {
    fn default() -> Self {
        Self::air()
    }
}

impl From<&Block> for Block {
    fn from(value: &Block) -> Self {
        value.clone()
    }
}

#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Custom {
    block: Block,
}

impl Custom {
    pub fn new(identifier: impl Into<String>) -> Self {
        Self {
            block: Block::new(identifier),
        }
    }

    pub fn with_property(mut self, name: impl Into<String>, value: impl Into<Property>) -> Self {
        self.block = self.block.with_property(name, value);
        self
    }
}

impl From<Custom> for Block {
    fn from(value: Custom) -> Self {
        value.block
    }
}

pub fn new(block: impl Into<Block>) -> Block {
    block.into()
}

include!("blocks_generated.rs");

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn typed_blocks_convert_to_owned_states() {
        let sand = new(Sand);
        assert_eq!(sand.identifier(), "minecraft:sand");
        assert!(sand.properties().is_empty());

        let log = new(OakLog::new(PillarAxis::Z));
        assert_eq!(log.identifier(), "minecraft:oak_log");
        assert_eq!(
            log.property("pillar_axis"),
            Some(&Property::String("z".to_owned()))
        );
    }

    #[test]
    fn custom_is_the_identifier_escape_hatch() {
        let custom = new(Custom::new("example:machine").with_property("active", true));
        assert_eq!(custom.identifier(), "example:machine");
        assert_eq!(custom.property("active"), Some(&Property::Bool(true)));
    }

    #[test]
    fn block_properties_round_trip() {
        let block = new(Custom::new("minecraft:oak_door"))
            .with_property("open_bit", true)
            .with_property("direction", 2i32)
            .with_property("age", 7u8)
            .with_property("pillar_axis", "y");
        let decoded =
            Block::from_nbt(block.identifier.clone(), &block.properties_nbt().unwrap()).unwrap();
        assert_eq!(decoded, block);
    }
}
