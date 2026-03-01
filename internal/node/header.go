package node

import "strings"

type Header struct {
	Username string
}

func build_header(h Header) string {
	var s strings.Builder
	s.WriteString(h.Username + " ")
	return s.String()
}

// header structure is already known so no funny buisness with reflection
func parse_header(msg string) (*Header, int) {
	var (
		h Header
		i int
	)
	get := delim(msg, ' ')

	h.Username, i = get()

	return &h, i
}

func delim(str string, delim byte) func() (string, int) {
	i := 0
	return func() (string, int) {
		start := i
		if str[start] == delim {
			i++
		}

		for ; str[i] != delim; i++ {
		}
		return str[start:i], i + 1
	}
}
