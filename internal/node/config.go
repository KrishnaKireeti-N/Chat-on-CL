package node

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Username string
	Socket   string
	Color    string
	Timeout  int64
}

func (c *Config) SetDefaults() {
	if c.Username == "" {
		var test string
		for test == "" {
			fmt.Print("Enter your username: ")
			fmt.Scanf("%v", &test)
			strings.TrimSuffix(test, "\n")
			strings.TrimSuffix(test, "\r")
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

func getUserPath() string {
	path, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	path = filepath.Join(path, "chat")

	os.MkdirAll(path, 0755)

	return path
}

func LoadConfig() *Config {
	var (
		path string  = getUserPath()
		c    *Config = &Config{}
	)

	fconfig, err := os.OpenFile(filepath.Join(path, "config.json"), os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fconfig.Close()

	// {"Username":"JACK","Socket":":23456","Color":"Blue","Timeout":200}
	json.NewDecoder(fconfig).Decode(c)

	return c
}

func SaveConfig(c *Config) {
	path := getUserPath()
	fconfig, err := os.OpenFile(filepath.Join(path, "config.json"), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fconfig.Close()

	if err := json.NewEncoder(fconfig).Encode(c); err != nil {
		panic(err)
	}
}
