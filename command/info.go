package command

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
)

type nomadVariables struct {
	Variables []nomadVarible `hcl:"variable,block"`
}

type nomadVarible struct {
	Key         string            `hcl:"key,label"`
	Default     string            `hcl:"default"`
	Description string            `hcl:"description,optional"`
	Type        string            `hcl:"type"`
	Meta        map[string]string `hcl:"meta,optional"`
}

type Info struct{}

func (f *Info) Help() string {
	helpText := `
Usage: bindle info <pack-name>

	Get information about a pack.

	Example: bindle info grafana
`
	return strings.TrimSpace(helpText)
}

func (f *Info) Synopsis() string {
	return "Get information about a pack"
}

func (f *Info) Name() string { return "info" }

func (f *Info) Run(args []string) int {
	packageArg := args[0]
	packageName := packageArg
	catalogsDir := ".bindle/catalogs/default"
	if strings.Contains(packageName, "/") {
		splitBySlash := strings.Split(packageName, "/")
		packageName = splitBySlash[1]

		catalogsDir = ".bindle/catalogs/" + splitBySlash[0]
	}

	packageSourceDir := catalogsDir + "/" + packageName
	topLevelVariablesPath := packageSourceDir + "/variables.hcl"

	parser := hclparse.NewParser()
	manifestHCLFile, diags := parser.ParseHCLFile(topLevelVariablesPath)

	res := nomadVariables{}
	if diags = gohcl.DecodeBody(manifestHCLFile.Body, nil, &res); diags.HasErrors() {
		log.Printf(diags.Error())
		os.Exit(1)
	}

	fmt.Println("Info for pack:", packageName)
	fmt.Println("")

	varString := ""
	fmt.Println("Variables:")
	for _, nomadVar := range res.Variables {
		fmt.Println("-----------")
		fmt.Println("Key: ", nomadVar.Key)
		fmt.Println("Type: ", nomadVar.Type)
		fmt.Println("Description: ", nomadVar.Description)
		fmt.Println("Default: ", nomadVar.Default)

		varString = varString + " -" + nomadVar.Key + "=" + nomadVar.Default
	}
	fmt.Println("=========")
	fmt.Println("")

	fmt.Println("Example installation:")
	fmt.Println("bindle install", packageName, varString)
	fmt.Println("")

	return 1
}
