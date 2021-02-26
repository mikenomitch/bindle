package command

import (
	"fmt"
	"strings"
)

type Source struct{}

func (f *Source) Help() string {
	helpText := `
Some helper text goes here
`
	return strings.TrimSpace(helpText)
}

func (f *Source) Synopsis() string {
	return "Add a new source for Bindle packages"
}

func (f *Source) Name() string { return "source" }

func (f *Source) Run(args []string) int {
	fmt.Println("Running Source")
	fmt.Println("Parse args")
	fmt.Println("If first arg is 'remove' then remove it from the file - short circuit on this")
	fmt.Println("First arg - Name, Second arg - URL")
	fmt.Println("Write these to a local file")
	fmt.Println("Return a success message")

	return 1
}
