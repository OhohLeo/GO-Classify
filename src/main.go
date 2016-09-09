package main

import (
	"github.com/codegangsta/cli"
	"github.com/ohohleo/classify/core"
	"log"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "Classify"
	app.Version = "0.0.1"
	app.Usage = "Collections' classification tool"
	app.Action = func(c *cli.Context) {
		classify, err := core.Start()
		if err != nil {
			log.Fatal(err.Error())
		}

		// Start server
		classify.Server.Start()
	}

	app.Run(os.Args)
}
