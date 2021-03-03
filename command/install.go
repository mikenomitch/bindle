package command

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	VariableFile string `hcl:"variable_file"`
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

func (f *Install) Run(args []string) int {
	packageName := args[0]
	log.Print("Installing ", packageName)

	bindleDir := ".bindle"
	packageDir := bindleDir + "/" + packageName

	// NOTE: FOR NOW JUST WIPE IT ALL CLEAN
	err := os.RemoveAll(packageDir)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error removing package %s", packageDir))
		return 1
	}

	os.Mkdir(packageDir, 0755)

	manifestPath := packageDir + "/manifest.hcl"
	topLevelVariablesPath := packageDir + "/variables.tf"
	argVariablesPath := packageDir + "/arg_variables.json"

	baseURL := getBaseUrl(packageName)
	packageURL := baseURL + "/" + packageName
	manifestURL := packageURL + "/manifest.hcl"
	topLevelVariablesURL := packageURL + "/variables.tf"

	err = utils.URLToFile(manifestURL, manifestPath)
	if err != nil {
		return 1
	}

	err = utils.URLToFile(topLevelVariablesURL, topLevelVariablesPath)
	if err != nil {
		return 1
	}

	err = writeCLIArgsToFile(args, argVariablesPath)
	if err != nil {
		return 1
	}

	parser := hclparse.NewParser()
	manifestHCLFile, diags := parser.ParseHCLFile(manifestPath)

	res := nomadResources{}
	if diags = gohcl.DecodeBody(manifestHCLFile.Body, nil, &res); diags.HasErrors() {
		log.Printf(diags.Error())
		return 1
	}

	for _, resource := range res.Resources {
		completedFilePath := packageDir + "/" + resource.TemplateFile
		templatePath := completedFilePath + ".template"
		variablesPath := packageDir + "/" + resource.VariableFile

		templateFileURL := packageURL + "/" + resource.TemplateFile
		variableURL := packageURL + "/" + resource.VariableFile

		err := utils.URLToFile(templateFileURL, templatePath)
		if err != nil {
			return 1
		}

		err = utils.URLToFile(variableURL, variablesPath)
		if err != nil {
			return 1
		}

		vars := make(map[string]string)
		vars["job_name"] = resource.Name

		variableFilePaths := []string{topLevelVariablesPath, variablesPath, argVariablesPath}
		job, errorA := template.RenderTemplate(templatePath, variableFilePaths, "", &vars)
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

func writeCLIArgsToFile(args []string, path string) error {
	vars := make(map[string]string)
	for i, arg := range args {
		if i > 0 && strings.Contains(arg, "=") {
			splitByEquals := strings.Split(arg, "=")
			varNameWithExtra := splitByEquals[0]
			varVal := splitByEquals[1]
			varName := utils.TrimLeftChar(varNameWithExtra)
			vars[varName] = varVal
		}
	}

	file, _ := json.MarshalIndent(vars, "", " ")

	err := ioutil.WriteFile(path, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// func configFromVariableURL(url string, config map[string]string) (map[string]string, error) {
// 	varBodyString, _ := utils.BodyFromURL(url)
// 	variableLines := strings.Split(varBodyString, "\n")

// 	for _, variableLine := range variableLines {
// 		linkChunks := strings.Split(variableLine, ",")
// 		variableName := linkChunks[0]

// 		if _, ok := config[variableName]; !ok {
// 			if len(linkChunks) > 1 {
// 				variableDefault := linkChunks[1]
// 				config[variableName] = variableDefault
// 			} else {
// 				err := fmt.Errorf("Missing value for %s", variableName)
// 				return config, err
// 			}
// 		}
// 	}

// 	return config, nil
// }

// func parseTemplateAndWriteFile(path string, templateBody string, config map[string]string) error {
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}

// 	tmpl, err := template.New(path).Parse(templateBody)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = tmpl.Execute(file, config)

// 	return nil
// }

func getBaseUrl(packageName string) string {
	sourcesFilePath := ".bindle/sources"
	// set the default
	baseURL := "https://raw.githubusercontent.com/mikenomitch/nomad-packages/main"

	file, err := os.Open(sourcesFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		prefix := fmt.Sprintf("%s,", packageName)

		if strings.HasPrefix(line, prefix) {
			byComma := strings.Split(line, ",")
			baseURL = byComma[1]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return baseURL
}
