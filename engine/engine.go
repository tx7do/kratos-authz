package engine

import (
	"context"
)

type Type int

const (
	CasbinEngine Type = 1
	OpaEngine    Type = 2
)

type Engine interface {
	Authorizer
	Writer
}

type Authorizer interface {
	ProjectsAuthorized(ctx context.Context, subjects Subjects, action Action, resource Resource, projects Projects) (Projects, error)

	FilterAuthorizedPairs(ctx context.Context, subjects Subjects, pairs Pairs) (Pairs, error)

	FilterAuthorizedProjects(ctx context.Context, subjects Subjects) (Projects, error)
}

type Writer interface {
	SetPolicies(ctx context.Context, policyMap map[string]interface{}, roleMap map[string]interface{}) error
}

type Subjects []string

func MakeSubjects(subs ...string) Subjects {
	return subs
}

type Project string
type Projects []Project

func MakeProjects(projects ...Project) Projects {
	return projects
}

type Action string
type Actions []Action

func MakeActions(actions ...Action) Actions {
	return actions
}

type Resource string
type Resources []Resource

func MakeResources(resources ...Resource) Resources {
	return resources
}

type Pair struct {
	Resource Resource `json:"resource"`
	Action   Action   `json:"action"`
}
type Pairs []Pair

func MakePair(res, act string) Pair {
	return Pair{Resource(res), Action(act)}
}
func MakePairs(pairs ...Pair) Pairs {
	return pairs
}
