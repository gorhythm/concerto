package http

import "context"

// ClientFinalizerFunc can be used to perform work at the end of a client HTTP request, after the response is returned.
// The principal intended use is for error logging.
// Additional response parameters are provided in the context under keys with the ContextKeyResponse prefix.
// Note: err may be nil. There maybe also no additional response parameters depending on when an error occurs.
type ClientFinalizerFunc func(ctx context.Context, err error)
