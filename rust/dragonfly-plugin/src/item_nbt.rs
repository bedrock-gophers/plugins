#![allow(dead_code)]

use std::collections::BTreeMap;

use crate::Value;

const TAG_END: u8 = 0;
const TAG_BYTE: u8 = 1;
const TAG_SHORT: u8 = 2;
const TAG_INT: u8 = 3;
const TAG_LONG: u8 = 4;
const TAG_FLOAT: u8 = 5;
const TAG_DOUBLE: u8 = 6;
const TAG_BYTE_ARRAY: u8 = 7;
const TAG_STRING: u8 = 8;
const TAG_LIST: u8 = 9;
const TAG_COMPOUND: u8 = 10;
const TAG_INT_ARRAY: u8 = 11;
const TAG_LONG_ARRAY: u8 = 12;

const MAX_BYTES: usize = 16 << 20;
const MAX_DEPTH: usize = 64;
const MAX_ELEMENTS: usize = 1 << 20;

#[derive(Clone, Debug, Eq, PartialEq)]
pub(crate) enum NbtError {
    TooLarge,
    TooDeep,
    TooManyElements,
    StringTooLong,
    MixedList,
    InvalidTag(u8),
    InvalidLength,
    UnexpectedEof,
    InvalidUtf8,
    DuplicateKey,
    RootNotCompound,
    TrailingData,
}

pub(crate) fn encode_values(values: &BTreeMap<String, Value>) -> Result<Vec<u8>, NbtError> {
    let mut writer = Writer::default();
    writer.byte(TAG_COMPOUND)?;
    writer.string("")?;
    writer.compound(values, 0)?;
    Ok(writer.bytes)
}

pub(crate) fn decode_values(bytes: &[u8]) -> Result<BTreeMap<String, Value>, NbtError> {
    if bytes.len() > MAX_BYTES {
        return Err(NbtError::TooLarge);
    }
    let mut reader = Reader {
        bytes,
        offset: 0,
        nodes: 0,
    };
    if reader.byte()? != TAG_COMPOUND {
        return Err(NbtError::RootNotCompound);
    }
    reader.string()?;
    let values = reader.compound(0)?;
    if reader.offset != bytes.len() {
        return Err(NbtError::TrailingData);
    }
    Ok(values)
}

#[derive(Default)]
struct Writer {
    bytes: Vec<u8>,
    nodes: usize,
}

impl Writer {
    fn append(&mut self, bytes: &[u8]) -> Result<(), NbtError> {
        let new_len = self
            .bytes
            .len()
            .checked_add(bytes.len())
            .ok_or(NbtError::TooLarge)?;
        if new_len > MAX_BYTES {
            return Err(NbtError::TooLarge);
        }
        self.bytes.extend_from_slice(bytes);
        Ok(())
    }

    fn byte(&mut self, value: u8) -> Result<(), NbtError> {
        self.append(&[value])
    }

    fn i16(&mut self, value: i16) -> Result<(), NbtError> {
        self.append(&value.to_le_bytes())
    }

    fn u16(&mut self, value: u16) -> Result<(), NbtError> {
        self.append(&value.to_le_bytes())
    }

    fn i32(&mut self, value: i32) -> Result<(), NbtError> {
        self.append(&value.to_le_bytes())
    }

    fn i64(&mut self, value: i64) -> Result<(), NbtError> {
        self.append(&value.to_le_bytes())
    }

    fn string(&mut self, value: &str) -> Result<(), NbtError> {
        let len = u16::try_from(value.len()).map_err(|_| NbtError::StringTooLong)?;
        self.u16(len)?;
        self.append(value.as_bytes())
    }

    fn count(&mut self, len: usize) -> Result<(), NbtError> {
        if len > MAX_ELEMENTS {
            return Err(NbtError::TooManyElements);
        }
        let len = i32::try_from(len).map_err(|_| NbtError::InvalidLength)?;
        self.i32(len)
    }

    fn node(&mut self) -> Result<(), NbtError> {
        self.nodes = self.nodes.checked_add(1).ok_or(NbtError::TooManyElements)?;
        if self.nodes > MAX_ELEMENTS {
            return Err(NbtError::TooManyElements);
        }
        Ok(())
    }

    fn compound(&mut self, values: &BTreeMap<String, Value>, depth: usize) -> Result<(), NbtError> {
        check_depth(depth)?;
        if values.len() > MAX_ELEMENTS {
            return Err(NbtError::TooManyElements);
        }
        for (name, value) in values {
            self.node()?;
            self.byte(tag(value))?;
            self.string(name)?;
            self.payload(value, depth + 1)?;
        }
        self.byte(TAG_END)
    }

