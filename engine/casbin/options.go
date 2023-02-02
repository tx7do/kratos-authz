package casbin

import (
	_ "embed"

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
