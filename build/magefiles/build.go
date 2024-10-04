// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

var ErrBuildFailed = errors.New("build failed")

// Full runs all prep steps and then builds the binary.
func (Build) Full() error {
	slog.Info("Starting full build.")

	// Make everything nice, neat, and proper
	mg.Deps(Preps.Generate)
	mg.Deps(Preps.Tidy)
	mg.Deps(Preps.Format)

	return buildProject()
}

// Fast just builds the binary and does not run any prep steps. It will fail if
// the prep steps have not run.
func (Build) Fast() error {
	return buildProject()
}

func (b Build) CI() error {
	if !isCI() {
		return ErrNotCI
	}

	mg.Deps(b.Full)

	return nil
}

//nolint:mnd
func buildProject() error {
	if err := os.RemoveAll(distPath); err != nil {
		return fmt.Errorf("could not clean dist directory: %w", err)
	}

	if err := os.Mkdir(distPath, 0o755); err != nil {
		return fmt.Errorf("could not create dist directory: %w", err)
	}

	// Set-up the build environment.
	buildEnv, err := generateBuildEnv()
	if err != nil {
		return errors.Join(ErrBuildFailed, err)
	}

	ldflags, err := getFlags()
	if err != nil {
		return errors.Join(ErrBuildFailed, err)
	}

	slog.Info("Running go build...",
		slog.String("output", buildEnv["OUTPUT"]),
		slog.String("ldflags", ldflags))

	if err := sh.RunV("go", "build", "-ldflags="+ldflags, "-o", buildEnv["OUTPUT"]); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}
