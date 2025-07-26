package main

import (
	"github.com/nunseik/forma/cmd"
	"embed"
)

//go:embed all:templates
var templatesFS embed.FS

func main() {
	cmd.Execute(templatesFS)
}
