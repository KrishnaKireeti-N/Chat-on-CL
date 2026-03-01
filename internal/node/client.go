// TODO: Implement Proper form for sharing data
package node

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net"
	"time"
)

type (
	Client struct {
		Name        string
		Senderstyle string
		Conn        net.Conn
		Users       map[string]*User
		timeout     int64
	}
	User struct {
		Name        string
		Senderstyle string
	}
	Message struct {
		header string
		body   string
	}
)

var (
	config *Config
)

func MakeUser(mode string, c *Config, v ...string) (*Client, error) {
	config = c
	fmt.Println("In MakeUser: ", config)
	modes := map[string]func(v ...string) (*Client, error){"listen": makeUserListen, "connect": makeUserConnect}
	for m, f := range modes {
		if m == mode {
			return f(v...)
		}
	}

	return nil, errors.New("Mode mismatch\n\tAvailable: \"listen\" \"connect\"")
}

func makeUserListen(v ...string) (*Client, error) {
	fmt.Println("localhost" + config.Socket)
	listener, err := net.Listen("tcp", "localhost"+config.Socket)
	if err != nil {
		return nil, fmt.Errorf("Coudln't listen on socket %s", config.Socket)
	}

	conn, err := listener.Accept()
	if err != nil {
		return nil, errors.New("Couldn't accept a connection:\n" + err.Error())
	}

	// SYN-ACK
	var data []byte = make([]byte, 256)
	n, err := conn.Read(data)
	data = data[:n]
	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	fmt.Println("LISTEN ACK", records, n)

	var buf []byte
	buf = fmt.Appendf(buf, "%v,%v", config.Username, config.Color)
	fmt.Println("LISTEN SYN", string(buf))
	conn.Write(buf)

	user := User{Name: records[0][0], Senderstyle: StrtoColor[records[0][1]]}
	c := &Client{
		Name:        config.Username,
		Senderstyle: StrtoColor[config.Color],
		Conn:        conn,
		Users:       map[string]*User{user.Name: &user},
		timeout:     config.Timeout,
	}

	return c, nil
}

func makeUserConnect(v ...string) (*Client, error) {
	conn, err := net.Dial("tcp", v[0]+config.Socket)
	if err != nil {
		return nil, errors.New("Couldn't establish a connection:\n" + err.Error())
	}

	// SYN
	var buf []byte
	buf = fmt.Appendf(buf, "%v,%v", config.Username, config.Color)
	fmt.Println("Connect SYN", string(buf))
	conn.Write(buf)

	// ACK
	var data []byte = make([]byte, 256)
	n, err := conn.Read(data)
	data = data[:n]
	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	fmt.Println("CONNECT ACK", records, n)

	user := User{Name: records[0][0], Senderstyle: StrtoColor[records[0][1]]}
	c := &Client{
		Name:        config.Username,
		Senderstyle: StrtoColor[config.Color],
		Conn:        conn,
		Users:       map[string]*User{user.Name: &user},
		timeout:     config.Timeout,
	}

	return c, nil
}

func (c *Client) Send(msg string) {
	h := Header{
		Username: c.Name,
	}
	header := build_header(h)
	c.Conn.Write([]byte(header + msg))
}

func (c *Client) SendTimeout(msg string) {
	c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(c.timeout) * time.Second))
	h := Header{
		Username: c.Name,
	}
	header := build_header(h)
	c.Conn.Write([]byte(header + msg))
	c.Conn.SetWriteDeadline(time.Time{})
}

func (c *Client) Recieve() (*Header, string) {
	buf := make([]byte, 256)
	n, err := c.Conn.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
	}
	msg := string(buf[:n])
	if msg == "" {
		return nil, "\r"
	}

	h, n := parse_header(msg)

	msg = msg[n:]
	switch msg {
	case "\r":
		c.Conn.Close()
		return h, "\r"
	case "":
		return nil, ""
	default:
		return h, msg
	}
}
