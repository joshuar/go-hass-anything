// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	_ "embed"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	appPathMatch = "../../apps/*/main.go"
	appFile      = "../../internal/agent/init.go"
)

//go:embed init.go.tmpl
var appTmpl string

//nolint:exhaustruct
func main() {
	templVars := struct {
		Dirs  []string
		Names []string
	}{}

	if len(os.Args) > 1 && os.Args[1] != "" {
		appPathMatch = filepath.Join(os.Args[1], "/*/main.go")
	}

	slog.Info("Looking for apps.",
		slog.String("path", appPathMatch))

	appMains, err := filepath.Glob(appPathMatch)
	if err != nil {
		slog.Error("Unable parse app pattern.",
			slog.Any("error", err))
		os.Exit(-1)
	}

	if appMains == nil {
		slog.Warn("No apps found.")
		os.Exit(0)
	}

	for _, main := range appMains {
		_, err = os.Stat(main)
		if err != nil {
			slog.Error("Could not read main.go for app. Ignoring app directory.",
				slog.String("directory", main),
				slog.Any("error", err))

			break
		}

		appDir := filepath.Dir(main)
		appDir = strings.ReplaceAll(appDir, "../", "")
		appName := filepath.Base(appDir)
		templVars.Dirs = append(templVars.Dirs, appDir)
		templVars.Names = append(templVars.Names, appName)
		slog.Info("Found app.",
			slog.String("name", appName),
			slog.String("directory", appDir))
	}

	initFile, err := os.Create(appFile)
	if err != nil {
		log.Fatalf("Unable to create file %s (%s). Exiting.", appFile, err.Error())
	}

	tmpl := template.Must(template.New("tmpl").Parse(appTmpl))

	if err := tmpl.Execute(initFile, templVars); err != nil {
		log.Fatalf("Unable to write out template to %s (%s). Exiting.", appFile, err.Error())
	}

	cmd := exec.Command("gofmt", "-w", appFile)
	if err := cmd.Run(); err != nil {
		log.Printf("Unable to run gofmt on %s. Formatting may be wonky.", appFile)
	}
}
