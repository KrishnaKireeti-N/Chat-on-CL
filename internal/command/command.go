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

		// tag (cached for performance)
		tag map[string]int

		// Function which is returned as per the command given
		call func()
	}
	option struct {
		// name, shorter-name, description
		name  string
		sname string
		desc  string

		// Args taken by the option & the container to be filled with the arguments
		args []string
		fill reflect.Value

		// tag (cached for performance)
		tag map[string]int
	}

	// interface so as to pass both command and option to helper-functions
	parsable interface {
		getArgs() []string
		getFill() reflect.Value
		getTags() map[string]int
	}

	Parser struct {
		args    []string
		cmds    []command
		options []option

		err error
	}
)

func (cmd *command) getArgs() []string       { return cmd.args }
func (cmd *command) getFill() reflect.Value  { return cmd.fill }
func (cmd *command) getTags() map[string]int { return cmd.tag }
func (opt *option) getArgs() []string        { return opt.args }
func (opt *option) getFill() reflect.Value   { return opt.fill }
func (opt *option) getTags() map[string]int  { return opt.tag }

// AddCmd adds a command to the command table which is used to parse a command
// and return the correspoding function after parsing is done
func (p *Parser) AddCmd(name string, sname string, desc string, args []string, fill any, call func()) {
	cmd := command{name: name, sname: sname, desc: desc, args: args, fill: reflect.ValueOf(fill), call: call}

	if len(args) != 0 {
		if filltype := reflect.TypeOf(fill); filltype.Kind() != reflect.Pointer {
			panic(fmt.Sprintf("Error: AddCmd(%v): fill must be a pointer", name))
		}
		if fillval := reflect.ValueOf(fill); fillval.Elem().Kind() != reflect.Struct || fillval.IsNil() {
			panic(fmt.Sprintf("Error: AddCmd(%v): fill must be a pointer to a struct", name))
		}
		if cmd.tag == nil {
			fill := cmd.fill.Elem()
			cmd.tag = make(map[string]int, fill.NumField())

			for i := 0; i < fill.NumField(); i++ {
				field := fill.Type().Field(i)

				if !field.IsExported() {
					continue
				}

				tag, ok := field.Tag.Lookup("command")
				if !ok {
					panic("Parser requires struct tags to be set for parsing arguments!")
				}
				cmd.tag[tag] = i
			}
		}
	}

	p.cmds = append(p.cmds, cmd)
}

// AddCmd adds a command to the command table which is used to parse a command
// and return the correspoding function after parsing is done
func (p *Parser) AddOption(name string, sname string, desc string, args []string, fill any) {
	opt := option{name: name, sname: sname, desc: desc, args: args, fill: reflect.ValueOf(fill)}

	if len(args) != 0 {
		if filltype := reflect.TypeOf(fill); filltype.Kind() != reflect.Pointer {
			panic(fmt.Sprintf("Error: AddOption(%v): fill must be a pointer", name))
		}
		if fillval := reflect.ValueOf(fill); fillval.Elem().Kind() != reflect.Struct || fillval.IsNil() {
			panic(fmt.Sprintf("Error: AddCmd(%v): fill must be a pointer to a struct", name))
		}
		if opt.tag == nil {
			fill := opt.fill.Elem()
			opt.tag = make(map[string]int, fill.NumField())

			for i := 0; i < fill.NumField(); i++ {
				field := fill.Type().Field(i)

				if !field.IsExported() {
					continue
				}

				tag, ok := field.Tag.Lookup("command")
				if !ok {
					panic("Parser requires struct tags to be set for parsing arguments!")
				}
				opt.tag[tag] = i
			}
		}
	}

	p.options = append(p.options, opt)
}

func NewParser(args []string, ncmds int, noptions int) *Parser {
	parser := &Parser{
		args: args,
	}
	parser.cmds = make([]command, 0, ncmds+1)
	parser.options = make([]option, 0, noptions)

	parser.cmds = append(parser.cmds, command{
		name: "help", sname: "h", desc: "Displays the help text", args: nil, call: cli_help(parser),
	})

	return parser
}

