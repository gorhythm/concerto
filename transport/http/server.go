package http

import (
	"context"
	"net/http"
)

// ServerFinalizerFunc can be used to perform work at the end of an HTTP request,
// after the response has been written to the client.
// The principal intended use is for request logging.
// In addition to the response code provided in the function signature,
// additional response parameters are provided in the context under keys with the ContextKeyResponse prefix.
type ServerFinalizerFunc func(ctx context.Context, code int, r *http.Request)
