// Author: {{ .Author }}
// Created on: {{ .Timestamp }}
package main

import (
	"fmt"
	"net/http"

	"github.com/{{ .Author }}/{{ .ProjectName }}/internal/server"
)

func main() {
	srv := server.New()
	fmt.Println("Starting {{ .ProjectName }} server on :8080...")
	http.ListenAndServe(":8080", srv)
}
