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
	"fmt"
	"strconv"
	"strings"
)

type (
	ArgType    int8
	Arg[T any] struct {
		Type ArgType

		Name  string
		Sname string // shorter name

		Desc string // What the command does

		// Args taken by the command & the containter to be filled with the arguments
		Args []string

		// Argument parser
		Parse func(ctx Context)

		// Function which is returned as per the command given
		Call func(*T)
	}

	// Parser GetX() Context
	Context struct {
		usage string
		args  []string
		idx   int
	}

	Parser[T any] struct {
		args []string
		cmds []Arg[T]

		idx int

		exit func(int)
	}
)

const (
	ArgCommand ArgType = iota
	ArgOption
)

func (cmd *Arg[T]) usage() string {
	if cmd.Sname == "" {
		return fmt.Sprintf("%v %v", cmd.Name, args_Arg(cmd))
	}
	return fmt.Sprintf("%v (%v) %v", cmd.Name, cmd.Sname, args_Arg(cmd))
}

// AddArg adds a command to the command table which is used to parse a command
// and return the correspoding function after parsing is done
//
// Arg with different ArgType and same name is UB
func (p *Parser[T]) AddArg(cmd Arg[T]) {
	p.cmds = append(p.cmds, cmd)
}

func NewParser[T any](args []string, exit func(int)) *Parser[T] {
	parser := &Parser[T]{
		args: args,
		idx:  0,
		exit: exit,
	}
	parser.cmds = make([]Arg[T], 0, 3)

	parser.cmds = append(parser.cmds, Arg[T]{
		Type: ArgCommand,
		Name: "help", Sname: "h", Desc: "Displays the help text", Args: nil,
		Parse: func(ctx Context) {},
		Call:  cli_help(parser),
	})

	return parser
}

func (p *Parser[T]) Parse() (func(*T), func(*T)) {
	var match_cmd *Arg[T]
	var match_opts []*Arg[T]

	if len(p.args) < 1 {
		fmt.Println("Use -h for help")
		p.exit(1)
		return nil, nil
	}

	for p.idx < len(p.args) {
		s := p.args[p.idx]
		p.idx += 1

		if strings.HasPrefix(s, "--") {
			opt, ok := p.parse_Arg(s)
			if !ok {
				fmt.Printf("Given option '%v' doesn't exist\n", s)
				p.exit(1)
			}

			match_opts = append(match_opts, opt)

			opt.Parse(Context{args: opt.Args, idx: 0, usage: opt.usage()})
		} else {
			cmd, ok := p.parse_Arg(s)
			if !ok {
				fmt.Printf("Given command '%v' doesn't exist\n", s)
				p.exit(1)
			}

			if match_cmd != nil {
				fmt.Println("Only 1 command! Use -h for help")
				p.exit(1)
			}

			match_cmd = cmd

			cmd.Parse(Context{args: cmd.Args, idx: 0, usage: cmd.usage()})
		}
		continue
	}

	if match_cmd == nil {
		fmt.Println("Atleast 1 command must be given!\n'help' is a command")
		p.exit(1)
		return nil, nil
	}

	match_opts_call := func(t *T) {
		for i := range match_opts {
			match_opts[i].Call(t)
		}
	}

	if match_cmd.Name == "help" {
		match_cmd.Call(nil)
	}
	return match_cmd.Call, match_opts_call
}

func (p *Parser[T]) parse_Arg(target string) (*Arg[T], bool) {
	target, _ = strings.CutPrefix(target, "--")
	target, _ = strings.CutPrefix(target, "-")
	for i := range p.cmds {
		cmd := &(p.cmds[i])
		if target == cmd.Sname || target == cmd.Name {
			return cmd, true
		}
	}

	return nil, false
}

func (p *Parser[T]) GetArgString(ctx Context) string {
	retstr := p.args[p.idx]

	if checkArgsAmount(retstr, ctx) {
		p.exit(1)
	}
	p.idx++

	p.idx += 1
	ctx.idx += 1
	return retstr
}

func (p *Parser[T]) GetArgInt(ctx Context) int64 {
	retstr := p.args[p.idx]

	if checkArgsAmount(retstr, ctx) {
		p.exit(1)
	}
	p.idx++

	retint, err := strconv.ParseInt(retstr, 10, 64)
	if err != nil {
		fmt.Printf(
			"Provide a proper argument for parameter '%v'\nGiven: %v\n",
			ctx.args[ctx.idx], retstr)
		p.exit(1)
	}
	ctx.idx += 1

	return retint
}

func GetArgX[T any, X any](p *Parser[T], Cast func(string) (X, error)) func(Context) X {
	return func(ctx Context) X {
		retstr := p.args[p.idx]

		if checkArgsAmount(retstr, ctx) {
			p.exit(1)
		}
		p.idx++

		retX, err := Cast(retstr)
		if err != nil {
			fmt.Printf(
				"Provide a proper argument for parameter '%v' | Given: %v\n%v\n",
				ctx.args[ctx.idx], retstr, err.Error())
			p.exit(1)
		}
		ctx.idx += 1

		return retX
	}
}

func checkCmdOrOpt(t string) bool {
	return strings.HasPrefix(t, "-") || strings.HasPrefix(t, "--")
}

// Return true if user gave wrong
func checkArgsAmount(retstr string, ctx Context) bool {
	if checkCmdOrOpt(retstr) {
		fmt.Printf("Not enough arguments!!!\n%v\n", ctx.usage)
		return true
	}

	if ctx.idx > len(ctx.args) {
		fmt.Printf("Too many arguments!!\n%v", ctx.usage)
		return true
	}

	return false
}

func cli_help[T any](a *Parser[T]) func(*T) {
	return func(_ *T) {
		fmt.Println("COMMANDS (specified as '-short <args>' or 'long <args>'):")
		for i := range a.cmds {
			cmd := &(a.cmds[i])
			if cmd.Type == ArgCommand {
				if cmd.Sname != "" {
					fmt.Printf("%v (%v) %v:\n\t%v\n", cmd.Name, cmd.Sname, args_Arg(cmd), cmd.Desc)
				} else {
					fmt.Printf("%v %v:\n\t%v\n", cmd.Name, args_Arg(cmd), cmd.Desc)
				}
			}
		}

		fmt.Println("\nOPTIONS (specified as '--option <args>'")
		for i := range a.cmds {
			opt := &(a.cmds[i])

			if opt.Type == ArgOption {
				if opt.Sname != "" {
					fmt.Printf("%v (%v) %v:\n\t%v\n", opt.Name, opt.Sname, args_Arg(opt), opt.Desc)
				} else {
					fmt.Printf("%v %v:\n\t%v\n", opt.Name, args_Arg(opt), opt.Desc)
				}
			}
		}

		a.exit(0)
	}
}

func args_Arg[T any](cmd *Arg[T]) string {
	var args strings.Builder
	args.Grow(len(cmd.Args) * 6)
	for _, arg := range cmd.Args {
		fmt.Fprintf(&args, "<%v> ", arg)
	}

	return args.String()
}
