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
	var resp *resty.Response

	var err error

	var client *resty.Client

	var found bool

	client, found = ContextGetClient(ctx)
	if !found {
		client = NewAPIClient()
	}

	webRequest := client.R().
		SetContext(ctx)

	switch req := request.(type) {
	case PostRequest:
		slog.Log(ctx, LevelTrace, "api request", "method", "POST", "body", req.RequestBody(), "sent_at", time.Now())

		resp, err = webRequest.SetBody(req.RequestBody()).Post(req.URL())
	case GetRequest:
		slog.Log(ctx, LevelTrace, "api request", "method", "GET", "sent_at", time.Now())

		resp, err = webRequest.Get(req.URL())
	}

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	slog.Log(ctx, LevelTrace, "api response",
		"statuscode", resp.StatusCode(),
		"status", resp.Status(),
		"time", resp.Time(),
		"received_at", resp.ReceivedAt(),
		"body", resp.Body())

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
