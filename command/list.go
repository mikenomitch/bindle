package command

import (
	"fmt"
	"strings"
)

type List struct{}

func (f *List) Help() string {
	helpText := `
Some helper text goes here
`
	return strings.TrimSpace(helpText)
}

func (f *List) Synopsis() string {
	return "List the packages in your sources"
}

func (f *List) Name() string { return "list" }

func (f *List) Run(args []string) int {
	fmt.Println("List is Running - ignore for now")
	return 1
}
