package util

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type ContextKey string

func (k ContextKey) String() string {
	return string(k)
}

func ExtractFromRequest[T any](r *http.Request, key string) (res T, err error) {
	res, ok := r.Context().Value(ContextKey(key)).(T)
	if !ok {
		err = errors.New("value not found")
		return
	}
	return
}

func InsertToRequest[T any](key string, value T) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKey(key), value)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
