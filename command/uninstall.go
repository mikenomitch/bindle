package command

import (
	"fmt"
	"os"
	"strings"
)

type Uninstall struct{}

func (f *Uninstall) Help() string {
	helpText := `
Some helper text goes here
`
	return strings.TrimSpace(helpText)
}

func (f *Uninstall) Synopsis() string {
	return "Uninstall Bindle package. (Kind of... doesn't actually remove from nomad)"
}

func (f *Uninstall) Name() string { return "uninstall" }

func (f *Uninstall) Run(args []string) int {
	packageName := args[0]
	bindleDir := ".bindle"
	packageDir := fmt.Sprintf("%s/%s", bindleDir, packageName)

	// TODO: SOMEHOW IDENTIFY AND STOP THE JOB

	err := os.RemoveAll(packageDir)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error removing package %s", packageDir))
		return 1
	}

	return 0
}
