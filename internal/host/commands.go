package host

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

type commandRuntime interface {
	Commands() ([]native.Command, error)
	HandleCommand(index uint64, input native.CommandInput) (native.CommandOutput, error)
	CommandEnumOptions(index, overload, parameter uint64, input native.CommandEnumContext) ([]string, error)
}

// RegisterCommands publishes all enabled plugin commands in Dragonfly's command registry.
func RegisterCommands(runtime commandRuntime, players *Players) error {
	commands, err := runtime.Commands()
	if err != nil {
		return err
	}
	for _, command := range commands {
		for _, alias := range append([]string{command.Name}, command.Aliases...) {
			if _, exists := cmd.ByAlias(strings.ToLower(alias)); exists {
				return fmt.Errorf("register plugin command %q: alias %q already registered", command.Name, alias)
			}
		}
		runnables, err := commandRunnables(runtime, players, command)
		if err != nil {
			return fmt.Errorf("register plugin command %q: %w", command.Name, err)
		}
		cmd.Register(cmd.New(command.Name, command.Description, command.Aliases, runnables...))
	}
	return nil
}

func commandRunnables(runtime commandRuntime, players *Players, command native.Command) ([]cmd.Runnable, error) {
	overloads := command.Overloads
	if len(overloads) == 0 {
		overloads = []native.CommandOverload{{}}
	}
	typeCounts := map[string]int{}
	for _, overload := range overloads {
		for _, parameter := range overload.Parameters {
			typeCounts[command.Name+"_"+parameter.Name]++
		}
	}
	runnables := make([]cmd.Runnable, 0, len(overloads))
	for overloadIndex, overload := range overloads {
		if err := validateCommandParameters(overload.Parameters); err != nil {
			return nil, fmt.Errorf("overload %d: %w", overloadIndex, err)
		}
		arguments := pluginArguments{
			descriptors: append([]native.CommandParameter(nil), overload.Parameters...),
			runtime:     runtime,
			players:     players,
			command:     command.Index,
			overload:    uint64(overloadIndex),
			typeCounts:  typeCounts,
			commandName: command.Name,
		}
		runnables = append(runnables, pluginCommand{
			pluginCommandBase: pluginCommandBase{
				runtime:  runtime,
				players:  players,
				index:    command.Index,
				overload: uint64(overloadIndex),
			},
			Arguments: arguments,
		})
	}
	return runnables, nil
}

func validateCommandParameters(parameters []native.CommandParameter) error {
	optional := false
	for index, parameter := range parameters {
		if parameter.Kind < native.CommandParameterSubcommand || parameter.Kind > native.CommandParameterVector {
			return fmt.Errorf("parameter %q has unknown kind %d", parameter.Name, parameter.Kind)
		}
		if !parameter.Optional && optional {
			return fmt.Errorf("parameter %q is required after an optional parameter", parameter.Name)
		}
		optional = optional || parameter.Optional
		if parameter.Kind == native.CommandParameterRawText && index != len(parameters)-1 {
			return fmt.Errorf("raw-text parameter %q must be last", parameter.Name)
		}
	}
	return nil
}

type pluginCommandBase struct {
	runtime  commandRuntime
	players  *Players
	index    uint64
	overload uint64
}

func (c pluginCommandBase) dispatch(source cmd.Source, output *cmd.Output, arguments pluginArguments, tx *world.Tx) {
	c.players.WithInvocation(tx, func(invocation native.InvocationID) {
		c.dispatchActive(invocation, source, output, arguments.transport(invocation))
	})
}

func (c pluginCommandBase) dispatchActive(invocation native.InvocationID, source cmd.Source, output *cmd.Output, arguments []string) {
	sourceName := fmt.Sprintf("%T", source)
	if named, ok := source.(cmd.NamedTarget); ok {
		sourceName = named.Name()
	}
	input := native.CommandInput{
		Invocation:    invocation,
		Overload:      c.overload,
		Source:        sourceName,
		Arguments:     arguments,
		SourceKind:    native.CommandSourceConsole,
		OnlinePlayers: c.players.CommandSnapshots(),
	}
	position := source.Position()
	input.SourcePosition = native.Vec3{X: position[0], Y: position[1], Z: position[2]}
	if sourcePlayer, ok := source.(interface{ UUID() uuid.UUID }); ok {
		playerUUID := sourcePlayer.UUID()
		if id, found := c.players.ResolveUUID(playerUUID); found {
			input.SourcePlayer = &id
			input.SourceKind = native.CommandSourcePlayer
		}
	}
	result, err := c.runtime.HandleCommand(c.index, input)
	if err != nil {
		slog.Error("native plugin command failed", "command", c.index, "overload", c.overload, "source", sourceName, "error", err)
		output.Error("Command failed.")
		return
	}
	if result.Message == "" {
		return
	}
	if result.Failed {
		output.Error(result.Message)
	} else {
		output.Print(result.Message)
	}
}

