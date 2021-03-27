package command

import (
	"fmt"
	"os"
	"strings"
)

type Uninstall struct{}

func (f *Uninstall) Help() string {
	helpText := `
Usage: bindle init package-name

	Currently only removes the templates locally.

	Will ideally be extended to stop running jobs.
`
	return strings.TrimSpace(helpText)
}

func (f *Uninstall) Synopsis() string {
	return "Uninstall a Nomad Package and stop running related jobs. (Not working yet)"
}

func (f *Uninstall) Name() string { return "uninstall" }

func (f *Uninstall) Run(args []string) int {
	packageName := args[0]
	bindleDir := ".bindle/installs"
	packageDir := fmt.Sprintf("%s/%s", bindleDir, packageName)

	// TODO: SOMEHOW IDENTIFY AND STOP THE JOB

	err := os.RemoveAll(packageDir)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error removing package %s", packageDir))
		return 1
	}

	return 0
}
