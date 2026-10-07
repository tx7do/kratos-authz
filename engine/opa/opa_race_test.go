package opa_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tx7do/kratos-authz/engine"
	"github.com/tx7do/kratos-authz/engine/opa"
)

// TestSetPoliciesConcurrentWithReads guards against the data race between
// SetPolicies (policy reload) and the read paths (IsAuthorized /
// ProjectsAuthorized / FilterAuthorizedPairs / FilterAuthorizedProjects).
// Must be run with -race.
func TestSetPoliciesConcurrentWithReads(t *testing.T) {
	ctx := context.Background()

	s, err := opa.NewEngine(ctx)
	require.NoError(t, err, "init state")

	policies := func(seq int) engine.PolicyMap {
		return engine.PolicyMap{
			"policies": map[string]interface{}{
				fmt.Sprintf("pol-%d", seq): map[string]interface{}{
					"members": engine.MakeSubjects("user:local:alice"),
					"statements": map[string]interface{}{
						"s-1": map[string]interface{}{
							"resources": engine.MakeResources("*"),
							"actions":   engine.MakeActions(engine.Action("iam:projects:delete")),
							"effect":    "allow",
							"projects":  engine.MakeProjects("proj-1"),
						},
					},
				},
			},
			"roles": map[string]interface{}{},
		}
	}

	require.NoError(t, s.SetPolicies(ctx, policies(0), nil), "load initial policies")

	const reloads = 50
	stop := make(chan struct{})
	var wg sync.WaitGroup

	// writer keeps reloading policies while readers evaluate queries
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= reloads; i++ {
			if err := s.SetPolicies(ctx, policies(i), nil); err != nil {
				t.Errorf("reload policies: %v", err)
				return
			}
		}
		close(stop)
	}()

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			subjects := engine.MakeSubjects("user:local:alice")
			pairs := engine.MakePairs(engine.MakePair("iam:projects", "iam:projects:delete"))
			for {
				select {
				case <-stop:
					return
				default:
				}
				_, _ = s.IsAuthorized(ctx, "user:local:alice", "iam:projects:delete", "iam:projects", "proj-1")
				_, _ = s.IsAuthorized(ctx, "user:local:alice", "iam:projects:delete", "iam:projects", "")
				_, _ = s.ProjectsAuthorized(ctx, subjects, "iam:projects:delete", "iam:projects", engine.MakeProjects("proj-1"))
				_, _ = s.FilterAuthorizedPairs(ctx, subjects, pairs)
				_, _ = s.FilterAuthorizedProjects(ctx, subjects)
			}
		}()
	}

	wg.Wait()
}
