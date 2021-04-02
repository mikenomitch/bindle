package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/hashicorp/levant/template"
	"github.com/mikenomitch/bindle/utils"
)

type nomadResources struct {
	Resources []nomadResource `hcl:"nomad_resource,block"`
}

type nomadResource struct {
	Name         string `hcl:"key,label"`
	TemplateFile string `hcl:"template_file"`
	VariableFile string `hcl:"variable_file,optional"`
	Description  string `hcl:"description"`
	Type         string `hcl:"type"`
}

type Install struct{}

func (f *Install) Help() string {
	helpText := `
Usage: bindle install <package-name> [variables]

	Install Packages on your Nomad Cluster.

	Example: bindle install grafana -datacenters=us-east-1,us-east-2 -memory=500

	Assumes NOMAD_ADDR is configured properly and Nomad is up and running.

	Use "bindle list" to see which packages are availible.
`
	return strings.TrimSpace(helpText)
}

func (f *Install) Synopsis() string {
	return "Install Packages on your Nomad Cluster"
}

func (f *Install) Name() string { return "install" }

func (f *Install) Run(args []string) int {
	packageArg := args[0]

	fmt.Print("Installing Package: ", packageArg)

	packageName := packageArg
	catalogsDir := ".bindle/catalogs/default"
	if strings.Contains(packageName, "/") {
		splitBySlash := strings.Split(packageName, "/")
		packageName = splitBySlash[1]

		catalogsDir = ".bindle/catalogs/" + splitBySlash[0]
	}

	installsDir := ".bindle/installs"

	packageSourceDir := catalogsDir + "/" + packageName
	packageInstallDir := installsDir + "/" + packageArg

	manifestPath := packageSourceDir + "/manifest.hcl"
	topLevelVariablesPath := packageSourceDir + "/variables.hcl"

	err := os.RemoveAll(packageInstallDir)
	utils.Handle(err, "error removing old data")

	err = utils.Mkdir(packageInstallDir)
	utils.Handle(err, "error making installs dir for package")

	parser := hclparse.NewParser()
	manifestHCLFile, diags := parser.ParseHCLFile(manifestPath)

	res := nomadResources{}
	if diags = gohcl.DecodeBody(manifestHCLFile.Body, nil, &res); diags.HasErrors() {
		log.Printf(diags.Error())
		os.Exit(1)
	}

	for _, resource := range res.Resources {
		variablesPathForTemplate := topLevelVariablesPath
		if resource.VariableFile != "" {
			variablesPathForTemplate = packageSourceDir + "/" + resource.VariableFile
		}

		completedFilePath := packageInstallDir + "/" + resource.TemplateFile
		templatePath := packageSourceDir + "/" + resource.TemplateFile

		variableFilePaths := []string{topLevelVariablesPath, variablesPathForTemplate}

		cliVars := parseVariablesFromCliArgs(args)
		job, errorA := template.RenderTemplate(templatePath, variableFilePaths, "", &cliVars)
		utils.Handle(errorA, "Error rendering template")

		writer, err := os.OpenFile(completedFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		_, err = job.WriteTo(writer)
		utils.Handle(err, "Error writing template to disk")

		cmd := exec.Command("nomad", "run", completedFilePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}

	fmt.Print(fmt.Sprintf("Sent job to Nomad for package: \"%s\".\n\nBindle does not properly handle errors, so you will have to validate the deployment yourself.", packageName))

	return 1
}

// This is a pretty sad implementation
func parseVariablesFromCliArgs(args []string) map[string]string {
	vars := make(map[string]string)

	for i, arg := range args {
		if i > 0 && strings.Contains(arg, "=") && strings.HasPrefix(arg, "-") {
			splitByEquals := strings.Split(arg, "=")
			varNameWithExtra := splitByEquals[0]
			varVal := splitByEquals[1]
			varName := utils.TrimLeftChar(varNameWithExtra)
			vars[varName] = varVal
		}
	}

	return vars
}
