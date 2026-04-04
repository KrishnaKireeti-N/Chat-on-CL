package app

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/command"
	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/node"
	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	client *node.Client
	config *node.Config
}

func cmd_parser_init() *command.Parser[App] {
	var (
		// Commands
		args_connect struct {
			Host string `command:"host"`
		}

		// Options
		args_option struct {
			Color string `command:"color"`
		}
	)

	arg_p := command.NewParser[App](os.Args[1:], os.Exit)
	arg_p.AddArg(command.Arg[App]{
		Type: command.ArgCommand,
		Name: "listen", Sname: "l", Desc: "Listen for incoming connections", Args: []string{},
		Parse: func(ctx command.Context) {},
		Call: func(a *App) {
			var err error

			a.client, err = node.MakeUserListen(a.config)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	})
	arg_p.AddArg(command.Arg[App]{
		Type: command.ArgCommand,
		Name: "connect", Sname: "c", Desc: "Connect to the given user", Args: []string{"host"},
		Parse: func(ctx command.Context) {
			args_connect.Host = command.GetArgX(arg_p, func(s string) (string, error) {
				if strings.HasSuffix(s, ":") {
					return s, fmt.Errorf("Remove the ':' at the end")
				}
				return s, nil
			})(ctx)
		},
		Call: func(a *App) {
			var err error

			a.client, err = node.MakeUserConnect(a.config, args_connect.Host)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	})
	arg_p.AddArg(command.Arg[App]{
		Type: command.ArgOption,
		Name: "color", Sname: "c", Desc: "Set the color of your messages (Red, Green, Blue, ...)",
		Args: []string{"color"},
		Parse: func(ctx command.Context) {
			args_option.Color = command.GetArgX(arg_p, func(s string) (string, error) {
				switch s {
				case "Red":
					return "#FF0000", nil
				case "Green":
					return "$00FF00", nil
				case "Blue":
					return "#0000FF", nil
				case "Yellow":
					return "#FFFF00", nil
				case "Purple":
					return "#800080", nil
				case "Cyan":
					return "#48D1CC", nil // mediumturquoise
				case "White":
					return "#FFFFFF", nil
				}
				return "", fmt.Errorf("The color is not supported")
			})(ctx)
		},
		Call: func(a *App) {
			a.config.Color = args_option.Color
		},
	})

	return arg_p
}

func NewApp() *App {
	parser := cmd_parser_init()
	cmd, set_opts := parser.Parse()

	a := &App{}

	a.config = node.LoadConfig()
	defer node.SaveConfig(a.config)

	// Managing the config according to the options provided
	a.config.SetDefaults()
	set_opts(a)
	fmt.Println(a.config)

	cmd(a)

	return a
}

func (a *App) Run() {
	p := tea.NewProgram(ui.InitialModel(a.client))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	fmt.Println("Connection Ended!")
}
