package command

import (
	"os"
	"strings"

	"github.com/mikenomitch/bindle/utils"
)

type Source struct{}

func (f *Source) Help() string {
	helpText := `
Usage: bindle source <catalog-name> <catalog-url>

	Adds a source catalog for Nomad packages.

	<catalog-name>: An identifier for the catalog. Example: "personal-packages"
	<catalog-url>: A git repo URL for the catalog. Example: "https://github.com/mikenomitch/nomad-packages"

	Re-calling source will re-fetch from the URL and replace the existing catalog.
`
	return strings.TrimSpace(helpText)
}

func (f *Source) Synopsis() string {
	return "Add a new source for a Nomad packages"
}

func (f *Source) Name() string { return "source" }

func (f *Source) Run(args []string) int {
	catalogsFilePath := ".bindle/catalogs"
	catalogName := args[0]
	catalogDir := catalogsFilePath + "/" + catalogName

	gitRepoUrl := args[1]

	err := os.RemoveAll(catalogDir)
	utils.Handle(err, "error removing old source data")

	err = utils.CloneRepoToDir(gitRepoUrl, catalogDir)
	utils.Handle(err, "error cloning catalog")

	utils.Log("Successfully added source \"" + catalogName + "\"")

	return 0
}
