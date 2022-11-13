package zanzibar

import (
	"context"
	"github.com/tx7do/kratos-authz/engine"
)

var _ engine.Engine = (*State)(nil)

type State struct{}

func (s *State) ProjectsAuthorized(_ context.Context, _ engine.Subjects, _ engine.Action, _ engine.Resource, _ engine.Projects) (engine.Projects, error) {
	return engine.Projects{}, nil
}

func (s *State) FilterAuthorizedPairs(_ context.Context, _ engine.Subjects, _ engine.Pairs) (engine.Pairs, error) {
	return engine.Pairs{}, nil
}

func (s *State) FilterAuthorizedProjects(_ context.Context, _ engine.Subjects) (engine.Projects, error) {
	return engine.Projects{}, nil
}

func (s *State) IsAuthorized(_ context.Context, _ engine.Subject, _ engine.Action, _ engine.Resource, _ engine.Project) (bool, error) {
	return true, nil
}

func (s *State) SetPolicies(_ context.Context, _ engine.PolicyMap, _ engine.RoleMap) error {
	return nil
}
