// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package web

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/rs/zerolog/log"
)

type Request interface {
	Builder() *requests.Builder
	Timeout() time.Duration
}

type genericResponse struct {
	err     error
	headers map[string][]string
	body    *bytes.Buffer
}

func (r *genericResponse) Body() *bytes.Buffer {
	return r.body
}

func (r *genericResponse) Headers() map[string][]string {
	return r.headers
}

func (r *genericResponse) Error() error {
	return r.err
}

func ExecuteRequest(ctx context.Context, request Request, responseCh chan Response) {
	var resultBody bytes.Buffer
	resultHeaders := make(map[string][]string)
	defer close(responseCh)
	requestCtx, cancel := context.WithTimeout(ctx, request.Timeout())
	defer cancel()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := request.Builder().
			ToBytesBuffer(&resultBody).
			CopyHeaders(resultHeaders).
			CheckStatus(http.StatusOK).
			Fetch(requestCtx)
		if err != nil {
			log.Info().Err(err).Msg("failed request")
			cancel()
			responseCh <- &genericResponse{
				err: err,
			}
		} else {
			responseCh <- &genericResponse{
				body:    &resultBody,
				headers: resultHeaders,
			}
		}
	}()
	wg.Wait()
}