type pluginCommand struct {
	pluginCommandBase `cmd:"-"`
	Arguments         pluginArguments `cmd:"arguments"`
}

func (c pluginCommand) Run(source cmd.Source, output *cmd.Output, tx *world.Tx) {
	c.dispatch(source, output, c.Arguments, tx)
}

func (c pluginCommand) DescribeParams(source cmd.Source) []cmd.ParamInfo {
	return c.Arguments.describe(source)
}

// pluginArguments is the one physical Dragonfly parameter used to parse one
// complete plugin overload. Logical parameters remain individually described
// through ParamDescriber and individually transported to the plugin.
type pluginArguments struct {
	descriptors []native.CommandParameter
	values      []string
	runtime     commandRuntime
	players     *Players
	command     uint64
	overload    uint64
	typeCounts  map[string]int
	commandName string
}

func (a pluginArguments) Parse(line *cmd.Line, value reflect.Value) error {
	if len(a.descriptors) == 0 && line.Len() != 0 {
		return line.UsageError()
	}
	values := make([]string, 0, len(a.descriptors))
	consumed := 0
	for _, descriptor := range a.descriptors {
		available := line.Len() - consumed
		if available == 0 && descriptor.Optional {
			break
		}
		var parsed string
		var count int
		var err error
		switch descriptor.Kind {
		case native.CommandParameterRawText:
			if available == 0 {
				parsed = ""
			} else {
				parts, _ := line.NextN(line.Len())
				parsed = strings.Join(parts[consumed:], " ")
				count = available
			}
		case native.CommandParameterVector:
			parsed, count, err = parseVector(line, consumed)
		default:
			parsed, count, err = parseCommandArgument(line, consumed, descriptor)
		}
		if err != nil {
			return err
		}
		values = append(values, parsed)
		consumed += count
	}
	// Dragonfly consumes one argument after Parameter.Parse returns. Consume
	// all preceding arguments here and leave the last one to its parser.
	if consumed > 1 {
		line.RemoveN(consumed - 1)
	}
	a.values = values
	value.Set(reflect.ValueOf(a))
	return nil
}

func parseCommandArgument(line *cmd.Line, offset int, descriptor native.CommandParameter) (string, int, error) {
	arguments, ok := line.NextN(offset + 1)
	if !ok {
		return "", 0, line.UsageError()
	}
	argument := arguments[offset]
	switch descriptor.Kind {
	case native.CommandParameterSubcommand:
		if !strings.EqualFold(argument, descriptor.Name) {
			return "", 0, cmd.MessageParameterInvalid.F(argument)
		}
		argument = descriptor.Name
	case native.CommandParameterEnum:
		option, found := enumOption(descriptor.Values, argument)
		if !found {
			return "", 0, cmd.MessageParameterInvalid.F(argument)
		}
		argument = option
	case native.CommandParameterInteger:
		if _, err := strconv.ParseInt(argument, 10, 64); err != nil {
			return "", 0, cmd.MessageNumberInvalid.F(argument)
		}
	case native.CommandParameterFloat:
		if _, err := strconv.ParseFloat(argument, 64); err != nil {
			return "", 0, cmd.MessageNumberInvalid.F(argument)
		}
	case native.CommandParameterBool:
		parsed, err := strconv.ParseBool(argument)
		if err != nil {
			return "", 0, cmd.MessageBooleanInvalid.F(argument)
		}
		argument = strconv.FormatBool(parsed)
	case native.CommandParameterDynamicEnum, native.CommandParameterString, native.CommandParameterPlayer:
	default:
		return "", 0, fmt.Errorf("unknown plugin parameter kind %d", descriptor.Kind)
	}
	return argument, 1, nil
}

func parseVector(line *cmd.Line, offset int) (string, int, error) {
	arguments, ok := line.NextN(offset + 3)
	if !ok {
		return "", 0, line.UsageError()
	}
	parts := arguments[offset : offset+3]
	for _, part := range parts {
		if _, err := strconv.ParseFloat(part, 64); err != nil {
			return "", 0, cmd.MessageNumberInvalid.F(part)
		}
	}
	return strings.Join(parts, " "), 3, nil
}

func enumOption(options []string, argument string) (string, bool) {
	for _, option := range options {
		if strings.EqualFold(option, argument) {
			return option, true
		}
	}
	return "", false
}

