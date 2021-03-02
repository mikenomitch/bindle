package command

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	levantCommand "github.com/hashicorp/levant/command"
	"github.com/mikenomitch/bindle/utils"
)

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
	log.Print("Installing", packageName)

	bindleDir := ".bindle"
	packageDir := bindleDir + "/" + packageName
	os.Mkdir(packageDir, 0755)

	topLevelVariablesPath := packageDir + "/variables.tf"

	baseURL := getBaseUrl(packageName)
	packageURL := baseURL + "/" + packageName
	manifestURL := packageURL + "/manifest"
	topLevelVariablesURL := packageURL + "/variables.tf"

	manifestBody, err := utils.BodyFromURL(manifestURL)
	if err != nil {
		return 1
	}

	err = utils.URLToFile(topLevelVariablesURL, topLevelVariablesPath)
	if err != nil {
		return 1
	}

	templatesToDownload := strings.Split(manifestBody, "\n")

	for _, name := range templatesToDownload {
		completedFilePath := packageDir + "/" + name
		templatePath := completedFilePath + ".template"
		variablesPath := packageDir + "/variables.tf"

		templateFileURL := packageURL + "/" + name
		variableURL := packageURL + "/variables.tf"

		err := utils.URLToFile(templateFileURL, templatePath)
		if err != nil {
			return 1
		}

		err = utils.URLToFile(variableURL, variablesPath)
		if err != nil {
			return 1
		}

		// TODO: ADD TOP LEVEL VARS FILE
		// TODO: ADD VARS FILE FROM THE FLAG

		c := levantCommand.RenderCommand{}
		args := []string{templatePath, "--out", completedFilePath, "--var-file", variablesPath}
		c.Run(args)

		// cmd := exec.Command("nomad", "run", completedFilePath)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// _ = cmd.Run()
	}

	log.Print(fmt.Sprintf("Successfully installed %s", packageName))
	return 1
}

func addCliArgsToConfig(config map[string]string, args []string) (map[string]string, error) {
	for i, arg := range args {
		if i > 0 {
			splitByEquals := strings.Split(arg, "=")
			varNameWithExtra := splitByEquals[0]
			varVal := splitByEquals[1]
			varName := trimLeftChar(varNameWithExtra)
			config[varName] = varVal
		}
	}

	return config, nil
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

func BodyFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}