    fn payload(&mut self, value: &Value, depth: usize) -> Result<(), NbtError> {
        check_depth(depth)?;
        match value {
            Value::Byte(value) => self.byte(value.to_le_bytes()[0]),
            Value::Short(value) => self.i16(*value),
            Value::Int(value) => self.i32(*value),
            Value::Long(value) => self.i64(*value),
            Value::Float(value) => self.append(&value.to_le_bytes()),
            Value::Double(value) => self.append(&value.to_le_bytes()),
            Value::String(value) => self.string(value),
            Value::ByteArray(values) => {
                self.count(values.len())?;
                self.append(values)
            }
            Value::IntArray(values) => {
                self.count(values.len())?;
                for value in values {
                    self.i32(*value)?;
                }
                Ok(())
            }
            Value::LongArray(values) => {
                self.count(values.len())?;
                for value in values {
                    self.i64(*value)?;
                }
                Ok(())
            }
            Value::List(values) => self.list(values, depth),
            Value::Compound(values) => self.compound(values, depth),
        }
    }

    fn list(&mut self, values: &[Value], depth: usize) -> Result<(), NbtError> {
        check_depth(depth)?;
        let element_tag = values.first().map_or(TAG_END, tag);
        if values.iter().any(|value| tag(value) != element_tag) {
            return Err(NbtError::MixedList);
        }
        self.byte(element_tag)?;
        self.count(values.len())?;
        for value in values {
            self.node()?;
            self.payload(value, depth + 1)?;
        }
        Ok(())
    }
}

struct Reader<'a> {
    bytes: &'a [u8],
    offset: usize,
    nodes: usize,
}

