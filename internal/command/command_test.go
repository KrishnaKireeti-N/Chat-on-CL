package command

import (
	"fmt"
	"os"
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

func TestParse(t *testing.T) {
	// Better implementation
	color := map[string]string{
		"Red":   "#ff0000",
		"Green": "00ff00",
		"Blue":  "#0000ff",
	}

	args_list := struct {
		Item string `command:"item"`
	}{}
	args_option := struct {
		Color string `command:"color"`
		priv  string `command:"color"`
	}{priv: "private"}

	p := NewParser([]string{"-l", "colors", "--color", "Red"})

	p.AddCmd("colors", "", "Prints all the colors available", nil, nil, func() { fmt.Println("Red Blue") })
	p.AddCmd("list", "l", "List down the default configuration", []string{"item"}, &args_list, func() {
		switch args_list.Item {
		case "colors":
			for c := Color(0); c != Invalid; c++ {
				t.Log(c)
			}
		}
	})
	p.AddOption("color", "", "Set the color of the text", []string{"color"}, &args_option)

	c, err := p.Parse()
	if err != nil {
		t.Error(err.Error())
		os.Exit(1)
	}
	c()

	t.Log(args_list)
	t.Log(color[args_option.Color])
}
