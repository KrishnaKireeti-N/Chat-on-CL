// TODO: Implement Proper form for sharing data
package client

import (
	"errors"
	"fmt"
	"net"
)

type (
	User struct {
		Name        string
		Senderstyle Color
		Conn        net.Conn
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

var (
	config Config
)

func MakeUser(mode string, c Config, v ...string) (User, error) {
	config = c
	modes := map[string]func(v ...string) (User, error){"listen": makeUserListen, "connect": makeUserConnect}
	for m, f := range modes {
		if m == mode {
			return f(v...)
		}
	}

	return User{}, errors.New("Mode mismatch\n\tAvailable: \"listen\" \"connect\"")
}

func makeUserListen(v ...string) (User, error) {
	fmt.Println("localhost" + config.Socket)
	listener, err := net.Listen("tcp", "localhost"+config.Socket)
	if err != nil {
		return User{}, fmt.Errorf("Coudln't listen on socket %s", config.Socket)
	}

	conn, err := listener.Accept()
	if err != nil {
		return User{}, errors.New("Couldn't accept a connection:\n" + err.Error())
	}

	return User{Name: config.Name, Senderstyle: Blue, Conn: conn}, nil
}

func makeUserConnect(v ...string) (User, error) {
	conn, err := net.Dial("tcp", v[0]+config.Socket)
	if err != nil {
		return User{}, errors.New("Couldn't establish a connection:\n" + err.Error())
	}

	return User{Name: config.Name, Senderstyle: Red, Conn: conn}, nil
}

// TODO: Handshake for establsihing any prior data
func (u User) Handshake() Data {
	m := Data{}

	return m
}

func (c User) Send(msg string) {
	c.Conn.Write([]byte(msg))
}

func (c User) Recieve() string {
	buf := make([]byte, 256)
	n, err := c.Conn.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
	}

	msg := string(buf[:n])
	if msg == "0" {
		c.Conn.Close()
		return "0"
	}

	return msg
}