impl Reader<'_> {
    fn take(&mut self, len: usize) -> Result<&[u8], NbtError> {
        let end = self
            .offset
            .checked_add(len)
            .ok_or(NbtError::UnexpectedEof)?;
        let value = self
            .bytes
            .get(self.offset..end)
            .ok_or(NbtError::UnexpectedEof)?;
        self.offset = end;
        Ok(value)
    }

    fn byte(&mut self) -> Result<u8, NbtError> {
        Ok(self.take(1)?[0])
    }

    fn i16(&mut self) -> Result<i16, NbtError> {
        let bytes: [u8; 2] = self
            .take(2)?
            .try_into()
            .map_err(|_| NbtError::UnexpectedEof)?;
        Ok(i16::from_le_bytes(bytes))
    }

    fn u16(&mut self) -> Result<u16, NbtError> {
        let bytes: [u8; 2] = self
            .take(2)?
            .try_into()
            .map_err(|_| NbtError::UnexpectedEof)?;
        Ok(u16::from_le_bytes(bytes))
    }

    fn i32(&mut self) -> Result<i32, NbtError> {
        let bytes: [u8; 4] = self
            .take(4)?
            .try_into()
            .map_err(|_| NbtError::UnexpectedEof)?;
        Ok(i32::from_le_bytes(bytes))
    }

    fn i64(&mut self) -> Result<i64, NbtError> {
        let bytes: [u8; 8] = self
            .take(8)?
            .try_into()
            .map_err(|_| NbtError::UnexpectedEof)?;
        Ok(i64::from_le_bytes(bytes))
    }

    fn f32(&mut self) -> Result<f32, NbtError> {
        let bytes: [u8; 4] = self
            .take(4)?
            .try_into()
            .map_err(|_| NbtError::UnexpectedEof)?;
        Ok(f32::from_le_bytes(bytes))
    }

    fn f64(&mut self) -> Result<f64, NbtError> {
        let bytes: [u8; 8] = self
            .take(8)?
            .try_into()
            .map_err(|_| NbtError::UnexpectedEof)?;
        Ok(f64::from_le_bytes(bytes))
    }

    fn string(&mut self) -> Result<String, NbtError> {
        let len = usize::from(self.u16()?);
        let bytes = self.take(len)?;
        let value = std::str::from_utf8(bytes).map_err(|_| NbtError::InvalidUtf8)?;
        Ok(value.to_owned())
    }

    fn count(&mut self) -> Result<usize, NbtError> {
        let value = self.i32()?;
        let value = usize::try_from(value).map_err(|_| NbtError::InvalidLength)?;
        if value > MAX_ELEMENTS {
            return Err(NbtError::TooManyElements);
        }
        Ok(value)
    }

    fn node(&mut self) -> Result<(), NbtError> {
        self.nodes = self.nodes.checked_add(1).ok_or(NbtError::TooManyElements)?;
        if self.nodes > MAX_ELEMENTS {
            return Err(NbtError::TooManyElements);
        }
        Ok(())
    }

    fn compound(&mut self, depth: usize) -> Result<BTreeMap<String, Value>, NbtError> {
        check_depth(depth)?;
        let mut values = BTreeMap::new();
        loop {
            let tag = self.byte()?;
            if tag == TAG_END {
                return Ok(values);
            }
            validate_tag(tag)?;
            self.node()?;
            let name = self.string()?;
            let value = self.payload(tag, depth + 1)?;
            if values.insert(name, value).is_some() {
                return Err(NbtError::DuplicateKey);
            }
        }
    }

    fn payload(&mut self, tag: u8, depth: usize) -> Result<Value, NbtError> {
        check_depth(depth)?;
        match tag {
            TAG_BYTE => Ok(Value::Byte(i8::from_le_bytes([self.byte()?]))),
            TAG_SHORT => Ok(Value::Short(self.i16()?)),
            TAG_INT => Ok(Value::Int(self.i32()?)),
            TAG_LONG => Ok(Value::Long(self.i64()?)),
            TAG_FLOAT => Ok(Value::Float(self.f32()?)),
            TAG_DOUBLE => Ok(Value::Double(self.f64()?)),
            TAG_BYTE_ARRAY => {
                let len = self.count()?;
                Ok(Value::ByteArray(self.take(len)?.to_vec()))
            }
            TAG_STRING => Ok(Value::String(self.string()?)),
            TAG_LIST => self.list(depth),
            TAG_COMPOUND => Ok(Value::Compound(self.compound(depth)?)),
            TAG_INT_ARRAY => {
                let len = self.count()?;
                let mut values = Vec::with_capacity(len);
                for _ in 0..len {
                    values.push(self.i32()?);
                }
                Ok(Value::IntArray(values))
            }
            TAG_LONG_ARRAY => {
                let len = self.count()?;
                let mut values = Vec::with_capacity(len);
                for _ in 0..len {
                    values.push(self.i64()?);
                }
                Ok(Value::LongArray(values))
            }
            other => Err(NbtError::InvalidTag(other)),
        }
    }

    fn list(&mut self, depth: usize) -> Result<Value, NbtError> {
        check_depth(depth)?;
        let element_tag = self.byte()?;
        let len = self.count()?;
        if element_tag == TAG_END {
            if len == 0 {
                return Ok(Value::List(Vec::new()));
            }
            return Err(NbtError::InvalidTag(element_tag));
        }
        validate_tag(element_tag)?;
        let mut values = Vec::with_capacity(len);
        for _ in 0..len {
            self.node()?;
            values.push(self.payload(element_tag, depth + 1)?);
        }
        Ok(Value::List(values))
    }
}

fn tag(value: &Value) -> u8 {
    match value {
        Value::Byte(_) => TAG_BYTE,
        Value::Short(_) => TAG_SHORT,
        Value::Int(_) => TAG_INT,
        Value::Long(_) => TAG_LONG,
        Value::Float(_) => TAG_FLOAT,
        Value::Double(_) => TAG_DOUBLE,
        Value::ByteArray(_) => TAG_BYTE_ARRAY,
        Value::String(_) => TAG_STRING,
        Value::List(_) => TAG_LIST,
        Value::Compound(_) => TAG_COMPOUND,
        Value::IntArray(_) => TAG_INT_ARRAY,
        Value::LongArray(_) => TAG_LONG_ARRAY,
    }
}

fn validate_tag(tag: u8) -> Result<(), NbtError> {
    if (TAG_BYTE..=TAG_LONG_ARRAY).contains(&tag) {
        Ok(())
    } else {
        Err(NbtError::InvalidTag(tag))
    }
}

