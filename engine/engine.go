package engine

import (
	"context"
)

type Engine interface {
	Authorizer
	Writer
}

type Authorizer interface {
	ProjectsAuthorized(context.Context, Subjects, Action, Resource, Projects) (Projects, error)

	FilterAuthorizedPairs(context.Context, Subjects, Pairs) (Pairs, error)

	FilterAuthorizedProjects(context.Context, Subjects) (Projects, error)

	IsAuthorized(context.Context, Subject, Action, Resource, Project) (bool, error)
}

type Writer interface {
	SetPolicies(context.Context, PolicyMap, RoleMap) error
}
