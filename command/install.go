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
Some helper text goes here
`
	return strings.TrimSpace(helpText)
}

func (f *Install) Synopsis() string {
	return "Install Nomad jobs automatically"
}

func (f *Install) Name() string { return "install" }

type installFlags map[string]string

func (i *installFlags) String() string {
	return "the return value!"
}

func (i *installFlags) Set(value string) error {
	utils.Log("IN HERE: " + value)

	valSplitByEquals := strings.Split(value, "=")
	i.Get()[valSplitByEquals[0]] = valSplitByEquals[1]

	return nil
}

func (i *installFlags) Get() map[string]string {
	return *i
}

func (f *Install) Run(args []string) int {
	packageName := args[0]
	log.Print("Installing Package: ", packageName)

	catalogsDir := ".bindle/catalogs/default"
	installsDir := ".bindle/installs"

	packageSourceDir := catalogsDir + "/" + packageName
	packageInstallDir := installsDir + "/" + packageName

	manifestPath := packageSourceDir + "/manifest.hcl"
	topLevelVariablesPath := packageSourceDir + "/variables.hcl"

	err := utils.Mkdir(packageInstallDir)
	utils.Handle(err, "error making installs dir for package")

	argVariablesPath := packageInstallDir + "/arg_variables.json"
	utils.CreateEmptyFile(argVariablesPath)

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

		if errorA != nil {
			log.Printf("error rendering template: %s", err)
			return 1
		}

		writer, err := os.OpenFile(completedFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		_, err = job.WriteTo(writer)
		if err != nil {
			return 1
		}

		cmd := exec.Command("nomad", "run", completedFilePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}

	log.Print(fmt.Sprintf("Successfully installed %s", packageName))
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
