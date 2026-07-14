using System.Globalization;
using System.Diagnostics.CodeAnalysis;
using System.Reflection;
using Dragonfly.Native;

namespace Dragonfly;

internal sealed class CommandBinding
{
    private CommandBinding(Cmd.Runnable runnable, CommandField[] fields)
    {
        Runnable = runnable;
        Fields = fields;
    }

    internal object Gate { get; } = new();
    internal Cmd.Runnable Runnable { get; }
    internal CommandField[] Fields { get; }

    [UnconditionalSuppressMessage("Trimming", "IL2075", Justification = "The plugin generator roots runnable fields for NativeAOT.")]
    internal static CommandBinding Create(Cmd.Runnable runnable)
    {
        ArgumentNullException.ThrowIfNull(runnable);
        var fields = runnable.GetType()
            .GetFields(BindingFlags.Instance | BindingFlags.Public)
            .Select(field => CommandField.Create(field, runnable))
            .Where(field => field is not null)
            .Cast<CommandField>()
            .ToArray();
        var optional = false;
        for (var index = 0; index < fields.Length; index++)
        {
            var field = fields[index];
            if (optional && !field.Optional)
                throw new ArgumentException($"non-optional command field {field.Name} follows an optional field");
            optional |= field.Optional;
            if (field.Kind == Abi.CommandParameterRawText && index != fields.Length - 1)
                throw new ArgumentException($"varargs command field {field.Name} must be last");
        }
        return new CommandBinding(runnable, fields);
    }

    internal string Usage(string command)
    {
        var values = Fields.Select(field => field.Optional
            ? $"[{field.Name}{field.Suffix}]"
            : $"<{field.Name}{field.Suffix}>");
        return $"/{command} {string.Join(' ', values)}".TrimEnd();
    }

    internal void Bind(IReadOnlyList<string> arguments, Cmd.Source source)
    {
        var argument = 0;
        foreach (var field in Fields)
        {
            if (argument == arguments.Count)
            {
                if (!field.Optional) throw new ArgumentException($"missing command argument {field.Name}");
                field.Clear(Runnable);
                continue;
            }
            field.Set(Runnable, arguments[argument++], source);
        }
        if (argument != arguments.Count) throw new ArgumentException("too many command arguments");
    }
}

internal sealed class CommandField
{
    private CommandField(
        FieldInfo field,
        Type valueType,
        uint kind,
        string name,
        string suffix,
        bool optional,
        string[] values,
        Cmd.Enum? dynamicEnum)
    {
        Field = field;
        ValueType = valueType;
        Kind = kind;
        Name = name;
        Suffix = suffix;
        Optional = optional;
        Values = values;
        DynamicEnum = dynamicEnum;
    }

    internal FieldInfo Field { get; }
    internal Type ValueType { get; }
    internal uint Kind { get; }
    internal string Name { get; }
    internal string Suffix { get; }
    internal bool Optional { get; }
    internal string[] Values { get; }
    internal Cmd.Enum? DynamicEnum { get; }

    [UnconditionalSuppressMessage("Trimming", "IL2062", Justification = "The plugin generator roots command field types for NativeAOT.")]
    [UnconditionalSuppressMessage("Trimming", "IL2072", Justification = "The plugin generator roots command field types for NativeAOT.")]
    internal static CommandField? Create(FieldInfo field, Cmd.Runnable runnable)
    {
        var tag = field.GetCustomAttribute<Cmd.TagAttribute>();
        if (tag?.Name == "-") return null;
        if (field.IsInitOnly) throw new ArgumentException($"command field {field.Name} must be mutable");

        var fieldType = field.FieldType;
        var optional = fieldType.IsGenericType && fieldType.GetGenericTypeDefinition() == typeof(Cmd.Optional<>);
        var valueType = optional ? fieldType.GetGenericArguments()[0] : fieldType;
        var kind = ParameterKind(valueType);
        var values = valueType.IsEnum
            ? System.Enum.GetNames(valueType).Select(value => value.ToLowerInvariant()).ToArray()
            : [];
        Cmd.Enum? dynamicEnum = null;
        if (kind == Abi.CommandParameterDynamicEnum)
        {
            dynamicEnum = field.GetValue(runnable) as Cmd.Enum;
            dynamicEnum ??= Activator.CreateInstance(valueType) as Cmd.Enum;
            if (dynamicEnum is null)
                throw new ArgumentException($"dynamic enum field {field.Name} requires a parameterless constructor");
        }
        var name = (string.IsNullOrEmpty(tag?.Name) ? field.Name : tag.Name).ToLowerInvariant();
        return new CommandField(
            field,
            valueType,
            kind,
            name,
            tag?.Suffix ?? string.Empty,
            optional,
            values,
            dynamicEnum);
    }

