// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package preferences

import (
	"errors"
	"fmt"
	"os"
)

var ErrUnknownPref = errors.New("unknown preference")

//nolint:unused // these are reserved for future use
var (
	gitVersion, gitCommit, gitTreeState, buildDate string
	AppVersion                                     = gitVersion
)

// Preference represents a single preference in a preferences file.
type Preference struct {
	// Value is the actual preference value.
	Value any `toml:"value"`
	// Description is a string that describes the preference, and may be used
	// for display purposes.
	Description string `toml:"description,omitempty"`
	// Secret is a flag that indicates whether this preference represents a
	// secret. The value has no effect on the preference encoding in the TOML,
	// only on how to display the preference to the user (masked or plaintext).
	Secret bool `toml:"-"`
}

func checkPath(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create new directory: %w", err)
		}
	}

	return nil
}
