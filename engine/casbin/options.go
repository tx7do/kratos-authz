package casbin

import (
	_ "embed"
	"github.com/tx7do/kratos-authz/engine"

	"github.com/casbin/casbin/v2/model"
)

//go:embed model/rbac.conf
var DefaultRbacModel string

//go:embed model/rbac_with_domains.conf
var DefaultRbacWithDomainModel string

//go:embed model/abac.conf
var DefaultAbacModel string

//go:embed model/acl.conf
var DefaultAclModel string

//go:embed model/restfull.conf
var DefaultRestfullModel string

//go:embed model/restfull_with_role.conf
var DefaultRestfullWithRoleModel string

type OptFunc func(*State)

func WithModel(model model.Model) OptFunc {
	return func(s *State) {
		s.model = model
	}
}

func WithStringModel(str string) OptFunc {
	return func(s *State) {
		s.model, _ = model.NewModelFromString(str)
	}
}

func WithFileModel(path string) OptFunc {
	return func(s *State) {
		s.model, _ = model.NewModelFromFile(path)
	}
}

func WithPolicyAdapter(policy *Adapter) OptFunc {
	return func(s *State) {
		s.policy = policy
	}
}

func WithPolices(policies map[string]interface{}) OptFunc {
	return func(s *State) {
		if s.policy == nil {
			s.policy = newAdapter()
		}
		s.policy.SetPolicies(policies)
	}
}

func WithProjects(projects engine.Projects) OptFunc {
	return func(s *State) {
		s.projects = projects
	}
}

func WithWildcardItem(item string) OptFunc {
	return func(s *State) {
		s.wildcardItem = item
	}
}

func WithAuthorizedProjectsMatcher(matcher string) OptFunc {
	return func(s *State) {
		s.authorizedProjectsMatcher = matcher
	}
}
