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
	"os/exec"
	"slices"
	"strings"
	"syscall"

	"github.com/magefile/mage/sh"
)

const (
	pkgBase = "github.com/joshuar/go-hass-anything/pkg/preferences"
)

var ErrNotCI = errors.New("not in CI environment")

// isCI checks whether we are currently running as part of a CI pipeline (i.e.
// in a GitHub runner).
func isCI() bool {
	return os.Getenv("CI") != ""
}

// isRoot checks whether we are running as the root user or with elevated
// privileges.
func isRoot() bool {
	euid := syscall.Geteuid()
	uid := syscall.Getuid()
	egid := syscall.Getegid()
	gid := syscall.Getgid()
	if uid != euid || gid != egid || uid == 0 {
		return true
	}
	return false
}

// SudoWrap will "wrap" the given command with sudo if needed.
func SudoWrap(cmd string, args ...string) error {
	if isRoot() {
		return sh.RunV(cmd, args...)
	}
	return sh.RunV("sudo", slices.Concat([]string{cmd}, args)...)
}

// FoundOrInstalled checks for existence then installs a file if it's not there.
func FoundOrInstalled(executableName, installURL string) (isInstalled bool) {
	_, missing := exec.LookPath(executableName)
	if missing != nil {
		slog.Info("Installing tool.", "tool", executableName, "url", installURL)
		err := sh.Run("go", "install", installURL)
		if err != nil {
			return false
		}
	}
	return true
}

// GetFlags gets all the compile flags to set the version and stuff.
func GetFlags() (string, error) {
	var version, hash, date string
	var err error
	if version, err = Version(); err != nil {
		return "", fmt.Errorf("failed to retrieve version from git: %w", err)
	}
	if hash, err = GitHash(); err != nil {
		return "", fmt.Errorf("failed to retrieve hash from git: %w", err)
	}
	if date, err = BuildDate(); err != nil {
		return "", fmt.Errorf("failed to retrieve build date from git: %w", err)
	}

	var flags strings.Builder
	flags.WriteString("-X " + pkgBase + ".gitVersion=" + version)
	flags.WriteString(" ")
	flags.WriteString("-X " + pkgBase + ".gitCommit=" + hash)
	flags.WriteString(" ")
	flags.WriteString("-X " + pkgBase + ".buildDate=" + date)
	return flags.String(), nil
}

// Version returns a string that can be used as a version string.
func Version() (string, error) {
	// Use the version already set in the environment (i.e., by the CI run).
	if version, ok := os.LookupEnv("APPVERSION"); ok {
		return version, nil
	}
	// Else, derive a version from git.
	version, err := sh.Output("git", "describe", "--tags", "--always", "--dirty")
	if err != nil {
		return "", err
	}
	return version, nil
}

// hash returns the git hash for the current repo or "" if none.
func GitHash() (string, error) {
	hash, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	return hash, nil
}

// BuildDate returns the build date.
func BuildDate() (string, error) {
	date, err := sh.Output("git", "log", "--date=iso8601-strict", "-1", "--pretty=%ct")
	if err != nil {
		return "", err
	}
	return date, nil
}