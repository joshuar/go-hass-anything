// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package web

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/joshuar/go-hass-anything/v12/internal/logging"
)

const (
	defaultTimeout = 30 * time.Second
)

var (
	ErrResponseMalformed = errors.New("malformed response")
	ErrRequestFailed     = errors.New("request failed")
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

// JSONResponse is a JSON encoded response. It will be automatically
// unmarshaled.
type JSONResponse interface {
	json.Unmarshaler
}

// GenericResponse represents a non-JSON encoded response. The UnMarshal method
// will be called on the response to unmarshal into the appropriate format.
type GenericResponse interface {
	Unmarshal(data []byte) error
}

// ExecuteRequest performs a given request and stores the response in the given
// response object. If possible, the response will be automatically unmarshaled.
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
		return errors.Join(ErrRequestFailed, err)
	case resp.IsError():
		return fmt.Errorf("%w: %v", ErrRequestFailed, resp.Error())
	case resp == nil:
		return ErrRequestFailed
	}

	logging.FromContext(ctx).
		LogAttrs(ctx, logging.LevelTrace,
			"Received response.",
			slog.Int("statuscode", resp.StatusCode()),
			slog.String("status", resp.Status()),
			slog.String("protocol", resp.Proto()),
			slog.Duration("time", resp.Time()),
			slog.String("body", string(resp.Body())))

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

// Option represents an option that can be used to configure Resty client
// behavior.
type Option func(*resty.Client) *resty.Client

// SetDefaultRetry sets up the client to retry 3 times with an exponential
// backoff in-between, on any request error.
func SetDefaultRetry() Option {
	return func(client *resty.Client) *resty.Client {
		client = client.SetTimeout(defaultTimeout).
			SetRetryCount(3).
			SetRetryWaitTime(5 * time.Second).
			SetRetryMaxWaitTime(20 * time.Second).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return r.IsError() || err != nil
			})

		return client
	}
}

// SetDefaultTimeout sets a default timeout for requests.
func SetDefaultTimeout() Option {
	return func(client *resty.Client) *resty.Client {
		client = client.SetTimeout(defaultTimeout)
		return client
	}
}

// WithInsecureConnections will tell the client to ignore TLS certificate
// issues.
func WithInsecureConnections() Option {
	return func(client *resty.Client) *resty.Client {
		client = client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // #nosec G402
		return client
	}
}

// WithTLSConfig configures the client with the given tls.Config settings.
func WithTLSConfig(config *tls.Config) Option {
	return func(client *resty.Client) *resty.Client {
		client = client.SetTLSClientConfig(config)
		return client
	}
}

// NewAPIClient will create a new Resty client for use in an application. Options can
// be provided to configure the client behavior (some options can also be set
// per-request as per Resty documentation).
func NewAPIClient(options ...Option) *resty.Client {
	client := resty.New()

	for _, option := range options {
		client = option(client)
	}

	return client
}
