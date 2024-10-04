//go:build tools
// +build tools

// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	_ "github.com/davecgh/go-spew/spew"
	_ "github.com/google/go-licenses"
	_ "github.com/magefile/mage"
	_ "github.com/sigstore/cosign/v2/cmd/cosign"
	_ "github.com/yassinebenaid/godump"
	_ "go.uber.org/nilaway/cmd/nilaway"
	_ "golang.org/x/tools/cmd/stringer"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
