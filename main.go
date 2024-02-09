// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"log"
	"syscall"

	"github.com/joshuar/go-hass-anything/v5/cmd"
)

func init() {
	euid := syscall.Geteuid()
	uid := syscall.Getuid()
	egid := syscall.Getegid()
	gid := syscall.Getgid()
	if uid != euid || gid != egid || uid == 0 {
		log.Fatalf("go-hass-anything should not be run with additional privileges or as root.")
	}
}

func main() {
	cmd.Execute()
}
