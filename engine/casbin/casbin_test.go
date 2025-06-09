package casbin

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tx7do/kratos-authz/engine"
)

var (
	allProjects = engine.Projects{
		"(unassigned)",
		"project1",
		"project2",
		"project3",
		"project4",
		"project5",
		"project6",
	}
)

func TestFilterAuthorizedPairs(t *testing.T) {
	ctx := context.Background()
	s, err := NewEngine(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	policies := map[string]interface{}{
		"policies": []PolicyRule{
			{PType: "p", V0: "bobo", V1: "/api/*", V2: "(GET)|(POST)", V3: "*"},
			{PType: "p", V0: "bobo01", V1: "/api/users", V2: "GET", V3: "*"},
			{PType: "p", V0: "admin_role", V1: "/api/*", V2: "(GET)|(POST)", V3: "*"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "*"},
		},
		"projects": engine.MakeProjects(),
	}

	err = s.SetPolicies(ctx, policies, nil)
	assert.Nil(t, err)

	tests := []struct {
		authorityId string
		path        string
		action      string
		equal       engine.Pairs
	}{
		{
			authorityId: "admin",
			path:        "/api/login",
			action:      "POST",
			equal:       engine.Pairs{engine.MakePair("/api/login", "POST")},
		},
		{
			authorityId: "admin",
			path:        "/api/logout",
			action:      "POST",
			equal:       engine.Pairs{engine.MakePair("/api/logout", "POST")},
		},
		{
			authorityId: "bobo",
			path:        "/api/login",
			action:      "POST",
			equal:       engine.Pairs{engine.MakePair("/api/login", "POST")},
		},
		{
			authorityId: "bobo01",
			path:        "/api/login",
			action:      "POST",
			equal:       engine.Pairs{},
		},
		{
			authorityId: "bobo01",
			path:        "/api/users",
			action:      "GET",
			equal:       engine.Pairs{engine.MakePair("/api/users", "GET")},
		},
		{
			authorityId: "bobo01",
			path:        "/api/users",
			action:      "POST",
			equal:       engine.Pairs{},
		},
	}

	for _, test := range tests {
		t.Run(test.authorityId, func(t *testing.T) {
			subjects := engine.MakeSubjects(engine.Subject(test.authorityId))
			pairs := engine.MakePairs(engine.MakePair(test.path, test.action))
			r, err := s.FilterAuthorizedPairs(ctx, subjects, pairs)
			assert.Nil(t, err)
			assert.EqualValues(t, test.equal, r)
			//fmt.Println(r, err)
		})
	}
}

func TestFilterAuthorizedProjects(t *testing.T) {
	ctx := context.Background()
	s, err := NewEngine(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	policies := map[string]interface{}{
		"policies": []PolicyRule{
			{PType: "p", V0: "bobo", V1: "/api/*", V2: "(GET)|(POST)", V3: "project1"},
			{PType: "p", V0: "bobo", V1: "/api/*", V2: "(GET)|(POST)", V3: "project2"},
			{PType: "p", V0: "bobo01", V1: "/api/users", V2: "GET", V3: "*"},
			{PType: "p", V0: "admin_role", V1: "/api/*", V2: "(GET)|(POST)", V3: "*"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "*"},
		},
		"projects": allProjects,
	}

	err = s.SetPolicies(ctx, policies, nil)
	assert.Nil(t, err)

	subjects := engine.Subjects{"bobo"}

	r, err := s.FilterAuthorizedProjects(ctx, subjects)
	assert.Nil(t, err)
	fmt.Println(r)

	tests := []struct {
		subjects engine.Subjects
		equal    engine.Projects
	}{
		{
			subjects: engine.MakeSubjects("bobo"),
			equal:    engine.Projects{"project1", "project2"},
		},
		{
			subjects: engine.MakeSubjects("bobo01"),
			equal:    allProjects,
		},
		{
			subjects: engine.MakeSubjects("admin"),
			equal:    allProjects,
		},
		{
			subjects: engine.MakeSubjects("admin_role"),
			equal:    allProjects,
		},
	}

	for _, test := range tests {
		t.Run(string(test.subjects[0]), func(t *testing.T) {
			r, err := s.FilterAuthorizedProjects(ctx, test.subjects)
			assert.Nil(t, err)
			assert.EqualValues(t, test.equal, r)
			//fmt.Println(r, err)
		})
	}
}

func TestProjectsAuthorized(t *testing.T) {
	ctx := context.Background()
	s, err := NewEngine(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	policies := map[string]interface{}{
		"policies": []PolicyRule{
			{PType: "p", V0: "bobo", V1: "/api/*", V2: "(GET)|(POST)", V3: "project1"},
			{PType: "p", V0: "bobo", V1: "/api/*", V2: "(GET)|(POST)", V3: "project2"},
			{PType: "p", V0: "bobo01", V1: "/api/users", V2: "GET", V3: "*"},
			{PType: "p", V0: "admin_role", V1: "/api/*", V2: "(GET)|(POST)", V3: "*"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "project1"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "project2"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "project3"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "project4"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "project5"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "project6"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "(unassigned)"},
		},
		"projects": allProjects,
	}

	err = s.SetPolicies(ctx, policies, nil)
	assert.Nil(t, err)

	subjects := engine.Subjects{"bobo"}
	action := engine.Action("GET")
	resource := engine.Resource("/api/users")
	projects := engine.Projects{"project1"}
	r, err := s.ProjectsAuthorized(ctx, subjects, action, resource, projects)
	assert.Nil(t, err)
	fmt.Println(r)

	tests := []struct {
		subjects engine.Subjects
		action   engine.Action
		resource engine.Resource
		projects engine.Projects
		equal    engine.Projects
	}{
		{
			subjects: engine.MakeSubjects("bobo"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects("project1"),
			equal:    engine.Projects{"project1"},
		},
		{
			subjects: engine.MakeSubjects("bobo"),
			action:   engine.Action("POST"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects("project1"),
			equal:    engine.Projects{"project1"},
		},
		{
			subjects: engine.MakeSubjects("bobo"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/projects"),
			projects: engine.MakeProjects("project1"),
			equal:    engine.Projects{"project1"},
		},
		{
			subjects: engine.MakeSubjects("bobo"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects("project2"),
			equal:    engine.Projects{"project2"},
		},
		{
			subjects: engine.MakeSubjects("bobo"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects("project3"),
			equal:    engine.Projects{},
		},
		{
			subjects: engine.MakeSubjects("bobo"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api1/users"),
			projects: engine.MakeProjects("project1"),
			equal:    engine.Projects{},
		},
		{
			subjects: engine.MakeSubjects("bobo"),
			action:   engine.Action("DELETE"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects("project1"),
			equal:    engine.Projects{},
		},
		{
			subjects: engine.MakeSubjects("bobo999"),
			action:   engine.Action("DELETE"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects("project1"),
			equal:    engine.Projects{},
		},
		{
			subjects: engine.MakeSubjects("bobo01"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects("project1"),
			equal:    engine.Projects{"project1"},
		},
		{
			subjects: engine.MakeSubjects("bobo01"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects(allProjects...),
			equal:    allProjects,
		},
		{
			subjects: engine.MakeSubjects("admin"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects("project1"),
			equal:    engine.Projects{"project1"},
		},
		{
			subjects: engine.MakeSubjects("admin"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects(allProjects...),
			equal:    allProjects,
		},
		{
			subjects: engine.MakeSubjects("admin_role"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects("project1"),
			equal:    engine.Projects{"project1"},
		},
		{
			subjects: engine.MakeSubjects("admin_role"),
			action:   engine.Action("GET"),
			resource: engine.Resource("/api/users"),
			projects: engine.MakeProjects(allProjects...),
			equal:    allProjects,
		},
	}

	for _, test := range tests {
		t.Run(string(test.subjects[0]), func(t *testing.T) {
			r, err := s.ProjectsAuthorized(ctx, test.subjects, test.action, test.resource, test.projects)
			assert.Nil(t, err)
			assert.EqualValues(t, test.equal, r)
			//fmt.Println(r, err)
		})
	}
}