func (p *Parser) Parse() (func(), error) {
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
			opt, ok := p.parse_option(s[2:])
			if !ok {
				return nil, fmt.Errorf("Given option '%v' doesn't exist", s)
			}
			if len(p.args)-i-1 < len(opt.args) {
				return nil, fmt.Errorf("Not enough arguements!!!\n%v", fmt.Sprintf("%v %v\n", opt.name, args_option(opt)))
			}
			err := p.parse_args(p.args[i+1:i+len(opt.args)+1], opt)
			if err != nil {
				return nil, fmt.Errorf(err.Error()+"%v", fmt.Sprintf("%v %v:\n\t%v\n", opt.name, args_option(opt), opt.desc))
			}

			i += len(opt.args)
			continue
		} else {
			if found_cmd {
				return nil, errors.New("Only 1 command! Use -h for help")
			}
			found_cmd = true

			cmd, ok := p.parse_cmd(s)
			if !ok {
				return nil, fmt.Errorf("Given command '%v' doesn't exist", s)
			}

			exec = cmd.call
			if cmd.args == nil {
				break
			}

			cmd_usage := func() string {
				if cmd.sname == "" {
					return fmt.Sprintf("%v %v", cmd.name, args_cmd(cmd))
				}
				return fmt.Sprintf("%v (%v) %v", cmd.name, cmd.sname, args_cmd(cmd))
			}
			if len(p.args)-i-1 < len(cmd.args) {
				return nil, fmt.Errorf("Not enough arguements!!!\n%v\n", cmd_usage())
			}
			err := p.parse_args(p.args[i+1:i+len(cmd.args)+1], cmd)
			if err != nil {
				return nil, fmt.Errorf("Wrong Arguments!!!\n%v\n", cmd_usage())
			}

			i += len(cmd.args)
			continue
		}
	}

	return exec, nil
}

func (p *Parser) parse_option(target_option string) (*option, bool) {
	for i := range p.options {
		opt := &(p.options[i])
		if target_option == opt.name || target_option == opt.sname {
			return opt, true
		}
	}

	return nil, false
}

func (p *Parser) parse_cmd(target_cmd string) (*command, bool) {
	for i := range p.cmds {
		cmd := &(p.cmds[i])
		if target_cmd[1:] == cmd.sname || target_cmd == cmd.name {
			return cmd, true
		}
	}

	return nil, false
}

func (p *Parser) parse_args(target_args []string, v parsable) error {
	var args []string = v.getArgs()

	if len(target_args) != len(args) {
		return fmt.Errorf("All arguments must be provided!\n")
	}

	for i := range target_args {
		arg := target_args[i]
		err := p.parse_struct(arg, args[i], v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parse_struct(arg string, fill_arg string, v parsable) error {
	fill := v.getFill().Elem()
	fieldi := v.getTags()[fill_arg]

	switch fill.Field(fieldi).Type().Kind() {
	case reflect.Int:
		x, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("Provide a proper integer for argument '%v'\nGiven integer: %v\n", fill_arg, arg)
		}
		fill.Field(fieldi).SetInt(x)
	case reflect.String:
		fill.Field(fieldi).SetString(arg)
	}

	return nil
}

func cli_help(a *Parser) func() {
	return func() {
		fmt.Println("COMMANDS (specified as '-short <args>' or 'long <args>'):")
		for i := range a.cmds {
			cmd := &(a.cmds[i])

			if cmd.sname != "" {
				fmt.Printf("%v (%v) %v:\n\t%v\n", cmd.name, cmd.sname, args_cmd(cmd), cmd.desc)
			} else {
				fmt.Printf("%v %v:\n\t%v\n", cmd.name, args_cmd(cmd), cmd.desc)
			}
		}

		fmt.Println("\nOPTIONS (specified as '--option <args>'")
		for i := range a.options {
			opt := &(a.options[i])

			if opt.sname != "" {
				fmt.Printf("%v (%v) %v:\n\t%v\n", opt.name, opt.sname, args_option(opt), opt.desc)
			} else {
				fmt.Printf("%v %v:\n\t%v\n", opt.name, args_option(opt), opt.desc)
			}
		}
	}
}

func args_cmd(cmd *command) string {
	var args strings.Builder
	args.Grow(len(cmd.args) * 6)
	for _, arg := range cmd.args {
		fmt.Fprintf(&args, "<%v> ", arg)
	}

	return args.String()
}
func args_option(opt *option) string {
	var args strings.Builder
	args.Grow(len(opt.args) * 6)
	for _, arg := range opt.args {
		fmt.Fprintf(&args, "<%v> ", arg)
	}

	return args.String()
}
