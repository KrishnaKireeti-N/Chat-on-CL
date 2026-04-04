package command

import (
	"fmt"
	"testing"
)

type App struct {
}

func Exit(x int) {
	fmt.Println("Exit")
}

func TestParse(t *testing.T) {
	p := NewParser[App]([]string{"-h", "--opt"}, Exit)
	args_log := struct{ depth int64 }{}
	p.AddArg(Arg[App]{
		Type:  ArgCommand,
		Name:  "log",
		Sname: "",
		Desc:  "Logs the data",
		Args:  []string{"depth"},
		Parse: func(ctx Context) {
			args_log.depth = p.GetArgInt(ctx)
		},
		Call: func(a *App) {
			fmt.Println("Log")
		},
	})
	p.AddArg(Arg[App]{
		Type:  ArgOption,
		Name:  "opt",
		Sname: "",
		Desc:  "idk",
		Args:  []string{},
		Parse: func(ctx Context) {
		},
		Call: func(a *App) {
			fmt.Println("opt")
		},
	})

	c1, c2 := p.Parse()
	a := &App{}
	c1(a)
	c2(a)
}

func BenchmarkParse(b *testing.B) {
}
