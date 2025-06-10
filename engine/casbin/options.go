package casbin

import (
	"github.com/casbin/casbin/v2/model"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/tx7do/kratos-authz/engine"
	"github.com/tx7do/kratos-authz/engine/casbin/assets"
)

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

func WithDefaultModel(name string) OptFunc {
	return func(s *State) {
		var str string
		switch name {
		case "rbac":
			str = assets.DefaultRbacModel

		case "rbac_with_domains":
			str = assets.DefaultRbacWithDomainModel

		case "abac":
			str = assets.DefaultAbacModel

		case "acl":
			str = assets.DefaultAclModel

		case "restfull":
			str = assets.DefaultRestfullModel

		case "restfull_with_role":
			str = assets.DefaultRestfullWithRoleModel
		}

		s.model, _ = model.NewModelFromString(str)
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

func WithLogger(logger log.Logger) OptFunc {
	return func(s *State) {
		s.log = log.NewHelper(log.With(logger, "module", "casbin.authz.engine"))
	}
}
