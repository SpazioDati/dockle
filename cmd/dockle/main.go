package main

import (
	l "log"
	"os"

	"github.com/SpazioDati/dockle/pkg"
	"github.com/SpazioDati/dockle/pkg/log"
)

func main() {
	app := pkg.NewApp()
	err := app.Run(os.Args)

	if err != nil {
		if log.Logger != nil {
			log.Fatal(err)
		}
		l.Fatal(err)
	}
}
