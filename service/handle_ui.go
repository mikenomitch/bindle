package service

import (
	"net/http"
)

func HandleUI(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "frontend/dist/index.html")
}
