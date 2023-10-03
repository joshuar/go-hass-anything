// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	_ "embed"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var appPathMatch = "../../internal/apps/*"
var appFile = "../../internal/agent/init.go"

//go:embed init.go.tmpl
var appTmpl string

func main() {
	var apps []string
	var err error

	a, err := filepath.Glob(appPathMatch)
	if err != nil {
		log.Fatalf("Unable to match apps with %s (%s). Exiting.", appPathMatch, err.Error())
		os.Exit(-1)
	}

	for _, v := range a {
		app := filepath.Base(v)
		log.Printf("checking %s...", filepath.Join(v, app+".go"))
		if _, err = os.Stat(filepath.Join(v, app+".go")); !os.IsNotExist(err) {
			log.Printf("Found app %s.", app)
			apps = append(apps, app)
		}
	}

	t := template.Must(template.New("tmpl").Parse(appTmpl))

	f, err := os.Create(appFile)
	if err != nil {
		log.Fatalf("Unable to create file %s (%s). Exiting.", appFile, err.Error())
	}

	if err := t.Execute(f, apps); err != nil {
		log.Fatalf("Unable to write out template to %s (%s). Exiting.", appFile, err.Error())
	}

	cmd := exec.Command("gofmt", "-w", appFile)
	if err := cmd.Run(); err != nil {
		log.Printf("Unable to run gofmt on %s. Formatting may be wonky.", appFile)
	}
}