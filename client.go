// TODO: Implement Proper form for sharing data
package main

import (
	"errors"
	"fmt"
	"net"
)

type (
	User struct {
		name        string
		senderstyle Color
		conn        net.Conn
	}
	Data struct {
		name        string
		senderstyle Color
	}
	Message struct {
		header string
		body   string
	}
)

func makeUser(mode string, v ...string) (User, error) {
	modes := map[string]func(v ...string) (User, error){"listen": makeUserListen, "connect": makeUserConnect}
	for m, f := range modes {
		if m == mode {
			return f(v...)
		}
	}

	return User{}, errors.New("Mode mismatch\n\tAvailable: \"listen\" \"connect\"")
}

func makeUserListen(v ...string) (User, error) {
	if len(v) != 1 {
		return User{}, errors.New("Provide proper arguments (mode, name)")
	}

	listener, err := net.Listen("tcp", config.Socket)
	if err != nil {
		return User{}, fmt.Errorf("Coudln't listen on socket %s", config.Socket)
	}

	conn, err := listener.Accept()
	if err != nil {
		return User{}, errors.New("Couldn't accept a connection:\n" + err.Error())
	}

	return User{name: v[0], senderstyle: Blue, conn: conn}, nil
}

func makeUserConnect(v ...string) (User, error) {
	if len(v) != 2 {
		return User{}, errors.New("Provide proper arguments (mode, name, endip)")
	}

	conn, err := net.Dial("tcp", v[1])
	if err != nil {
		return User{}, errors.New("Couldn't establish a connection:\n" + err.Error())
	}

	return User{name: v[0], senderstyle: Red, conn: conn}, nil
}

// TODO: Handshake for establsihing any prior data
func (u User) Handshake() Data {
	m := Data{}

	return m
}

func (c User) Send(msg string) {
	c.conn.Write([]byte(msg))
}

func (c User) Recieve() string {
	buf := make([]byte, 256)
	n, err := c.conn.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
	}

	msg := string(buf[:n])
	if msg == "0" {
		c.conn.Close()
		return "0"
	}

	return msg
}
