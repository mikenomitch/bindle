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

type Vars struct{}

func (f *Vars) Help() string {
	helpText := `
Usage: bindle vars <pack-name>

	Get information on variables in a pack.

	Example: bindle vars grafana
`
	return strings.TrimSpace(helpText)
}

func (f *Vars) Synopsis() string {
	return "Get information on variables in a pack"
}

func (f *Vars) Name() string { return "vars" }

func (f *Vars) Run(args []string) int {
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

	fmt.Println("Variables for ", packageName)
	for _, nomadVar := range res.Variables {
		fmt.Println("-----------")
		fmt.Println("Key: ", nomadVar.Key)
		fmt.Println("Type: ", nomadVar.Type)
		fmt.Println("Description: ", nomadVar.Description)
		fmt.Println("Default: ", nomadVar.Default)
	}
	fmt.Println("=========")

	return 1
}
