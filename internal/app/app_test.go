package app

import (
	"os"
	"testing"
)

func TestAppInit(t *testing.T) {
	os.Args = os.Args[4:]
	t.Log(os.Args)
	a := NewApp()
	t.Log(&a)
}
