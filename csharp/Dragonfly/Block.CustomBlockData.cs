using System.Collections;
using System.Text.Json;

namespace Dragonfly;

public static partial class Block
{
    public sealed class CustomBlockData
    {
        private readonly Dictionary<string, object> _server = new(StringComparer.Ordinal);
        private readonly Dictionary<string, object[]> _states = new(StringComparer.Ordinal);
        private readonly List<PermutationData> _permutations = [];

        public CustomBlockData SetNativeProperty(string name, object value) { Set(_server, name, value); return this; }

        public CustomBlockData AddState(string name, params object[] values)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(name);
            ArgumentNullException.ThrowIfNull(values);
            if (values.Length == 0) throw new ArgumentException("custom block states require at least one value", nameof(values));
            if (values.Any(value => value is not bool and not byte and not int and not string))
                throw new ArgumentException("custom block state values must be bool, byte, int or string", nameof(values));
            _states[name] = values.ToArray();
            return this;
        }

        public CustomBlockData AddPermutation(string condition, IReadOnlyDictionary<string, object> properties)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(condition);
            ArgumentNullException.ThrowIfNull(properties);
            _permutations.Add(new PermutationData(condition, properties));
            return this;
        }

        internal byte[] ToJson()
        {
            using var buffer = new MemoryStream();
            using (var writer = new Utf8JsonWriter(buffer))
            {
                writer.WriteStartObject();
                writer.WritePropertyName("server"); WriteDictionary(writer, _server);
                writer.WritePropertyName("states");
                writer.WriteStartObject();
                foreach (var (name, values) in _states) { writer.WritePropertyName(name); WriteValue(writer, values); }
                writer.WriteEndObject();
                writer.WritePropertyName("permutations");
                writer.WriteStartArray();
                foreach (var permutation in _permutations)
                {
                    writer.WriteStartObject();
                    writer.WriteString("condition", permutation.Condition);
                    writer.WritePropertyName("properties"); WriteDictionary(writer, permutation.Properties);
                    writer.WriteEndObject();
                }
                writer.WriteEndArray();
                writer.WriteEndObject();
            }
            return buffer.ToArray();
        }

        internal IReadOnlyDictionary<string, object?> DefaultState() =>
            _states.ToDictionary(entry => entry.Key, entry => (object?)entry.Value[0], StringComparer.Ordinal);

        private static void Set(IDictionary<string, object> target, string name, object value)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(name);
            ArgumentNullException.ThrowIfNull(value);
            if (name.Length > 128) throw new ArgumentOutOfRangeException(nameof(name));
            target[name] = value;
        }

        private static void WriteDictionary(Utf8JsonWriter writer, IEnumerable<KeyValuePair<string, object>> values)
        {
            writer.WriteStartObject();
            foreach (var (key, value) in values) { writer.WritePropertyName(key); WriteValue(writer, value); }
            writer.WriteEndObject();
        }

        private static void WriteValue(Utf8JsonWriter writer, object value)
        {
            switch (value)
            {
                case string v: writer.WriteStringValue(v); break;
                case bool v: writer.WriteBooleanValue(v); break;
                case byte v: writer.WriteNumberValue(v); break;
                case sbyte v: writer.WriteNumberValue(v); break;
                case short v: writer.WriteNumberValue(v); break;
                case ushort v: writer.WriteNumberValue(v); break;
                case int v: writer.WriteNumberValue(v); break;
                case uint v: writer.WriteNumberValue(v); break;
                case long v: writer.WriteNumberValue(v); break;
                case ulong v: writer.WriteNumberValue(v); break;
                case float v when float.IsFinite(v): writer.WriteNumberValue(v); break;
                case double v when double.IsFinite(v): writer.WriteNumberValue(v); break;
                case decimal v: writer.WriteNumberValue(v); break;
                case IReadOnlyDictionary<string, object> v: WriteDictionary(writer, v); break;
                case IDictionary v:
                    writer.WriteStartObject();
                    foreach (DictionaryEntry entry in v)
                    {
                        if (entry.Key is not string key || entry.Value is null) throw new ArgumentException("custom block dictionaries require string keys and non-null values");
                        writer.WritePropertyName(key); WriteValue(writer, entry.Value);
                    }
                    writer.WriteEndObject(); break;
                case IEnumerable v:
                    writer.WriteStartArray();
                    foreach (var entry in v) WriteValue(writer, entry ?? throw new ArgumentException("custom block arrays cannot contain null values"));
                    writer.WriteEndArray(); break;
                default: throw new ArgumentException($"unsupported custom block property value type {value.GetType().FullName}");
            }
        }

        private sealed record PermutationData(string Condition, IReadOnlyDictionary<string, object> Properties);
    }
}
