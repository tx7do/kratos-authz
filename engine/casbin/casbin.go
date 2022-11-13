package casbin

import (
	"context"

	stdCasbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"

	"github.com/tx7do/kratos-authz/engine"
)

var _ engine.Engine = (*State)(nil)

type State struct {
	model    model.Model
	policy   *Adapter
	enforcer *stdCasbin.SyncedEnforcer
	projects engine.Projects

	wildcardItem              string
	authorizedProjectsMatcher string
}

func New(_ context.Context, opts ...OptFunc) (*State, error) {
	s := State{
		policy:                    newAdapter(),
		projects:                  engine.Projects{},
		wildcardItem:              "*",
		authorizedProjectsMatcher: "g(r.sub, p.sub, p.dom) && (keyMatch(r.dom, p.dom) || p.dom == '*')",
	}

	for _, opt := range opts {
		opt(&s)
	}

	var err error

	if s.model == nil {
		s.model, err = model.NewModelFromString(DefaultRestfullWithRoleModel)
		if err != nil {
			return nil, err
		}
	}

	s.enforcer, err = stdCasbin.NewSyncedEnforcer(s.model, s.policy)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *State) ProjectsAuthorized(_ context.Context, subjects engine.Subjects, action engine.Action, resource engine.Resource, projects engine.Projects) (engine.Projects, error) {
	result := make(engine.Projects, 0, len(projects))

	var err error
	var allowed bool
	for _, project := range projects {
		for _, subject := range subjects {
			if allowed, err = s.enforcer.Enforce(string(subject), string(resource), string(action), string(project)); err != nil {
				//fmt.Println(allowed, err)
				return nil, err
			} else if allowed {
				result = append(result, project)
			}
		}
	}

	return result, nil
}

func (s *State) FilterAuthorizedPairs(_ context.Context, subjects engine.Subjects, pairs engine.Pairs) (engine.Pairs, error) {
	result := make(engine.Pairs, 0, len(pairs))

	project := engine.Project(s.wildcardItem)

	var err error
	var allowed bool
	for _, p := range pairs {
		for _, subject := range subjects {
			if allowed, err = s.enforcer.Enforce(string(subject), string(p.Resource), string(p.Action), string(project)); err != nil {
				//fmt.Println(allowed, err)
				return nil, err
			} else if allowed {
				result = append(result, p)
			}
		}
	}
	return result, nil
}

func (s *State) FilterAuthorizedProjects(_ context.Context, subjects engine.Subjects) (engine.Projects, error) {
	result := make(engine.Projects, 0, len(s.projects))

	resource := engine.Resource(s.wildcardItem)
	action := engine.Action(s.wildcardItem)

	var err error
	var allowed bool
	for _, project := range s.projects {
		for _, subject := range subjects {
			if allowed, err = s.enforcer.EnforceWithMatcher(s.authorizedProjectsMatcher, string(subject), string(resource), string(action), string(project)); err != nil {
				//fmt.Println(allowed, err)
				return nil, err
			} else if allowed {
				result = append(result, project)
			}
		}
	}

	return result, nil
}

func (s *State) IsAuthorized(_ context.Context, subject engine.Subject, action engine.Action, resource engine.Resource, project engine.Project) (bool, error) {
	if len(project) == 0 {
		project = engine.Project(s.wildcardItem)
	}

	var err error
	var allowed bool
	if allowed, err = s.enforcer.Enforce(string(subject), string(resource), string(action), string(project)); err != nil {
		//fmt.Println(allowed, err)
		return false, err
	} else if allowed {
		return true, nil
	}
	return false, nil
}

func (s *State) SetPolicies(_ context.Context, policyMap engine.PolicyMap, _ engine.RoleMap) error {
	s.policy.SetPolicies(policyMap)
	err := s.enforcer.LoadPolicy()
	//fmt.Println(err, s.enforcer.GetAllSubjects(), s.enforcer.GetAllRoles())

	projects, ok := policyMap["projects"]
	if ok {
		switch t := projects.(type) {
		case engine.Projects:
			s.projects = t
		}
	}

	return err
}
