// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package helpers

import (
	"context"
	"sync"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/lthibault/jitterbug/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FormatName will take a string s and format it with appropriate spacing
// between words and capitalised the first letter of each word. For example
// someString becomes Some String. The new string is then an appropriate format
// to be used as a name in Home Assistant.
func FormatName(s string) string {
	c := cases.Title(language.AmericanEnglish)
	return c.String(strcase.ToDelimited(s, ' '))
}

// FormatID will take a string s and format it as snake_case. The new string is
// then an appropriate format to be used as a unique ID in Home Assistant.
func FormatID(s string) string {
	return strcase.ToSnake(s)
}

// PollSensors is a helper function that will call the passed `updater()`
// function around each `interval` duration within the `stdev` duration window.
// Effectively, `updater()` will get called sometime near `interval`, but not
// exactly on it. This can help avoid a "thundering herd" problem of sensors all
// trying to update at the same time.
func PollSensors(ctx context.Context, updater func(), interval, stdev time.Duration) {
	updater()
	ticker := jitterbug.New(
		interval,
		&jitterbug.Norm{Stdev: stdev},
	)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			case <-ticker.C:
				updater()
			}
		}
	}()
	wg.Wait()
}
