package main

import (
	"fmt"
	"log"
	"os"
)

// settings
var config = struct {
	Socket string
}{":23456"}

// global variables
var (
	user User
)

func init() {
	if len(os.Args) < 2 {
		fmt.Printf("Use %v -h for help", os.Args[0])
		os.Exit(1)
	}

	i := 1
	switch os.Args[i] {
	case "-h":
		fmt.Printf(
			`%v [OPTION] <args>
	-s <username>:
		Listen for connections
	-c <username> <ip:port/name>:
		Connect to a certain 'address' and start chatting
	-h:
		Show this help page
`,
			os.Args[0],
		)
		os.Exit(0)
	case "-s":
		i++
		username := os.Args[i]

		var err error
		user, err = makeUser("listen", username)
		if err != nil {
			log.Fatal(err.Error())
		}
	case "-c":
		i++
		username := os.Args[i]
		i++
		endip := os.Args[i]

		var err error
		user, err = makeUser("connect", username, endip+config.Socket)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
