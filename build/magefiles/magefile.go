// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

const (
	appName     = "go-hass-anything"
	platformENV = "TARGETPLATFORM"
)

func init() {
	version, _ := getVersion() //nolint:errcheck
	hash, _ := getGitHash()    //nolint:errcheck

	// platform, arch, ver := parseBuildPlatform()

	// set global logger with custom options
	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
	})).
		With(
			slog.Group("git",
				slog.String("version", version),
				slog.String("hash", hash),
			),
			// slog.Group("build",
			// 	slog.String("os", platform),
			// 	slog.String("arch", arch),
			// 	slog.String("revision", ver),
			// ),
		)

	slog.SetDefault(logger)
}