fn check_depth(depth: usize) -> Result<(), NbtError> {
    if depth > MAX_DEPTH {
        Err(NbtError::TooDeep)
    } else {
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn round_trips_every_nbt_tag_little_endian() {
        let nested = BTreeMap::from([("answer".to_owned(), Value::Int(42))]);
        let values = BTreeMap::from([
            ("byte".to_owned(), Value::Byte(-4)),
            ("short".to_owned(), Value::Short(-300)),
            ("int".to_owned(), Value::Int(0x0102_0304)),
            ("long".to_owned(), Value::Long(-9_000_000_000)),
            ("float".to_owned(), Value::Float(1.25)),
            ("double".to_owned(), Value::Double(-8.5)),
            ("string".to_owned(), Value::String("Bedrock".to_owned())),
            ("bytes".to_owned(), Value::ByteArray(vec![0, 127, 255])),
            ("ints".to_owned(), Value::IntArray(vec![-1, 2, 3])),
            ("longs".to_owned(), Value::LongArray(vec![-4, 5, 6])),
            (
                "list".to_owned(),
                Value::List(vec![
                    Value::String("a".to_owned()),
                    Value::String("b".to_owned()),
                ]),
            ),
            ("compound".to_owned(), Value::Compound(nested)),
        ]);

        let encoded = encode_values(&values).unwrap();
        assert_eq!(&encoded[..3], &[TAG_COMPOUND, 0, 0]);
        assert_eq!(decode_values(&encoded).unwrap(), values);
    }

    #[test]
    fn bool_converts_to_the_standard_byte_tag() {
        let values = BTreeMap::from([("enabled".to_owned(), Value::from(true))]);
        let encoded = encode_values(&values).unwrap();
        assert_eq!(encoded[3], TAG_BYTE);
        assert_eq!(
            decode_values(&encoded).unwrap().get("enabled"),
            Some(&Value::Byte(1))
        );
    }

    #[test]
    fn integers_are_fixed_little_endian() {
        let values = BTreeMap::from([("x".to_owned(), Value::Int(0x0102_0304))]);
        assert_eq!(
            encode_values(&values).unwrap(),
            [
                TAG_COMPOUND,
                0,
                0,
                TAG_INT,
                1,
                0,
                b'x',
                0x04,
                0x03,
                0x02,
                0x01,
                TAG_END,
            ]
        );
    }

    #[test]
    fn rejects_mixed_lists() {
        let values = BTreeMap::from([(
            "mixed".to_owned(),
            Value::List(vec![Value::Int(1), Value::String("two".to_owned())]),
        )]);
        assert_eq!(encode_values(&values), Err(NbtError::MixedList));
    }

    #[test]
    fn rejects_malformed_and_trailing_input() {
        assert_eq!(decode_values(&[]), Err(NbtError::UnexpectedEof));
        assert_eq!(
            decode_values(&[TAG_INT, 0, 0]),
            Err(NbtError::RootNotCompound)
        );
        assert_eq!(
            decode_values(&[TAG_COMPOUND, 0, 0, TAG_STRING, 1, 0, b'x', 4, 0]),
            Err(NbtError::UnexpectedEof)
        );
        let mut valid = encode_values(&BTreeMap::new()).unwrap();
        valid.push(0);
        assert_eq!(decode_values(&valid), Err(NbtError::TrailingData));
    }

    #[test]
    fn rejects_deep_and_oversized_data() {
        let mut value = Value::Int(1);
        for _ in 0..=MAX_DEPTH {
            value = Value::List(vec![value]);
        }
        let values = BTreeMap::from([("deep".to_owned(), value)]);
        assert_eq!(encode_values(&values), Err(NbtError::TooDeep));

        let mut deep_input = vec![TAG_COMPOUND, 0, 0];
        for _ in 0..=MAX_DEPTH {
            deep_input.extend_from_slice(&[TAG_COMPOUND, 1, 0, b'x']);
        }
        deep_input.extend(std::iter::repeat_n(TAG_END, MAX_DEPTH + 2));
        assert_eq!(decode_values(&deep_input), Err(NbtError::TooDeep));

        assert_eq!(
            decode_values(&vec![0; MAX_BYTES + 1]),
            Err(NbtError::TooLarge)
        );
        let oversized =
            BTreeMap::from([("bytes".to_owned(), Value::ByteArray(vec![0; MAX_BYTES]))]);
        assert_eq!(encode_values(&oversized), Err(NbtError::TooManyElements));
    }

    #[test]
    fn rejects_negative_lengths_and_duplicate_keys() {
        let negative_list = [
            TAG_COMPOUND,
            0,
            0,
            TAG_LIST,
            1,
            0,
            b'x',
            TAG_INT,
            0xff,
            0xff,
            0xff,
            0xff,
        ];
        assert_eq!(decode_values(&negative_list), Err(NbtError::InvalidLength));

        let duplicate = [
            TAG_COMPOUND,
            0,
            0,
            TAG_BYTE,
            1,
            0,
            b'x',
            1,
            TAG_BYTE,
            1,
            0,
            b'x',
            2,
            TAG_END,
        ];
        assert_eq!(decode_values(&duplicate), Err(NbtError::DuplicateKey));
    }
}
