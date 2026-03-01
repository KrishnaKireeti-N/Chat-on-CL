package command

// Some of the Tests are AI generated!

import (
	"fmt"
	"strconv"
	"testing"
)

type Color int

const (
	Red Color = iota
	Blue
	Invalid
)

func Int(c string) Color {
	switch c {
	case "Red":
		return Red
	case "Blue":
		return Blue
	default:
		return Invalid
	}
}

func (c Color) String() string {
	switch c {
	case Red:
		return "Red"
	case Blue:
		return "Blue"
	default:
		return "Invalid"
	}
}

func TestParse1(t *testing.T) {
	args_list := struct {
		Item string `command:"item"`
	}{}
	args_option := struct {
		Color string `command:"color"`
		priv  string `command:"color"`
	}{priv: "private"}

	//p := NewParser([]string{"-h"}, 2, 1)
	p := NewParser([]string{"-l", "colors", "--color", "Red"}, 2, 1)

	p.AddCmd("colors", "", "Prints all the colors available", nil, nil, func() { fmt.Println("Red Blue") })
	p.AddCmd("list", "l", "List down the default configuration", []string{"item"}, &args_list, func() {
		switch args_list.Item {
		case "colors":
			return
		default:
			t.Fatal("case `colors` should run")
		}
	})
	p.AddOption("color", "", "Set the color of the text", []string{"color"}, &args_option)

	c, err := p.Parse()
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(args_option)
	c()
}

type testCmd struct {
	A int    `command:"a"`
	B string `command:"b"`
}

type testOpt struct {
	X int `command:"x"`
}

func TestParse2(t *testing.T) {
	args := []string{"--opt", "5", "run", "1", "ok"}

	parser := NewParser(args, 1, 1)

	cmdStruct := &testCmd{}
	optStruct := &testOpt{}

	parser.AddOption("opt", "", "test option", []string{"x"}, optStruct)
	parser.AddCmd("run", "", "test command", []string{"a", "b"}, cmdStruct, func() {})

	_, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if optStruct.X != 5 {
		t.Fatalf("expected option X=5, got %d", optStruct.X)
	}
}

func TestInvalidCommand(t *testing.T) {
	args := []string{"invalid"}

	parser := NewParser(args, 1, 0)

	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
}

func BenchmarkParse1(b *testing.B) {
	for b.Loop() {
		args := []string{"run", "42", "world"}

		parser := NewParser(args, 1, 0)

		cmdStruct := &testCmd{}
		parser.AddCmd("run", "", "test command", []string{"a", "b"}, cmdStruct, func() {})

		_, err := parser.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParse100(b *testing.B) {
	for b.Loop() {
		args := []string{"cmd99", "1", "x"}

		parser := NewParser(args, 100, 0)

		for j := range 100 {
			name := "cmd" + strconv.Itoa(j)
			parser.AddCmd(name, "", "desc", []string{"a", "b"}, &testCmd{}, func() {})
		}

		_, err := parser.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseOption1(b *testing.B) {
	for b.Loop() {
		args := []string{"--opt", "3", "run", "1", "hello"}

		parser := NewParser(args, 1, 1)

		parser.AddOption("opt", "", "desc", []string{"x"}, &testOpt{})
		parser.AddCmd("run", "", "desc", []string{"a", "b"}, &testCmd{}, func() {})

		_, err := parser.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}
