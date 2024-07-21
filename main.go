// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//revive:disable:unused-receiver
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"syscall"

	"github.com/alecthomas/kong"

	"github.com/joshuar/go-hass-anything/v10/internal/agent"
	"github.com/joshuar/go-hass-anything/v10/internal/logging"
)

const (
	AgentID   = "com.github.joshuar.go-hass-anything"
	AgentName = "Go Hass Anything"
)

type Context struct {
	Profile   profileFlags
	LogLevel  string
	NoLogFile bool
}

type profileFlags logging.ProfileFlags

func (d profileFlags) AfterApply() error {
	err := logging.StartProfiling(logging.ProfileFlags(d))
	if err != nil {
		return fmt.Errorf("could not start profiling: %w", err)
	}

	return nil
}

type ResetCmd struct{}

func (r *ResetCmd) Help() string {
	return `
Reset will unregister apps from MQTT (if in use), delete the
configuration directory and remove the log file.
`
}

func (r *ResetCmd) Run(ctx *Context) error {
	agentCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	logger := logging.New(ctx.LogLevel, ctx.NoLogFile)
	agentCtx = logging.ToContext(agentCtx, logger)

	err := agent.ClearApps(agentCtx)
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

func (r *ConfigureCmd) Run(ctx *Context) error {
	agentCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	logger := logging.New(ctx.LogLevel, ctx.NoLogFile)
	agentCtx = logging.ToContext(agentCtx, logger)

	hassAgent := agent.NewAgent(agentCtx, AgentID, AgentName)
	hassAgent.Configure()

	return nil
}

type RunCmd struct{}

func (r *RunCmd) Help() string {
	return `
Run will run the agent and any configured apps.
`
}

func (r *RunCmd) Run(ctx *Context) error {
	agentCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	logger := logging.New(ctx.LogLevel, ctx.NoLogFile)
	agentCtx = logging.ToContext(agentCtx, logger)

	if err := agent.Run(agentCtx); err != nil {
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
		log.Fatalf("go-hass-anything should not be run with additional privileges or as root.")
	}
}

//nolint:tagalign
var CLI struct {
	Run       RunCmd       `cmd:"" help:"Run Go Hass Anything."`
	Configure ConfigureCmd `cmd:"" help:"Configure Go Hass Anything."`
	Reset     ResetCmd     `cmd:"" help:"Reset Go Hass Anything."`
	Profile   profileFlags `help:"Enable profiling."`
	LogLevel  string       `name:"log-level" enum:"info,debug,trace" default:"info" help:"Set logging level."`
	NoLogFile bool         `help:"Don't write to a log file."`
}

func main() {
	kong.Name(AgentName)
	ctx := kong.Parse(&CLI, kong.Bind(), kong.Vars{"defaultAppID": AgentID})

	err := ctx.Run(&Context{Profile: CLI.Profile, LogLevel: CLI.LogLevel, NoLogFile: CLI.NoLogFile})
	if CLI.Profile != nil {
		err = errors.Join(logging.StopProfiling(logging.ProfileFlags(CLI.Profile)), err)
	}

	ctx.FatalIfErrorf(err)
}
