package app

import (
	"fmt"
	"log"
	"os"

	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/client"
	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/command"
	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	user client.User
}

func NewApp() (a App) {
	var (
		args_listen struct {
			Username string `command:"username"`
		}
		args_connect struct {
			Username string `command:"username"`
			Host     string `command:"host"`
		}

		args_option struct {
			Color client.Color `command:"color"`
		}
	)

	arg_p := command.NewParser(os.Args[1:])
	arg_p.AddCmd("listen", "l", "Listen for incoming connections", []string{"username"}, &args_listen, func() {
		config := client.Config{Name: args_listen.Username, Socket: ":23456"}

		var err error
		a.user, err = client.MakeUser("listen", config)
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	arg_p.AddCmd("connect", "c", "Connect to the given user", []string{"host", "username"}, &args_connect, func() {
		config := client.Config{Name: args_listen.Username, Socket: ":23456"}

		var err error
		a.user, err = client.MakeUser("connect", config, args_connect.Host)
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	arg_p.AddOption("color", "c", "Set the color of your messages", []string{"color"}, &args_option)
	c, err := arg_p.Parse()
	if err != nil {
		log.Fatal(err.Error())
	}
	c()

	return a
}

func (a App) Run() {
	p := tea.NewProgram(ui.InitialModel(a.user))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	fmt.Println("Connection Ended!")
}
