package command

/*
This comment describes the intended use:
 An excecutable can have 'commands' and 'options'
 option  : These are values that are used to change the default
           value of the config of the excecutable, w.r.t to a
           command, for the running process
 command : Command is something that is excecuted, that does
           something real and uses the defaults of config, or
           uses the option value provided in command line. A
           command can & should have the privilige to change
           the defaults of a config but should be stated to the
           user

 The option can change the config at the global(process) level
 or just at the command level and hence is left to the
 implementation of the user
*/

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	command struct {
		// name, shorter-name, description
		name  string
		sname string
		desc  string

		// Args taken by the command & the containter to be filled with the arguments
		args []string
		fill reflect.Value

		// Function which is returned as per the command given
		call func()
	}
	option struct {
		// name, shorter-name, description
		name  string
		sname string
		desc  string

		// Arg (only 1) taken by the command & the containter to be filled with the arguments
		args []string
		fill reflect.Value
	}

	ArgParser struct {
		args    []string
		cmds    []command
		options []option

		err error
	}
)

// AddCmd adds a command to the command table which is used to parse a command
// and return the correspoding function after parsing is done
func (p *ArgParser) AddCmd(name string, sname string, desc string, args []string, fill any, call func()) {
	if len(args) != 0 {
		if filltype := reflect.TypeOf(fill); filltype.Kind() != reflect.Pointer {
			p.err = fmt.Errorf("Error: AddCmd(%v): fill must be a pointer", name)
			return
		}
		if fillval := reflect.ValueOf(fill); fillval.Elem().Kind() != reflect.Struct || fillval.IsNil() {
			p.err = fmt.Errorf("Error: AddCmd(%v): fill must be a pointer to a struct", name)
			return
		}
	}

	p.cmds = append(p.cmds, command{name: name, sname: sname, desc: desc, args: args, fill: reflect.ValueOf(fill), call: call})
}

// AddCmd adds a command to the command table which is used to parse a command
// and return the correspoding function after parsing is done
func (p *ArgParser) AddOption(name string, sname string, desc string, args []string, fill any) {
	if len(args) != 0 {
		if filltype := reflect.TypeOf(fill); filltype.Kind() != reflect.Pointer {
			p.err = fmt.Errorf("Error: AddCmd(%v): fill must be a pointer", name)
			return
		}
		if fillval := reflect.ValueOf(fill); fillval.Elem().Kind() != reflect.Struct || fillval.IsNil() {
			p.err = fmt.Errorf("Error: AddCmd(%v): fill must be a pointer to a struct", name)
			return
		}
	}

	p.options = append(p.options, option{name: name, sname: sname, desc: desc, args: args, fill: reflect.ValueOf(fill)})
}

func NewParser(args []string) *ArgParser {
	parser := ArgParser{
		args: args,
	}
	parser.cmds = append(parser.cmds, command{
		name: "help", sname: "h", desc: "Displays the help text", args: nil, call: cli_help(&parser),
	})
	return &parser
}

func (p *ArgParser) Parse() (func(), error) {
	if p.err != nil {
		return nil, p.err
	}

	found_cmd := false
	var exec func() = cli_help(p)

	if len(p.args) < 1 {
		return nil, errors.New("Provide sufficient arguments!!!")
	}

	for i := 0; i < len(p.args); i++ {
		s := p.args[i]
		if len(s) < 2 {
			return nil, errors.New("There seems to be extra arguments given!\nUse -h for help")
		}

		if s[:2] == "--" {
			opt, ok := parse_option(s[2:], p.options)
			if !ok {
				return nil, fmt.Errorf("Given option '%v' doesn't exist", s)
			}
			err := parse_args(p.args[i+1:i+len(opt.args)+1], opt.args, opt.fill)
			if err != nil {
				var args strings.Builder
				args.Grow(len(opt.args) * 6)
				for _, arg := range opt.args {
					fmt.Fprintf(&args, "<%v> ", arg)
				}
				return nil, fmt.Errorf(err.Error()+"%v", fmt.Sprintf("%v %v:\n\t%v\n", opt.name, args.String(), opt.desc))
			}

			i += len(opt.args)
			continue
		} else {
			if found_cmd {
				return nil, errors.New("Only 1 command! Use -h for help")
			}
			found_cmd = true

			cmd, ok := parse_cmd(s, p.cmds)
			if !ok {
				return nil, fmt.Errorf("Given command '%v' doesn't exist", s)
			}

			exec = cmd.call
			if cmd.args == nil {
				break
			}

			cmd_usage := func() string {
				var args strings.Builder
				args.Grow(len(cmd.args) * 6)
				fmt.Fprintf(&args, "%v ", cmd.name)
				for _, arg := range cmd.args {
					fmt.Fprintf(&args, "<%v> ", arg)
				}
				return args.String()
			}
			if len(p.args)-i-1 < len(cmd.args) {
				return nil, fmt.Errorf("Not enough arguements!!!\n%v", cmd_usage())
			}
			err := parse_args(p.args[i+1:i+len(cmd.args)+1], cmd.args, cmd.fill)
			if err != nil {
				return nil, fmt.Errorf("Wrong Arguments!!!\n%v", cmd_usage())
			}

			i += len(cmd.args)
			continue
		}
	}

	return exec, nil
}

func parse_option(target_option string, options []option) (option, bool) {
	for _, opt := range options {
		if target_option == opt.name || target_option == opt.sname {
			return opt, true
		}
	}

	return option{}, false
}

func parse_cmd(target_cmd string, cmds []command) (command, bool) {
	for _, cmd := range cmds {
		if target_cmd[1:] == cmd.sname || target_cmd == cmd.name {
			return cmd, true
		}
	}

	return command{}, false
}

func parse_args(target_args []string, args []string, fill reflect.Value) error {
	var (
		i   int
		arg string
	)
	for i, arg = range target_args {
		err := parse_struct(arg, args[i], fill)
		if err != nil {
			return err
		}
	}
	if i != len(args)-1 {
		return fmt.Errorf("All arguments must be provided!\n")
	}

	return nil
}

func parse_struct(target_arg string, arg string, pfill reflect.Value) error {
	fill := pfill.Elem()

	for i := 0; i < fill.NumField(); i++ {
		fieldvalue := fill.Field(i)
		field := fill.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		if field.Tag.Get("command") == arg {
			switch field.Type.Kind() {
			case reflect.Int:
				x, err := strconv.ParseInt(target_arg, 10, 64)
				if err != nil {
					return fmt.Errorf("Provide a proper integer for argument '%v'\nGiven integer: %v\n", arg, target_arg)
				}
				fieldvalue.SetInt(x)
			case reflect.String:
				fieldvalue.SetString(target_arg)
			}
		}
	}

	return nil
}

func cli_help(a *ArgParser) func() {
	return func() {
		fmt.Println("COMMANDS (can be specified as '-short' or 'long'):")
		for _, cmd := range a.cmds {

			var args strings.Builder
			args.Grow(len(cmd.args) * 6)
			for _, arg := range cmd.args {
				fmt.Fprintf(&args, "<%v> ", arg)
			}

			if cmd.sname != "" {
				fmt.Printf("%v (%v) %v:\n\t%v\n", cmd.name, cmd.sname, args.String(), cmd.desc)
			} else {
				fmt.Printf("%v %v:\n\t%v\n", cmd.name, args.String(), cmd.desc)
			}
		}
	}
}
