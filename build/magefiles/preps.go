// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"log/slog"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Preps mg.Namespace

var generators = map[string]string{
	"moq":      "github.com/matryer/moq@latest",
	"stringer": "golang.org/x/tools/cmd/stringer@latest",
}

// Tidy runs go mod tidy to update the go.mod and go.sum files.
func (Preps) Tidy() error {
	slog.Info("Running go mod tidy...")

	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	return nil
}

// Format prettifies your code in a standard way to prevent arguments over curly braces.
func (Preps) Format() error {
	slog.Info("Running go fmt...")

	if err := sh.RunV("go", "fmt", "./..."); err != nil {
		return fmt.Errorf("failed to run go fmt: %w", err)
	}

	return nil
}

// Generate ensures all machine-generated files (gotext, stringer, moq, etc.) are up to date.
func (Preps) Generate() error {
	for tool, url := range generators {
		if err := foundOrInstalled(tool, url); err != nil {
			return fmt.Errorf("unable to install %s: %w", tool, err)
		}
	}

	slog.Info("Running go generate...")

	if err := sh.RunV("go", "generate", "-v", "./..."); err != nil {
		return fmt.Errorf("failed to run go generate: %w", err)
	}

	return nil
}
