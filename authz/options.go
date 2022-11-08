package authz

import (
	"context"
	"github.com/tx7do/kratos-authz/engine"
	"github.com/tx7do/kratos-authz/engine/casbin"
	"github.com/tx7do/kratos-authz/engine/opa"
)

type Option func(*options)

type options struct {
	securityUserCreator SecurityUserCreator
	engine              engine.Engine
	domainSupport       bool
}

func WithSecurityUserCreator(securityUserCreator SecurityUserCreator) Option {
	return func(o *options) {
		o.securityUserCreator = securityUserCreator
	}
}

func WithPolicyEngine(ctx context.Context, engineType engine.Type) Option {
	return func(o *options) {
		switch engineType {
		case engine.CasbinEngine:
			o.engine, _ = casbin.New(ctx)
		case engine.OpaEngine:
			o.engine, _ = opa.New(ctx)
		}
	}
}

func WithDomainSupport() Option {
	return func(o *options) {
		o.domainSupport = true
	}
}

func WithPolicies(ctx context.Context, policyMap map[string]interface{}, roleMap map[string]interface{}) Option {
	return func(o *options) {
		_ = o.engine.SetPolicies(ctx, policyMap, roleMap)
	}
}
