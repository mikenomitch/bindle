package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Init struct{}

func (f *Init) Help() string {
	helpText := `
Init sets up bindle and the associated terraform resources.

It creates a .bindle directory, passes Nomad configuration variables into terraform configuration, initializes terraform.
`
	return strings.TrimSpace(helpText)
}

func (f *Init) Synopsis() string {
	return "Initializes .bindle directory and necessary files"
}

func (f *Init) Name() string { return "init" }

func (f *Init) Run(args []string) int {
	err := os.Mkdir(".bindle", 0755)
	if err != nil {
		fmt.Println("Bindle already initialized.")
		return 1
	}
	createEmptyFile(".bindle/sources")
	log.Println("Bindle successfully initialized.")

	return 0
}

// ======= INTERNAL =======

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createEmptyFile(name string) {
	d := []byte("")
	check(ioutil.WriteFile(name, d, 0644))
}
