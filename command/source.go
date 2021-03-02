package command

import (
	"fmt"
	"strings"

	"github.com/mikenomitch/bindle/utils"
)

type Source struct{}

func (f *Source) Help() string {
	helpText := `
Some helper text goes here
`
	return strings.TrimSpace(helpText)
}

func (f *Source) Synopsis() string {
	return "Add a new source for a Bindle package"
}

func (f *Source) Name() string { return "source" }

func (f *Source) Run(args []string) int {
	// TODO: Handle multiples and removals
	sourcesFilePath := ".bindle/sources"
	packageName := args[0]
	sourceForPackage := args[1]

	utils.WriteToFile(sourcesFilePath, fmt.Sprintf("%s,%s\n", packageName, sourceForPackage))
	fmt.Println("Saved new source for", packageName)

	return 1
}
