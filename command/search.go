package command

import (
	"fmt"
	"strings"
)

type Search struct{}

func (f *Search) Help() string {
	helpText := `
Some helper text goes here
`
	return strings.TrimSpace(helpText)
}

func (f *Search) Synopsis() string {
	return "Search for Bindle packages"
}

func (f *Search) Name() string { return "search" }

func (f *Search) Run(args []string) int {
	fmt.Println("Search is Running - ignore for now")
	return 1
}
