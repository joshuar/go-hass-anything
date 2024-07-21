// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"errors"
	"time"

	"github.com/lthibault/jitterbug/v2"
)

var ErrInvalidPollConfig = errors.New("invalid polling configuration")

// poll is a helper function that will call the passed `updater()`
// function around each `interval` duration within the `stdev` duration window.
// Effectively, `updater()` will get called sometime near `interval`, but not
// exactly on it. This can help avoid a "thundering herd" problem of sensors all
// trying to update at the same time.
//
//nolint:exhaustruct
func poll(ctx context.Context, updater func(), interval, jitter time.Duration) error {
	if interval <= 0 || jitter <= 0 {
		return ErrInvalidPollConfig
	}

	updater()

	ticker := jitterbug.New(
		interval,
		&jitterbug.Norm{Stdev: jitter},
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updater()
			}
		}
	}()

	return nil
}
