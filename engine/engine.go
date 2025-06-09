package engine

import (
	"context"
)

type Engine interface {
	Authorizer
	Writer
}

type Authorizer interface {
	Name() string

	ProjectsAuthorized(ctx context.Context, subjects Subjects, action Action, resource Resource, projects Projects) (Projects, error)

	FilterAuthorizedPairs(ctx context.Context, subjects Subjects, pairs Pairs) (Pairs, error)

	FilterAuthorizedProjects(ctx context.Context, subjects Subjects) (Projects, error)

	IsAuthorized(ctx context.Context, subjects Subject, action Action, resource Resource, project Project) (bool, error)
}

type Writer interface {
	SetPolicies(ctx context.Context, policies PolicyMap, roles RoleMap) error
}
