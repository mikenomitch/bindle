package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
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
	fmt.Println(args)

	packageName := args[0]

	bindleDir := ".bindle"
	packageDir := fmt.Sprintf("%s/%s", bindleDir, packageName)

	err := os.Mkdir(packageDir, 0755)
	if err != nil {
		fmt.Println(fmt.Sprintf("Package %s already downloaded.", packageDir))
		return 1
	}

	// TODO: Make the Base URL change-able based on the source files
	baseURL := "https://raw.githubusercontent.com/mikenomitch/nomad-packages/main"
	packageURL := fmt.Sprintf("%s/%s", baseURL, packageName)
	manifestURL := fmt.Sprintf("%s/%s", packageURL, "manifest")
	variablesURL := fmt.Sprintf("%s/%s", packageURL, "vars")

	bodyString, err := bodyFromURL(manifestURL)
	if err != nil {
		return 1
	}

	config := map[string]string{}
	configWithArgs, err := addCliArgsToConfig(config, args)
	variablesConfig, err := configFromVariableURL(variablesURL, configWithArgs)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	templatesToDownload := strings.Split(bodyString, "\n")

	log.Printf(templatesToDownload[0])

	for _, name := range templatesToDownload {
		fileURL := fmt.Sprintf("%s/%s", packageURL, name)

		fmt.Println(fileURL)

		templateFileBody, err := bodyFromURL(fileURL)
		if err != nil {
			return 1
		}

		fmt.Println(templateFileBody)

		completedFilePath := fmt.Sprintf("%s/%s", packageDir, name)

		err = parseTemplateAndWriteFile(
			completedFilePath,
			templateFileBody,
			variablesConfig,
		)
		if err != nil {
			return 1
		}

		// cmd := exec.Command("nomad", "run", completedFilePath)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// _ = cmd.Run()

		// Error status of two might have to be ignored?
		// if err != nil {
		// 	log.Print("error running job: ", err)
		// 	return 1
		// }
	}

	fmt.Println(fmt.Sprintf("Successfully installed %s", packageName))
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

func configFromVariableURL(url string, config map[string]string) (map[string]string, error) {
	varBodyString, _ := bodyFromURL(url)
	variableLines := strings.Split(varBodyString, "\n")

	for _, variableLine := range variableLines {
		linkChunks := strings.Split(variableLine, ",")
		variableName := linkChunks[0]

		if _, ok := config[variableName]; !ok {
			if len(linkChunks) > 1 {
				variableDefault := linkChunks[1]
				config[variableName] = variableDefault
			} else {
				err := fmt.Errorf("Missing value for %s", variableName)
				return config, err
			}
		}
	}

	return config, nil
}

func parseTemplateAndWriteFile(path string, templateBody string, config map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	tmpl, err := template.New(path).Parse(templateBody)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(file, config)

	return nil
}

func bodyFromURL(url string) (string, error) {
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
