package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
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
	return "Initializes .bindle directory, files, and subdirs"
}

func (f *Init) Name() string { return "init" }

func (f *Init) Run(args []string) int {
	fmt.Println("Adding a .bindle directory")
	err := os.Mkdir(".bindle", 0755)
	if err != nil {
		fmt.Println("Bindle already initialized.")
		return 1
	}

	fmt.Println("Adding a sources file")
	createEmptyFile(".bindle/sources")

	// ==== Write the initial TF configuration ===
	mainFile, err := os.Create(".bindle/main.tf")
	if err != nil {
		log.Println("create file: ", err)
		return 1
	}

	t, err := template.ParseFiles("./templates/main.tf")
	if err != nil {
		log.Print(err)
		return 1
	}

	config := map[string]string{
		"Address": "127.0.0.1:4646",
	}

	err = t.Execute(mainFile, config)
	if err != nil {
		log.Print("execute: ", err)
		return 1
	}

	fmt.Println("Run terraform init in the .bindle directory")
	cmd := exec.Command("terraform", "init")
	os.Chdir("./.bindle")

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Print("error initializing terraform: ", err)
		return 1
	}

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
