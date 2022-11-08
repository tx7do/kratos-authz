package authz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/tx7do/kratos-authz/engine"
	"github.com/tx7do/kratos-authz/engine/casbin"
)

type contextKey string

const (
	SecurityUserContextKey contextKey = "SecurityUser"
	reason                 string     = "FORBIDDEN"
)

var (
	ErrSecurityUserCreatorMissing = errors.Forbidden(reason, "SecurityUserCreator is required")
	ErrEngineMissing              = errors.Forbidden(reason, "Engine is missing")
	ErrSecurityParseFailed        = errors.Forbidden(reason, "Security Info fault")
	ErrUnauthorized               = errors.Forbidden(reason, "Unauthorized Access")
)

func isAllowedTuple4(ctx context.Context, e engine.Engine, sub, obj, act, dom string) (bool, error) {
	result, err := e.IsProjectAuthorized(ctx,
		engine.Subject(sub),
		engine.Action(act),
		engine.Resource(obj),
		engine.Project(dom),
	)
	//fmt.Println(result, sub, obj, act, dom)
	return result, err
}

func isAllowedTuple3(ctx context.Context, e engine.Engine, sub, obj, act string) (bool, error) {
	result, err := e.IsAuthorized(ctx,
		engine.Subject(sub),
		engine.Action(act),
		engine.Resource(obj),
	)
	return result, err
}

func isAllowed(ctx context.Context, e engine.Engine, request SecurityUser, withDomain bool) (bool, error) {
	if withDomain {
		return isAllowedTuple4(ctx, e, request.GetSubject(), request.GetObject(), request.GetAction(), request.GetDomain())
	} else {
		return isAllowedTuple3(ctx, e, request.GetSubject(), request.GetObject(), request.GetAction())
	}
}

func Server(opts ...Option) middleware.Middleware {
	o := &options{
		domainSupport: false,
	}

	for _, opt := range opts {
		opt(o)
	}

	if o.engine == nil {
		o.engine, _ = casbin.New(nil)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				allowed bool
				err     error
			)

			if o.engine == nil {
				return nil, ErrEngineMissing
			}

			if o.securityUserCreator == nil {
				return nil, ErrSecurityUserCreatorMissing
			}

			securityUser := o.securityUserCreator()
			if err := securityUser.ParseFromContext(ctx); err != nil {
				return nil, ErrSecurityParseFailed
			}

			if allowed, err = isAllowed(ctx, o.engine, securityUser, o.domainSupport); err != nil {
				return nil, err
			} else if !allowed {
				return nil, ErrUnauthorized
			}

			ctx = context.WithValue(ctx, SecurityUserContextKey, securityUser)
			return handler(ctx, req)
		}
	}
}

// SecurityUserFromContext extract SecurityUser from context
func SecurityUserFromContext(ctx context.Context) (SecurityUser, bool) {
	user, ok := ctx.Value(SecurityUserContextKey).(SecurityUser)
	return user, ok
}
