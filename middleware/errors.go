package middleware

import "github.com/go-kratos/kratos/v2/errors"

const (
	reason string = "FORBIDDEN"
)

var (
	ErrUnauthorized  = errors.Forbidden(reason, "unauthorized access")
	ErrMissingClaims = errors.Forbidden(reason, "missing authz claims")
	ErrInvalidClaims = errors.Forbidden(reason, "invalid authz claims")
)
