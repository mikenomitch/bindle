package service

import (
	"encoding/json"
	"net/http"

	"github.com/mikenomitch/bindle/utils"
)

type response struct {
	Packages []string `json:"packages"`
}

const catalogPath = ".bindle/catalogs/default"

func HandleCatalog(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Content-Type", "application/json")

	packages := utils.DirsInPath(catalogPath)

	json.NewEncoder(w).Encode(packages)
}
