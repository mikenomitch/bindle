package command

import (
	"os"
	"strings"

	"github.com/mikenomitch/bindle/utils"
)

type Init struct{}

func (f *Init) Help() string {
	helpText := `
Init sets up bindle and the associated terraform resources.

It creates a .bindle directory and pulls the base source.
`
	return strings.TrimSpace(helpText)
}

func (f *Init) Synopsis() string {
	return "Initializes .bindle directory and necessary files"
}

func (f *Init) Name() string { return "init" }

func (f *Init) Run(args []string) int {
	bindleDir := ".bindle"
	catalogsDir := bindleDir + "/catalogs"
	installsDir := bindleDir + "/installs"

	defaultCatalogRepo := "https://github.com/mikenomitch/nomad-packages"
	defaultCatalogSourceDir := catalogsDir + "/default"

	// TODO: make this not remove sources eventaully
	err := os.RemoveAll(bindleDir)
	utils.Handle(err, "error removing old data")

	err = os.Mkdir(bindleDir, 0755)
	utils.Handle(err, "error initializing")
	err = os.Mkdir(catalogsDir, 0755)
	utils.Handle(err, "error initializing")

	err = os.Mkdir(installsDir, 0755)
	utils.Handle(err, "error initializing")

	err = os.Mkdir(defaultCatalogSourceDir, 0755)
	utils.Handle(err, "error initializing")

	err = utils.CloneRepoToDir(defaultCatalogRepo, defaultCatalogSourceDir)
	utils.Handle(err, "error cloning default catalog")

	utils.Log("Bindle successfully initialized.")

	return 0
}