    [UnconditionalSuppressMessage("Trimming", "IL2072", Justification = "The plugin generator roots command field types for NativeAOT.")]
    internal void Clear(Cmd.Runnable runnable) => Field.SetValue(runnable, Activator.CreateInstance(Field.FieldType));

    [UnconditionalSuppressMessage("Trimming", "IL2072", Justification = "The plugin generator roots command field types for NativeAOT.")]
    internal void Set(Cmd.Runnable runnable, string argument, Cmd.Source source)
    {
        var value = Parse(argument, source);
        if (Optional) value = Activator.CreateInstance(Field.FieldType, value);
        Field.SetValue(runnable, value);
    }

    [UnconditionalSuppressMessage("Trimming", "IL2072", Justification = "The plugin generator roots command field types for NativeAOT.")]
    private object Parse(string argument, Cmd.Source source)
    {
        if (Kind == Abi.CommandParameterSubcommand)
        {
            if (!string.Equals(argument, Name, StringComparison.OrdinalIgnoreCase))
                throw new ArgumentException($"expected subcommand {Name}");
            return default(Cmd.SubCommand);
        }
        if (Kind == Abi.CommandParameterRawText) return new Cmd.Varargs(argument);
        if (Kind == Abi.CommandParameterPlayer)
            return Player.FromCommandArgument(argument, source is Player player ? player.Invocation : 0);
        if (Kind == Abi.CommandParameterVector)
        {
            var parts = argument.Split(' ', StringSplitOptions.RemoveEmptyEntries);
            if (parts.Length != 3) throw new ArgumentException($"invalid vector {argument}");
            return new Vector3(ParseDouble(parts[0]), ParseDouble(parts[1]), ParseDouble(parts[2]));
        }
        if (ValueType.IsEnum) return System.Enum.Parse(ValueType, argument, true);
        if (Kind == Abi.CommandParameterDynamicEnum)
        {
            var selected = DynamicEnum!.Options(source)
                .FirstOrDefault(option => string.Equals(option, argument, StringComparison.OrdinalIgnoreCase));
            if (selected is null) throw new ArgumentException($"invalid {DynamicEnum.Type()} value {argument}");
            return Activator.CreateInstance(ValueType, selected)
                ?? throw new ArgumentException($"dynamic enum {ValueType.Name} requires a string constructor");
        }
        if (ValueType == typeof(string)) return argument;
        if (ValueType == typeof(bool)) return bool.Parse(argument);
        if (ValueType == typeof(float)) return float.Parse(argument, CultureInfo.InvariantCulture);
        if (ValueType == typeof(double)) return ParseDouble(argument);
        return Convert.ChangeType(argument, ValueType, CultureInfo.InvariantCulture);
    }

    private static uint ParameterKind(Type type)
    {
        if (type == typeof(Cmd.SubCommand)) return Abi.CommandParameterSubcommand;
        if (type == typeof(Cmd.Varargs)) return Abi.CommandParameterRawText;
        if (type == typeof(Player)) return Abi.CommandParameterPlayer;
        if (type == typeof(Vector3)) return Abi.CommandParameterVector;
        if (type.IsEnum) return Abi.CommandParameterEnum;
        if (typeof(Cmd.Enum).IsAssignableFrom(type)) return Abi.CommandParameterDynamicEnum;
        if (type == typeof(string)) return Abi.CommandParameterString;
        if (type == typeof(bool)) return Abi.CommandParameterBool;
        if (type == typeof(float) || type == typeof(double) || type == typeof(decimal)) return Abi.CommandParameterFloat;
        if (type == typeof(byte) || type == typeof(sbyte) || type == typeof(short) || type == typeof(ushort) ||
            type == typeof(int) || type == typeof(uint) || type == typeof(long) || type == typeof(ulong))
            return Abi.CommandParameterInteger;
        throw new ArgumentException($"unsupported command field type {type}");
    }

    private static double ParseDouble(string value) => double.Parse(value, CultureInfo.InvariantCulture);
}
