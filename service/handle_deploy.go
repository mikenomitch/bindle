package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type DeployRequest struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
}

func HandleDeploy(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if (*req).Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(req.Body)

	var dr DeployRequest

	err := decoder.Decode(&dr)

	if err != nil {
		panic(err)
	}

	variableArgs := []string{"install", dr.Name}

	for k, v := range dr.Variables {
		flag := fmt.Sprintf("-%s=%s", k, v)
		variableArgs = append(variableArgs, flag)
	}

	cmd := exec.Command("bindle", variableArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