func (a pluginArguments) Type() string { return "arguments" }

func (a pluginArguments) transport(invocation native.InvocationID) []string {
	transported := append([]string(nil), a.values...)
	for index, descriptor := range a.descriptors[:len(transported)] {
		if descriptor.Kind != native.CommandParameterPlayer {
			continue
		}
		id, ok := a.players.ResolveName(transported[index])
		if !ok {
			transported[index] = "invalid"
			continue
		}
		connected, ok := a.players.ResolveID(id, invocation)
		if !ok {
			transported[index] = "invalid"
			continue
		}
		position := connected.Position()
		transported[index] = strings.Join([]string{
			hex.EncodeToString(id.UUID[:]),
			strconv.FormatUint(id.Generation, 10),
			strconv.FormatInt(connected.Latency().Milliseconds(), 10),
			strconv.FormatFloat(position[0], 'g', -1, 64),
			strconv.FormatFloat(position[1], 'g', -1, 64),
			strconv.FormatFloat(position[2], 'g', -1, 64),
			connected.Name(),
		}, ":")
	}
	return transported
}

func (a pluginArguments) describe(source cmd.Source) []cmd.ParamInfo {
	result := make([]cmd.ParamInfo, 0, len(a.descriptors))
	for index, descriptor := range a.descriptors {
		info := cmd.ParamInfo{Name: strings.ToLower(descriptor.Name), Optional: descriptor.Optional, Suffix: descriptor.Suffix}
		switch descriptor.Kind {
		case native.CommandParameterSubcommand:
			info.Value = cmd.SubCommand{}
		case native.CommandParameterEnum:
			info.Value = describedEnum{typeName: a.enumType(descriptor.Name, index), options: descriptor.Values}
		case native.CommandParameterDynamicEnum:
			info.Value = describedEnum{
				typeName: a.enumType(descriptor.Name, index), runtime: a.runtime, players: a.players,
				command: a.command, overload: a.overload, parameter: uint64(index),
			}
		case native.CommandParameterString:
			info.Value = ""
		case native.CommandParameterInteger:
			info.Value = int64(0)
		case native.CommandParameterFloat:
			info.Value = float64(0)
		case native.CommandParameterBool:
			info.Value = false
		case native.CommandParameterPlayer:
			info.Value = describedEnum{typeName: a.enumType(descriptor.Name, index), players: a.players, playerNames: true}
		case native.CommandParameterRawText:
			info.Value = cmd.Varargs("")
		case native.CommandParameterVector:
			info.Value = mgl64.Vec3{}
		}
		result = append(result, info)
	}
	return result
}

func (a pluginArguments) enumType(parameter string, _ int) string {
	name := strings.ToLower(a.commandName + "_" + parameter)
	if a.typeCounts[name] > 1 {
		return fmt.Sprintf("%s_%d", name, a.overload)
	}
	return name
}

type describedEnum struct {
	typeName    string
	options     []string
	runtime     commandRuntime
	players     *Players
	command     uint64
	overload    uint64
	parameter   uint64
	playerNames bool
}

func (e describedEnum) Type() string { return e.typeName }

func (e describedEnum) Options(source cmd.Source) []string {
	if e.playerNames {
		return lowercaseOptions(e.players.Names())
	}
	if e.runtime == nil {
		return lowercaseOptions(e.options)
	}
	sourceName := fmt.Sprintf("%T", source)
	if named, ok := source.(cmd.NamedTarget); ok {
		sourceName = named.Name()
	}
	input := native.CommandEnumContext{
		Source: sourceName, SourceKind: native.CommandSourceConsole, OnlinePlayers: e.players.CommandSnapshots(),
	}
	position := source.Position()
	input.SourcePosition = native.Vec3{X: position[0], Y: position[1], Z: position[2]}
	if sourcePlayer, ok := source.(interface{ UUID() uuid.UUID }); ok {
		if id, found := e.players.ResolveUUID(sourcePlayer.UUID()); found {
			input.SourcePlayer = &id
			input.SourceKind = native.CommandSourcePlayer
		}
	}
	options, err := e.runtime.CommandEnumOptions(e.command, e.overload, e.parameter, input)
	if err != nil {
		slog.Error("native plugin command enum failed", "command", e.command, "overload", e.overload, "parameter", e.parameter, "error", err)
		return nil
	}
	return lowercaseOptions(options)
}

func lowercaseOptions(options []string) []string {
	normalized := make([]string, len(options))
	for index, option := range options {
		normalized[index] = strings.ToLower(option)
	}
	return normalized
}
