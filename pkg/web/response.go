// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package web

import "bytes"

type Response interface {
	Body() *bytes.Buffer
	Headers() map[string][]string
	Error() error
}
