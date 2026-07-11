package host

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

type commandRuntime interface {
	Commands() ([]native.Command, error)
	HandleCommand(index uint64, input native.CommandInput) (native.CommandOutput, error)
}

// RegisterCommands publishes all enabled plugin commands in Dragonfly's command registry.
func RegisterCommands(runtime commandRuntime) error {
	commands, err := runtime.Commands()
	if err != nil {
		return err
	}
	for _, command := range commands {
		if _, exists := cmd.ByAlias(command.Name); exists {
			return fmt.Errorf("register plugin command %q: name already registered", command.Name)
		}
		runnables, err := commandRunnables(runtime, command)
		if err != nil {
			return fmt.Errorf("register plugin command %q: %w", command.Name, err)
		}
		cmd.Register(cmd.New(command.Name, command.Description, nil, runnables...))
	}
	return nil
}

func commandRunnables(runtime commandRuntime, command native.Command) ([]cmd.Runnable, error) {
	base := pluginCommandBase{runtime: runtime, index: command.Index}
	if len(command.Overloads) == 0 {
		return []cmd.Runnable{pluginCommand{pluginCommandBase: base}}, nil
	}
	runnables := make([]cmd.Runnable, 0, len(command.Overloads))
	for overloadIndex, overload := range command.Overloads {
		parameters := make([]pluginParameter, len(overload.Parameters))
		for index, parameter := range overload.Parameters {
			parameters[index] = pluginParameter{
				descriptor: parameter,
				enumType: fmt.Sprintf(
					"bedrock_gophers_%s_%d_%s",
					command.Name,
					overloadIndex,
					parameter.Name,
				),
			}
		}
		switch len(parameters) {
		case 0:
			runnables = append(runnables, pluginCommand{pluginCommandBase: base})
		case 1:
			runnables = append(runnables, pluginCommand1{pluginCommandBase: base, P1: parameters[0]})
		case 2:
			runnables = append(runnables, pluginCommand2{pluginCommandBase: base, P1: parameters[0], P2: parameters[1]})
		case 3:
			runnables = append(runnables, pluginCommand3{pluginCommandBase: base, P1: parameters[0], P2: parameters[1], P3: parameters[2]})
		case 4:
			runnables = append(runnables, pluginCommand4{pluginCommandBase: base, P1: parameters[0], P2: parameters[1], P3: parameters[2], P4: parameters[3]})
		default:
			return nil, fmt.Errorf("overloads support at most four typed parameters")
		}
	}
	return runnables, nil
}

type pluginCommandBase struct {
	runtime commandRuntime
	index   uint64
}

func (c pluginCommandBase) dispatch(source cmd.Source, output *cmd.Output, arguments []string, trailing cmd.Varargs) {
	if value := strings.TrimSpace(string(trailing)); value != "" {
		arguments = append(arguments, value)
	}
	sourceName := fmt.Sprintf("%T", source)
	if named, ok := source.(cmd.NamedTarget); ok {
		sourceName = named.Name()
	}
	result, err := c.runtime.HandleCommand(c.index, native.CommandInput{
		Source:    sourceName,
		Arguments: strings.Join(arguments, " "),
	})
	if err != nil {
		output.Error(err)
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
	Arguments         cmd.Varargs `cmd:"arguments"`
}

func (c pluginCommand) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	c.dispatch(source, output, nil, c.Arguments)
}

type pluginParameter struct {
	descriptor native.CommandParameter
	enumType   string
	selected   string
}

func (p pluginParameter) Parse(line *cmd.Line, value reflect.Value) error {
	argument, ok := line.Next()
	if !ok {
		return fmt.Errorf("missing %s", p.descriptor.Name)
	}
	switch p.descriptor.Kind {
	case native.CommandParameterSubcommand:
		if !strings.EqualFold(argument, p.descriptor.Name) {
			return fmt.Errorf("expected %s", p.descriptor.Name)
		}
		p.selected = p.descriptor.Name
	case native.CommandParameterEnum:
		for _, option := range p.descriptor.Values {
			if strings.EqualFold(argument, option) {
				p.selected = option
				value.Set(reflect.ValueOf(p))
				return nil
			}
		}
		return fmt.Errorf("%q is not a valid %s", argument, p.descriptor.Name)
	default:
		return fmt.Errorf("unknown plugin parameter kind %d", p.descriptor.Kind)
	}
	value.Set(reflect.ValueOf(p))
	return nil
}

func (p pluginParameter) Type() string { return p.descriptor.Name }

type describedEnum struct {
	typeName string
	options  []string
}

func (e describedEnum) Type() string                { return e.typeName }
func (e describedEnum) Options(cmd.Source) []string { return e.options }

func describe(parameters ...pluginParameter) []cmd.ParamInfo {
	result := make([]cmd.ParamInfo, 0, len(parameters)+1)
	for _, parameter := range parameters {
		if parameter.descriptor.Kind == native.CommandParameterSubcommand {
			result = append(result, cmd.ParamInfo{Name: parameter.descriptor.Name, Value: cmd.SubCommand{}})
		} else {
			result = append(result, cmd.ParamInfo{
				Name: parameter.descriptor.Name,
				Value: describedEnum{
					typeName: parameter.enumType,
					options:  parameter.descriptor.Values,
				},
			})
		}
	}
	result = append(result, cmd.ParamInfo{Name: "arguments", Value: cmd.Varargs(""), Optional: true})
	return result
}

type pluginCommand1 struct {
	pluginCommandBase `cmd:"-"`
	P1                pluginParameter
	Arguments         cmd.Varargs `cmd:"arguments"`
}

func (c pluginCommand1) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	c.dispatch(source, output, []string{c.P1.selected}, c.Arguments)
}
func (c pluginCommand1) DescribeParams(cmd.Source) []cmd.ParamInfo { return describe(c.P1) }

type pluginCommand2 struct {
	pluginCommandBase `cmd:"-"`
	P1, P2            pluginParameter
	Arguments         cmd.Varargs `cmd:"arguments"`
}

func (c pluginCommand2) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	c.dispatch(source, output, []string{c.P1.selected, c.P2.selected}, c.Arguments)
}
func (c pluginCommand2) DescribeParams(cmd.Source) []cmd.ParamInfo { return describe(c.P1, c.P2) }

type pluginCommand3 struct {
	pluginCommandBase `cmd:"-"`
	P1, P2, P3        pluginParameter
	Arguments         cmd.Varargs `cmd:"arguments"`
}

func (c pluginCommand3) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	c.dispatch(source, output, []string{c.P1.selected, c.P2.selected, c.P3.selected}, c.Arguments)
}
func (c pluginCommand3) DescribeParams(cmd.Source) []cmd.ParamInfo {
	return describe(c.P1, c.P2, c.P3)
}

type pluginCommand4 struct {
	pluginCommandBase `cmd:"-"`
	P1, P2, P3, P4    pluginParameter
	Arguments         cmd.Varargs `cmd:"arguments"`
}

func (c pluginCommand4) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	c.dispatch(source, output, []string{c.P1.selected, c.P2.selected, c.P3.selected, c.P4.selected}, c.Arguments)
}
func (c pluginCommand4) DescribeParams(cmd.Source) []cmd.ParamInfo {
	return describe(c.P1, c.P2, c.P3, c.P4)
}
