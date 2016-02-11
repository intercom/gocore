package coreapi

import (
	"errors"
	"net/http"
)

// BasicAuth coniguration
// Calling .Protect(ContextHandlerFunc) ensures the Basic Auth request parameters match those
// configured here
type BasicAuth struct {
	User string
	Pass string
}

// ErrAuthentication returned when Authentication has failed.
var ErrAuthentication = errors.New("Authentication Error")

// Protect a given ContextHandlerFunc with a Basic Auth check
func (auth *BasicAuth) Protect(next ContextHandlerFunc) ContextHandlerFunc {
	return ContextHandlerFunc(func(ctx *ContextHandler, w http.ResponseWriter, r *http.Request) {
		key, secret, ok := r.BasicAuth()
		if !ok {
			ctx.Metrics.IncrementCount("request.unauthorized")
			JSONErrorResponse(http.StatusForbidden, ErrAuthentication).WriteTo(w)
			return
		}
		if !auth.CheckBasicAuth(key, secret) {
			ctx.Metrics.IncrementCount("request.unauthorized")
			JSONErrorResponse(http.StatusForbidden, ErrAuthentication).WriteTo(w)
			return
		}
		next(ctx, w, r)
	})
}

// CheckBasicAuth checks that the passed key and secret match the configured User and Pass.
func (auth *BasicAuth) CheckBasicAuth(key, secret string) bool {
	if key != auth.User || secret != auth.Pass {
		return false
	}
	return true
}
