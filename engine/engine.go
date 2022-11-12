package engine

import (
	"context"
)

type Engine interface {
	Authorizer
	Writer
}

type Authorizer interface {
	ProjectsAuthorized(context.Context) (Projects, error)

	FilterAuthorizedPairs(context.Context) (Pairs, error)

	FilterAuthorizedProjects(context.Context) (Projects, error)

	IsAuthorized(context.Context) (bool, error)
}

type Writer interface {
	SetPolicies(ctx context.Context, policyMap map[string]interface{}, roleMap map[string]interface{}) error
}
