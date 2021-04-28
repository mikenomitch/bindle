package service

import (
	"log"
	"net/http"
)

func Run() {

	http.HandleFunc("/", HandleUI)
	http.HandleFunc("/health", HandleHealth)

	http.HandleFunc("/package", HandlePackage)
	http.HandleFunc("/catalog", HandleCatalog)

	http.HandleFunc("/deploy", HandleDeploy)

	log.Println("Listing for requests at http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
