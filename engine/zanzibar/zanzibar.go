package zanzibar

import (
	"context"
	"errors"

	"github.com/tx7do/kratos-authz/engine"
	"github.com/tx7do/kratos-authz/engine/zanzibar/keto"
	"github.com/tx7do/kratos-authz/engine/zanzibar/openfga"
)

var _ engine.Engine = (*State)(nil)

type State struct {
	ketoClient    *keto.Client
	openfgaClient *openfga.Client
}

func NewEngine(_ context.Context, opts ...OptFunc) (*State, error) {
	s := &State{}

	for _, opt := range opts {
		opt(s)
	}

	if s.openfgaClient == nil && s.ketoClient == nil {
		return nil, errors.New("zanzibar client is nil")
	}

	return s, nil
}

func (s *State) Name() string {
	return string(engine.Zanzibar)
}

func (s *State) ProjectsAuthorized(_ context.Context, _ engine.Subjects, _ engine.Action, _ engine.Resource, _ engine.Projects) (engine.Projects, error) {
	return engine.Projects{}, nil
}

func (s *State) FilterAuthorizedPairs(_ context.Context, _ engine.Subjects, _ engine.Pairs) (engine.Pairs, error) {
	return engine.Pairs{}, nil
}

func (s *State) FilterAuthorizedProjects(_ context.Context, _ engine.Subjects) (engine.Projects, error) {
	return engine.Projects{}, nil
}

func (s *State) IsAuthorized(ctx context.Context, subject engine.Subject, action engine.Action, resource engine.Resource, project engine.Project) (bool, error) {
	if s.ketoClient != nil {
		allow, err := s.ketoClient.GetCheck(ctx, string(project), string(resource), string(action), string(subject))
		return allow, err
	} else if s.openfgaClient != nil {
		allow, err := s.openfgaClient.GetCheck(ctx, string(resource), string(action), string(subject))
		return allow, err
	}
	return false, nil
}

func (s *State) SetPolicies(_ context.Context, _ engine.PolicyMap, _ engine.RoleMap) error {
	return nil
}
