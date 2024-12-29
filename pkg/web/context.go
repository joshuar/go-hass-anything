// Copyright 2024 Joshua Rich <joshua.rich@gmail.com>.
// SPDX-License-Identifier: MIT

package web

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type contextKey string

const clientContextKey contextKey = "client"

// ContextSetClient is used to store a Resty client for an application in the context.
func ContextSetClient(ctx context.Context, client *resty.Client) context.Context {
	return context.WithValue(ctx, clientContextKey, client)
}

// ContextGetClient is used to retrieve a Resty client from the context.
func ContextGetClient(ctx context.Context) (*resty.Client, bool) {
	client, ok := ctx.Value(clientContextKey).(*resty.Client)
	if !ok {
		return nil, false
	}

	return client, true
}
