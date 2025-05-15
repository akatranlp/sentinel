package csrf

import (
	"net/http"
)

// Option describes a functional option for configuring the CSRF handler.
type Option func(*csrf)

// ErrorHandler allows you to change the handler called when CSRF request
// processing encounters an invalid token or request. A typical use would be to
// provide a handler that returns a static HTML file with a HTTP 403 status. By
// default a HTTP 403 status and a plain text CSRF failure reason are served.
//
// Note that a custom error handler can also access the csrf.FailureReason(r)
// function to retrieve the CSRF validation reason from the request context.
func ErrorHandler(h http.Handler) Option {
	return func(cs *csrf) {
		cs.opts.ErrorHandler = h
	}
}

// RequestHeader allows you to change the request header the CSRF middleware
// inspects. The default is X-CSRF-Token.
func RequestHeader(header string) Option {
	return func(cs *csrf) {
		cs.opts.RequestHeader = header
	}
}

// FieldName allows you to change the name attribute of the hidden <input> field
// inspected by this package. The default is 'gorilla.csrf.Token'.
func FieldName(name string) Option {
	return func(cs *csrf) {
		cs.opts.FieldName = name
	}
}

// TrustedOrigins configures a set of origins (Referers) that are considered as trusted.
// This will allow cross-domain CSRF use-cases - e.g. where the front-end is served
// from a different domain than the API server - to correctly pass a CSRF check.
//
// You should only provide origins you own or have full control over.
func TrustedOrigins(origins []string) Option {
	return func(cs *csrf) {
		cs.opts.TrustedOrigins = origins
	}
}

// parseOptions parses the supplied options functions and returns a configured
// csrf handler.
func parseOptions(h http.Handler, opts ...Option) *csrf {
	// Set the handler to call after processing.
	cs := &csrf{
		h: h,
	}

	// Range over each options function and apply it
	// to our csrf type to configure it. Options functions are
	// applied in order, with any conflicting options overriding
	// earlier calls.
	for _, option := range opts {
		option(cs)
	}

	return cs
}
