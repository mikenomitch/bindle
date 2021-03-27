package command

import (
	"os"
	"strings"

	"github.com/mikenomitch/bindle/utils"
)

type Init struct{}

func (f *Init) Help() string {
	helpText := `
Usage: bindle init

	Init creates a .bindle directory and pulls the default catalog of Nomad packages.

	Calling init again will reload the default catalog.
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

	err := utils.Mkdir(bindleDir)
	utils.Handle(err, "error initializing")
	err = utils.Mkdir(catalogsDir)
	utils.Handle(err, "error initializing")

	err = utils.Mkdir(installsDir)
	utils.Handle(err, "error initializing")

	err = utils.Mkdir(defaultCatalogSourceDir)
	utils.Handle(err, "error initializing")

	err = os.RemoveAll(defaultCatalogSourceDir)
	utils.Handle(err, "error removing old data")

	err = utils.CloneRepoToDir(defaultCatalogRepo, defaultCatalogSourceDir)
	utils.Handle(err, "error cloning default catalog")

	utils.Log("Bindle successfully initialized.")

	return 0
}
