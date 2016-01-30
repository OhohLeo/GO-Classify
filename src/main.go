package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "Classify"
	app.Version = "0.0.1"
	app.Usage = "Collections' classification tool"
	app.Action = func(c *cli.Context) {
		Start()
	}

	app.Run(os.Args)
}
