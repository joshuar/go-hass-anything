// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//revive:disable:unused-receiver
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"

	"github.com/joshuar/go-hass-anything/v12/internal/agent"
	"github.com/joshuar/go-hass-anything/v12/internal/logging"
)

const (
	AgentID   = "com.github.joshuar.go-hass-anything"
	AgentName = "Go Hass Anything"
)

type Context struct {
	context.Context //nolint:containedctx // workaround for https://github.com/alecthomas/kong/issues/144
}

type ResetCmd struct{}

func (r *ResetCmd) Help() string {
	return `
Reset will unregister apps from MQTT (if in use), delete the
configuration directory and remove the log file.
`
}

func (r *ResetCmd) Run(ctx Context) error {
	err := agent.ClearApps(ctx)
	if err != nil {
		return fmt.Errorf("reset agent: %w", err)
	}

	return nil
}

type ConfigureCmd struct{}

func (r *ConfigureCmd) Help() string {
	return `
Configure will present a dialog to configure the agent and any apps that have user-facing configuration.
`
}

//nolint:unparam
func (r *ConfigureCmd) Run(ctx Context) error {
	hassAgent := agent.NewAgent(ctx, AgentID, AgentName)
	hassAgent.Configure()

	return nil
}

type RunCmd struct{}

func (r *RunCmd) Help() string {
	return `
Run will run the agent and any configured apps.
`
}

func (r *RunCmd) Run(ctx Context) error {
	if err := agent.Run(ctx); err != nil {
		return fmt.Errorf("run agent failed: %w", err)
	}

	return nil
}

func init() {
	euid := syscall.Geteuid()
	uid := syscall.Getuid()
	egid := syscall.Getegid()
	gid := syscall.Getgid()

	if uid != euid || gid != egid || uid == 0 {
		slog.Error("go-hass-anything should not be run with additional privileges or as root.")
		os.Exit(-1)
	}
}

//nolint:tagalign
var CLI struct {
	Run          RunCmd               `cmd:"" help:"Run Go Hass Anything."`
	Configure    ConfigureCmd         `cmd:"" help:"Configure Go Hass Anything."`
	Reset        ResetCmd             `cmd:"" help:"Reset Go Hass Anything."`
	ProfileFlags logging.ProfileFlags `name:"profile" help:"Enable profiling."`
	LogLevel     string               `name:"log-level" enum:"info,debug,trace" default:"info" help:"Set logging level."`
	NoLogFile    bool                 `help:"Don't write to a log file."`
}

func main() {
	kong.Name(AgentName)
	env := kong.Parse(&CLI, kong.Bind(), kong.Vars{"defaultAppID": AgentID})

	ctx, cancelFunc := context.WithCancel(context.Background())

	logger := logging.New(CLI.LogLevel, CLI.NoLogFile)
	ctx = logging.ToContext(ctx, logger)

	startProfiling()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer cancelFunc()
		<-c
		stopProfiling()
		os.Exit(-1)
	}()

	if err := env.Run(Context{ctx}); err != nil {
		logger.Error("Problem running.", slog.Any("error", err))
	}

	stopProfiling()
}

func startProfiling() {
	if CLI.ProfileFlags != nil {
		if err := logging.StartProfiling(CLI.ProfileFlags); err != nil {
			slog.Warn("Problem starting profiling.",
				slog.Any("error", err))
		}
	}
}

func stopProfiling() {
	if CLI.ProfileFlags != nil {
		if err := logging.StopProfiling(CLI.ProfileFlags); err != nil {
			slog.Error("Problem stopping profiling.", slog.Any("error", err))
		}
	}
}
