// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package web

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/philippta/trip"
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

func ExecuteRequest(ctx context.Context, request Request) chan Response {
	var (
		attempts = 3
		delay    = 150 * time.Millisecond
	)
	t := trip.Default(
		trip.Retry(attempts, delay, trip.RetryableStatusCodes...),
	)
	client := &http.Client{Transport: t}
	resp := &genericResponse{
		headers: make(map[string][]string),
		body:    &bytes.Buffer{},
	}
	responseCh := make(chan Response)
	requestCtx, cancel := context.WithTimeout(ctx, request.Timeout())
	go func() {
		defer close(responseCh)
		defer cancel()
		resp.err = request.Builder().
			ToBytesBuffer(resp.body).
			CopyHeaders(resp.headers).
			CheckStatus(http.StatusOK).
			Client(client).
			Fetch(requestCtx)
		responseCh <- resp
	}()
	return responseCh
}
