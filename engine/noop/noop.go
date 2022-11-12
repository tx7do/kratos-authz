package noop

import (
	"context"
	"github.com/tx7do/kratos-authz/engine"
)

var _ engine.Engine = (*State)(nil)

type State struct{}

func (s State) ProjectsAuthorized(context.Context) (engine.Projects, error) {
	return engine.Projects{}, nil
}

func (s State) FilterAuthorizedPairs(context.Context) (engine.Pairs, error) {
	return engine.Pairs{}, nil
}

func (s State) FilterAuthorizedProjects(context.Context) (engine.Projects, error) {
	return engine.Projects{}, nil
}

func (s State) IsAuthorized(context.Context) (bool, error) {
	return true, nil
}

func (s State) SetPolicies(context.Context, map[string]interface{}, map[string]interface{}) error {
	return nil
}
