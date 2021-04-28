package service

import (
	"io"
	"net/http"
)

func HandleUI(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "figure out how to serve the frontend/dist/index.html file here")
}
