package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Color
var StrtoColor map[string]string = map[string]string{
	"Red":    "#FF0000",
	"Green":  "$00FF00",
	"Blue":   "#0000FF",
	"Yellow": "#FFFF00",
	"Purple": "#800080",
	"Cyan":   "#48D1CC", // mediumturquoise
	"White":  "#FFFFFF",
}

type Config struct {
	Username string
	Socket   string
	Color    string
	Timeout  int64
}

func setDefaults(c *Config) {
	if c.Username == "" {
		var test string
		for test == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter your username: ")
			test, _ = reader.ReadString('\n')
			test = test[:len(test)-1]
		}
		c.Username = test
	}

	if c.Socket == "" {
		c.Socket = ":23456"
	}

	if c.Color == "" {
		c.Color = "White"
	}

	if c.Timeout == 0 {
		c.Timeout = 200
	}
}

func LoadConfig(opsys string) *Config {
	var (
		path string
		c    *Config = &Config{}
	)
	switch opsys {
	case "linux":
		home := os.Getenv("HOME")
		path = filepath.Join(home, ".config/chat")
		os.MkdirAll(path, 0755)
	}

	fconfig, err := os.OpenFile(filepath.Join(path, "config.json"), os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		panic("Couldn't Open the config file")
	}
	defer fconfig.Close()

	// {"Username":"JACK","Socket":":23456","Color":"Blue","Timeout":200}
	json.NewDecoder(fconfig).Decode(c)

	setDefaults(c)

	return c
}

func SaveConfig(c *Config, opsys string) {
	var (
		path string
	)
	switch opsys {
	case "linux":
		home := os.Getenv("HOME")
		path = filepath.Join(home, ".config/chat")
		os.MkdirAll(path, 0755)
	}

	fconfig, err := os.OpenFile(filepath.Join(path, "config.json"), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("Couldn't Open the config file")
	}
	defer fconfig.Close()

	json.NewEncoder(fconfig).Encode(c)
}
