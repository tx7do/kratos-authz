package casbin

import (
	"context"

	stdCasbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"

	"github.com/tx7do/kratos-authz/engine"
)

const (
	wildcardItem              = "*"
	authorizedProjectsMatcher = "g(r.sub, p.sub, p.dom) && (keyMatch(r.dom, p.dom) || p.dom == '*')"
)

var _ engine.Engine = (*State)(nil)

type State struct {
	model    model.Model
	policy   *Adapter
	enforcer *stdCasbin.SyncedEnforcer
	projects engine.Projects
}

func New(_ context.Context, opts ...OptFunc) (*State, error) {
	s := State{
		policy:   newAdapter(),
		projects: engine.Projects{},
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

func (s *State) ProjectsAuthorized(ctx context.Context) (engine.Projects, error) {
	claims, ok := engine.AuthClaimsFromContext(ctx)
	if !ok {
		return nil, engine.ErrMissingAuthClaims
	}

	if claims.Subjects == nil || claims.Action == nil || claims.Resource == nil || claims.Projects == nil {
		return nil, engine.ErrInvalidClaims
	}

	subjects := claims.Subjects
	projects := claims.Projects
	resource := claims.Resource
	action := claims.Action

	result := make(engine.Projects, 0, len(*projects))

	var err error
	var allowed bool
	for _, project := range *projects {
		for _, subject := range *subjects {
			if allowed, err = s.enforcer.Enforce(string(subject), string(*resource), string(*action), string(project)); err != nil {
				//fmt.Println(allowed, err)
				return nil, err
			} else if allowed {
				result = append(result, project)
			}
		}
	}

	return result, nil
}

func (s *State) FilterAuthorizedProjects(ctx context.Context) (engine.Projects, error) {
	claims, ok := engine.AuthClaimsFromContext(ctx)
	if !ok {
		return nil, engine.ErrMissingAuthClaims
	}

	if claims.Subjects == nil {
		return nil, engine.ErrInvalidClaims
	}

	subjects := claims.Subjects

	result := make(engine.Projects, 0, len(s.projects))

	var err error
	var allowed bool
	for _, project := range s.projects {
		for _, subject := range *subjects {
			if allowed, err = s.enforcer.EnforceWithMatcher(authorizedProjectsMatcher, string(subject), wildcardItem, wildcardItem, string(project)); err != nil {
				//fmt.Println(allowed, err)
				return nil, err
			} else if allowed {
				result = append(result, project)
			}
		}
	}

	return result, nil
}

func (s *State) FilterAuthorizedPairs(ctx context.Context) (engine.Pairs, error) {
	claims, ok := engine.AuthClaimsFromContext(ctx)
	if !ok {
		return nil, engine.ErrMissingAuthClaims
	}

	if claims.Subjects == nil || claims.Pairs == nil {
		return nil, engine.ErrInvalidClaims
	}

	subjects := claims.Subjects
	pairs := claims.Pairs

	result := make(engine.Pairs, 0, len(*pairs))

	var err error
	var allowed bool
	for _, p := range *pairs {
		for _, subject := range *subjects {
			if allowed, err = s.enforcer.Enforce(string(subject), string(p.Resource), string(p.Action), wildcardItem); err != nil {
				//fmt.Println(allowed, err)
				return nil, err
			} else if allowed {
				result = append(result, p)
			}
		}
	}
	return result, nil
}

func (s *State) IsAuthorized(ctx context.Context) (bool, error) {
	claims, ok := engine.AuthClaimsFromContext(ctx)
	if !ok {
		return false, engine.ErrMissingAuthClaims
	}

	if claims.Subject == nil || claims.Resource == nil || claims.Action == nil {
		return false, engine.ErrInvalidClaims
	}

	var project string
	if claims.Project == nil {
		project = wildcardItem
	} else if len(*claims.Project) > 0 {
		project = string(*claims.Project)
	}

	var err error
	var allowed bool
	if allowed, err = s.enforcer.Enforce(string(*claims.Subject), string(*claims.Resource), string(*claims.Action), project); err != nil {
		//fmt.Println(allowed, err)
		return false, err
	} else if allowed {
		return true, nil
	}
	return false, nil
}

func (s *State) SetPolicies(_ context.Context, policyMap map[string]interface{}, _ map[string]interface{}) error {
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
