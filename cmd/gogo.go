package main

import (
	"os"

	"github.com/dolab/gogo/cmd/commands"
	"github.com/golib/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gogo"
	app.Version = "1.0.0"
	app.Usage = "gogo COMMAND [ARGS]"

	app.Authors = []cli.Author{
		{
			Name:  "Spring MC",
			Email: "Heresy.MC@gmail.com",
		},
	}

	app.Commands = []cli.Command{
		commands.Application.Command(),
		commands.Component.Command(),
	}

	app.Run(os.Args)
}