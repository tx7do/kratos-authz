package noop

import (
	"context"
	"github.com/tx7do/kratos-authz/engine"
)

type Authorizer struct{}

var _ engine.Engine = (*Authorizer)(nil)

func (n Authorizer) ProjectsAuthorized(context.Context) (engine.Projects, error) {
	return engine.Projects{}, nil
}

func (n Authorizer) FilterAuthorizedPairs(context.Context) (engine.Pairs, error) {
	return engine.Pairs{}, nil
}

func (n Authorizer) FilterAuthorizedProjects(context.Context) (engine.Projects, error) {
	return engine.Projects{}, nil
}

func (n Authorizer) IsAuthorized(context.Context) (bool, error) {
	return true, nil
}

func (n Authorizer) SetPolicies(context.Context, map[string]interface{}, map[string]interface{}) error {
	return nil
}
