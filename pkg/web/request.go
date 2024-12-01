// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/joshuar/go-hass-anything/v12/internal/logging"
)

type contextKey string

const clientContextKey contextKey = "client"

var ErrResponseMalformed = errors.New("malformed response")
var (
	defaultTimeout = 30 * time.Second
	defaultRetry   = func(r *resty.Response, _ error) bool {
		return r.StatusCode() == http.StatusTooManyRequests
	}
)

// GetRequest is a HTTP GET request.
type GetRequest interface {
	URL() string
}

// PostRequest is a HTTP POST request with the request body provided by Body().
type PostRequest interface {
	GetRequest
	RequestBody() json.RawMessage
}

type JSONResponse interface {
	json.Unmarshaler
}

type GenericResponse interface {
	Unmarshal(data []byte) error
}

func ExecuteRequest(ctx context.Context, request, response any) error {
	var (
		resp   *resty.Response
		err    error
		client *resty.Client
		found  bool
	)

	client, found = ContextGetClient(ctx)
	if !found {
		client = NewAPIClient()
	}

	webRequest := client.R().
		SetContext(ctx)

	switch req := request.(type) {
	case PostRequest:
		logging.FromContext(ctx).
			LogAttrs(ctx, logging.LevelTrace,
				"Sending request.",
				slog.String("method", "POST"),
				slog.String("body", string(req.RequestBody())),
				slog.Time("sent_at", time.Now()))

		resp, err = webRequest.SetBody(req.RequestBody()).Post(req.URL())
	case GetRequest:
		logging.FromContext(ctx).
			LogAttrs(ctx, logging.LevelTrace,
				"Sending request.",
				slog.String("method", "GET"),
				slog.Time("sent_at", time.Now()))

		resp, err = webRequest.Get(req.URL())
	}

	switch {
	case err != nil:
		return fmt.Errorf("error sending request: %w", err)
	case resp == nil:
		return fmt.Errorf("unknown error sending request")
	case resp.IsError():
		return fmt.Errorf("received error response: %v", resp.Error())
	}

	logging.FromContext(ctx).
		LogAttrs(ctx, logging.LevelTrace,
			"Received response.",
			slog.Int("statuscode", resp.StatusCode()),
			slog.String("status", resp.Status()),
			slog.String("protocol", resp.Proto()),
			slog.Duration("time", resp.Time()),
			slog.String("body", string(resp.Body())))

	if resp.IsError() {
		return fmt.Errorf("received error response: %w", err)
	}

	switch res := response.(type) {
	case JSONResponse:
		if err := res.UnmarshalJSON(resp.Body()); err != nil {
			return errors.Join(ErrResponseMalformed, err)
		}
	case GenericResponse:
		if err := res.Unmarshal(resp.Body()); err != nil {
			return errors.Join(ErrResponseMalformed, err)
		}
	}

	return nil
}

func NewAPIClient() *resty.Client {
	return resty.New().
		SetTimeout(defaultTimeout).
		AddRetryCondition(defaultRetry)
}

func ContextSetClient(ctx context.Context, client *resty.Client) context.Context {
	return context.WithValue(ctx, clientContextKey, client)
}

func ContextGetClient(ctx context.Context) (*resty.Client, bool) {
	client, ok := ctx.Value(clientContextKey).(*resty.Client)
	if !ok {
		return nil, false
	}

	return client, true
}
