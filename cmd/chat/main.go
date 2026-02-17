/*
Goal:
chat -c UserB
UserA> Yo!
UserB> Sup!
*/
package main

import (
	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/app"
)

func main() {
	App := app.NewApp()
	App.Run()
}
