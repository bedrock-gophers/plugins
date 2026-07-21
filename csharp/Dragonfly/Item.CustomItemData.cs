using System.Collections;
using System.Text.Json;

namespace Dragonfly;

public static partial class Item
{
    public sealed class CustomItemData
    {
        private readonly Dictionary<string, object> _properties = new(StringComparer.Ordinal);
        private readonly Dictionary<string, object> _components = new(StringComparer.Ordinal);
        private readonly Dictionary<string, object> _server = new(StringComparer.Ordinal);

        public CustomItemData SetProperty(string name, object value)
        {
            ValidateName(name, nameof(name));
            ArgumentNullException.ThrowIfNull(value);
            _properties[name] = value;
            return this;
        }

        public CustomItemData AddComponent(string name, IReadOnlyDictionary<string, object>? values = null)
        {
            ValidateName(name, nameof(name));
            _components[name] = values ?? new Dictionary<string, object>();
            return this;
        }

        public CustomItemData RemoveProperty(string name)
        {
            ArgumentNullException.ThrowIfNull(name);
            _properties.Remove(name);
            return this;
        }

        public CustomItemData RemoveComponent(string name)
        {
            ArgumentNullException.ThrowIfNull(name);
            _components.Remove(name);
            return this;
        }

        public CustomItemData SetServerBehavior(string name, object value)
        {
            ValidateName(name, nameof(name));
            ArgumentNullException.ThrowIfNull(value);
            _server[name] = value;
            return this;
        }

        internal byte[] ToJson()
        {
            using var buffer = new MemoryStream();
            using (var writer = new Utf8JsonWriter(buffer))
            {
                writer.WriteStartObject();
                writer.WritePropertyName("properties");
                WriteDictionary(writer, _properties);
                writer.WritePropertyName("components");
                WriteDictionary(writer, _components);
                writer.WritePropertyName("server");
                WriteDictionary(writer, _server);
                writer.WriteEndObject();
            }
            return buffer.ToArray();
        }

        private static void WriteValue(Utf8JsonWriter writer, object value)
        {
            switch (value)
            {
                case string text: writer.WriteStringValue(text); return;
                case bool boolean: writer.WriteBooleanValue(boolean); return;
                case byte number: writer.WriteNumberValue(number); return;
                case sbyte number: writer.WriteNumberValue(number); return;
                case short number: writer.WriteNumberValue(number); return;
                case ushort number: writer.WriteNumberValue(number); return;
                case int number: writer.WriteNumberValue(number); return;
                case uint number: writer.WriteNumberValue(number); return;
                case long number: writer.WriteNumberValue(number); return;
                case ulong number: writer.WriteNumberValue(number); return;
                case float number when float.IsFinite(number): writer.WriteNumberValue(number); return;
                case double number when double.IsFinite(number): writer.WriteNumberValue(number); return;
                case decimal number: writer.WriteNumberValue(number); return;
                case IReadOnlyDictionary<string, object> dictionary:
                    WriteDictionary(writer, dictionary);
                    return;
                case IDictionary dictionary:
                    writer.WriteStartObject();
                    foreach (DictionaryEntry entry in dictionary)
                    {
                        if (entry.Key is not string key || entry.Value is null)
                            throw new ArgumentException("custom item dictionaries require string keys and non-null values");
                        writer.WritePropertyName(key);
                        WriteValue(writer, entry.Value);
                    }
                    writer.WriteEndObject();
                    return;
                case IEnumerable values:
                    writer.WriteStartArray();
                    foreach (var entry in values)
                    {
                        if (entry is null) throw new ArgumentException("custom item arrays cannot contain null values");
                        WriteValue(writer, entry);
                    }
                    writer.WriteEndArray();
                    return;
                default:
                    throw new ArgumentException($"unsupported custom item component value type {value.GetType().FullName}");
            }
        }

        private static void WriteDictionary(Utf8JsonWriter writer, IEnumerable<KeyValuePair<string, object>> values)
        {
            writer.WriteStartObject();
            foreach (var (key, value) in values)
            {
                writer.WritePropertyName(key);
                WriteValue(writer, value);
            }
            writer.WriteEndObject();
        }

        private static void ValidateName(string name, string parameter)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(name, parameter);
            if (name.Length > 128) throw new ArgumentOutOfRangeException(parameter);
        }
    }
}
