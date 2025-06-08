package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/tx7do/kratos-authz/engine"
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

			claims, ok := engine.AuthClaimsFromContext(ctx)
			if !ok {
				return nil, ErrMissingClaims
			}

			if claims.Action == nil || claims.Resource == nil {
				return nil, ErrInvalidClaims
			}

			var project engine.Project
			if claims.Project == nil {
				project = ""
			} else {
				project = *claims.Project
			}

			if claims.Subject != nil {
				allowed, err = authorizer.IsAuthorized(ctx, *claims.Subject, *claims.Action, *claims.Resource, project)
				if err != nil {
					return nil, err
				}
				if !allowed {
					return nil, ErrUnauthorized
				}
			} else if claims.Subjects != nil && len(*claims.Subjects) > 0 {
				for _, subject := range *claims.Subjects {
					allowed, err = authorizer.IsAuthorized(ctx, subject, *claims.Action, *claims.Resource, project)
					if err != nil {
						return nil, err
					}
					if allowed {
						break
					}
				}
				if !allowed {
					return nil, ErrUnauthorized
				}
			} else {
				return nil, ErrMissingSubject
			}

			return handler(ctx, req)
		}
	}
}
