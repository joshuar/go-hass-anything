//go:build tools
// +build tools

// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	_ "github.com/davecgh/go-spew/spew"
	_ "github.com/yassinebenaid/godump"
	_ "golang.org/x/tools/cmd/stringer"
)
