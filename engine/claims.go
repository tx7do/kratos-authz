package engine

import (
	"context"
)

type ctxKey string

var (
	authClaimsContextKey = ctxKey("authz-claims")
)

type AuthClaims struct {
	Subjects *[]string
	Pairs    *Pairs
	Projects *[]string

	Subject  *Subject
	Action   *Action
	Resource *Resource
	Project  *Project
}

// ContextWithAuthClaims injects the provided AuthClaims into the parent context.
func ContextWithAuthClaims(parent context.Context, claims *AuthClaims) context.Context {
	return context.WithValue(parent, authClaimsContextKey, claims)
}

// AuthClaimsFromContext extracts the AuthClaims from the provided ctx (if any).
func AuthClaimsFromContext(ctx context.Context) (*AuthClaims, bool) {
	claims, ok := ctx.Value(authClaimsContextKey).(*AuthClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}
