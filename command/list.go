package command

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mikenomitch/bindle/utils"
)

type List struct{}

func (f *List) Help() string {
	helpText := `
Usage: bindle list

	Lists the packages availible to install.

	To add more packages, use the 'bindle souce' command.
`
	return strings.TrimSpace(helpText)
}

func (f *List) Synopsis() string {
	return "List the packages in your sources"
}

func (f *List) Name() string { return "list" }

func (f *List) Run(args []string) int {
	files, err := ioutil.ReadDir(".bindle/catalogs")
	utils.Handle(err, "Error reading catalogs")

	fmt.Println("Availible packages:\n ")

	for _, f := range files {
		if f.IsDir() {
			subFiles, _ := ioutil.ReadDir(".bindle/catalogs/" + f.Name())
			for _, sf := range subFiles {
				if sf.IsDir() && !(strings.HasPrefix(sf.Name(), ".")) {
					if f.Name() == "default" {
						fmt.Println(sf.Name())
					} else {
						fmt.Println(f.Name() + "/" + sf.Name())
					}
				}
			}
			fmt.Println("")
		}
	}

	fmt.Println("\nUse the 'source' command to add new catalogs")

	return 0
}
