package zanzibar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tx7do/kratos-authz/engine"
)

func TestOpenFga(t *testing.T) {
	// Fix: WithOpenFga expects (string, string, *string, *string, *string, *string, *string)
	var token *string = nil
	var clientId *string = nil
	var clientSecret *string = nil
	var apiAudience *string = nil
	var apiTokenIssuer *string = nil
	s, err := NewEngine(t.Context(), WithOpenFga(
		"http", "127.0.0.1:8080",
		token,
		clientId, clientSecret,
		apiAudience, apiTokenIssuer,
	))
	assert.Nil(t, err)
	assert.NotNil(t, s)

	tests := []struct {
		subject  engine.Subject
		action   engine.Action
		resource engine.Resource
		project  engine.Project
		allowed  bool
	}{
		{
			resource: "document:Z",
			action:   "reader",
			subject:  "user:anne",
			allowed:  true,
		},
		{
			resource: "document:Z",
			action:   "reader",
			subject:  "user:kitty",
			allowed:  false,
		},
		{
			resource: "document:Z",
			action:   "writer",
			subject:  "user:anne",
			allowed:  false,
		},
		{
			resource: "document:Y",
			action:   "reader",
			subject:  "user:anne",
			allowed:  false,
		},
	}

	for _, test := range tests {
		t.Run(string(test.subject), func(t *testing.T) {
			allowed, err := s.IsAuthorized(t.Context(), test.subject, test.action, test.resource, test.project)
			assert.Nil(t, err)
			assert.Equal(t, test.allowed, allowed)
			//fmt.Println(r, err)
		})
	}
}

func TestKeto(t *testing.T) {
	s, err := NewEngine(t.Context(), WithKeto(
		"127.0.0.1:4466",
		"127.0.0.1:4467",
		true,
	))
	assert.Nil(t, err)
	assert.NotNil(t, s)

	tests := []struct {
		subject  engine.Subject
		action   engine.Action
		resource engine.Resource
		project  engine.Project
		allowed  bool
	}{
		{
			project:  "app",
			resource: "my-first-blog-post",
			action:   "read",
			subject:  "alice",
			allowed:  true,
		},
		{
			project:  "app1",
			resource: "my-first-blog-post",
			action:   "read",
			subject:  "alice",
			allowed:  false,
		},
		{
			project:  "app",
			resource: "obj1",
			action:   "read",
			subject:  "alice",
			allowed:  false,
		},
	}

	for _, test := range tests {
		t.Run(string(test.subject), func(t *testing.T) {
			allowed, err := s.IsAuthorized(t.Context(), test.subject, test.action, test.resource, test.project)
			assert.Nil(t, err)
			assert.Equal(t, test.allowed, allowed)
			//fmt.Println(r, err)
		})
	}
}
