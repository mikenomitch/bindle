package command

import (
	"fmt"
	"strings"

	"github.com/mikenomitch/bindle/service"
	"github.com/mikenomitch/bindle/utils"
	"github.com/skratchdot/open-golang/open"
)

type UI struct{}

func (f *UI) Help() string {
	helpText := `
Usage: bindle ui

	Launches a user interface for bindle.
`
	return strings.TrimSpace(helpText)
}

func (f *UI) Synopsis() string {
	return "Launches a user interface for bindle"
}

func (f *UI) Name() string { return "ui" }

func (f *UI) Run(args []string) int {
	if err := open.Start("http://localhost:9000"); err != nil {
		utils.Handle(err, "Error opening URL")
		return 1
	}

	fmt.Println("Starting Bindle server.")
	service.Run()

	return 0
}
