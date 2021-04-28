package service

import (
	"io"
	"net/http"
)

func HandleHealth(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "ok\n")
}
