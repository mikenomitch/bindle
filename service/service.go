package service

import (
	"log"
	"net/http"
)

func Run() {
	fs := http.FileServer(http.Dir("./frontend/dist/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/health", HandleHealth)
	http.HandleFunc("/package", HandlePackage)
	http.HandleFunc("/catalog", HandleCatalog)
	http.HandleFunc("/deploy", HandleDeploy)

	http.HandleFunc("/", HandleUI)
	log.Println("Listing for requests at http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
