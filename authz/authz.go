package authz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/tx7do/kratos-authz/engine"
)

type contextKey string

const (
	reason string = "FORBIDDEN"
)

var (
	ErrUnauthorized = errors.Forbidden(reason, "unauthorized access")
)

func Server(authorizer engine.Authorizer, opts ...Option) middleware.Middleware {
	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	if authorizer == nil {
		return nil
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				allowed bool
				err     error
			)

			allowed, err = authorizer.IsAuthorized(ctx)
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, ErrUnauthorized
			}

			return handler(ctx, req)
		}
	}
}
