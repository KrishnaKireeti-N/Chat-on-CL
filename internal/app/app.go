package app

import (
	"fmt"
	"log"
	"os"

	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/command"
	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/node"
	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	client *node.Client
}

func NewApp() *App {
	var (
		// Commands
		args_listen struct {
		}
		args_connect struct {
			Host string `command:"host"`
		}

		// Options
		args_option struct {
			Color string `command:"color"`
		}

		// App
		a App
	)

	// linux
	config := node.LoadConfig("linux")
	defer node.SaveConfig(config, "linux")

	arg_p := command.NewParser(os.Args[1:], 2, 1)
	arg_p.AddCmd("listen", "l", "Listen for incoming connections", []string{}, &args_listen, func() {
		var (
			err error
		)
		a.client, err = node.MakeUser("listen", config)
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	arg_p.AddCmd("connect", "c", "Connect to the given user", []string{"host"}, &args_connect, func() {
		var (
			err error
		)
		a.client, err = node.MakeUser("connect", config, args_connect.Host)
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	arg_p.AddOption("color", "c", "Set the color of your messages", []string{"color"}, &args_option)

	c, err := arg_p.Parse()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Managing the config according to the options provided
	if color := args_option.Color; color != "" {
		config.Color = color
	}

	fmt.Println(config)

	c()

	return &a
}

func (a *App) Run() {
	p := tea.NewProgram(ui.InitialModel(a.client))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	fmt.Println("Connection Ended!")
}
