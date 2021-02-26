package main

import (
	"log"
	"os"

	cmd "github.com/mikenomitch/bindle/command"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("bindle", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"install": func() (cli.Command, error) {
			return &cmd.Install{}, nil
		},
		"uninstall": func() (cli.Command, error) {
			return &cmd.Uninstall{}, nil
		},
		"search": func() (cli.Command, error) {
			return &cmd.Search{}, nil
		},
		"init": func() (cli.Command, error) {
			return &cmd.Init{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
	// parse the arguments and pass to command
	// fall back to some message
}
