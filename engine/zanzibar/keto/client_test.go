package keto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_REST(t *testing.T) {
	ctx := context.Background()

	cli := NewClient("http://127.0.0.1:4466", "http://127.0.0.1:4467", false)
	assert.NotNil(t, cli)

	err := cli.CreateRelationTuple(ctx, "app", "my-first-blog-post", "read", "alice")
	assert.Nil(t, err)

	allowed, err := cli.GetCheck(ctx, "app", "my-first-blog-post", "read", "alice")
	assert.Nil(t, err)
	assert.True(t, allowed)

	doTestData(t, cli)
}

func TestClient_GRPC(t *testing.T) {
	ctx := context.Background()

	cli := NewClient("127.0.0.1:4466", "127.0.0.1:4467", true)
	assert.NotNil(t, cli)

	err := cli.CreateRelationTuple(ctx, "app", "my-first-blog-post", "read", "alice")
	assert.Nil(t, err)

	allowed, err := cli.GetCheck(ctx, "app", "my-first-blog-post", "read", "alice")
	assert.Nil(t, err)
	assert.True(t, allowed)

	doTestData(t, cli)
}

func doTestData(t *testing.T, cli *Client) {
	ctx := context.Background()

	testDatas := []struct {
		namespace string
		object    string
		relation  string
		subjectId string
		allowed   bool
	}{
		{
			namespace: "app",
			object:    "my-first-blog-post",
			relation:  "read",
			subjectId: "alice",
			allowed:   true,
		},
		{
			namespace: "app1",
			object:    "my-first-blog-post",
			relation:  "read",
			subjectId: "alice",
			allowed:   false,
		},
		{
			namespace: "app",
			object:    "obj1",
			relation:  "read",
			subjectId: "alice",
			allowed:   false,
		},
	}
	for _, test := range testDatas {
		t.Run(test.object, func(t *testing.T) {
			allowed, err := cli.GetCheck(ctx, test.namespace, test.object, test.relation, test.subjectId)
			assert.Nil(t, err)
			assert.Equal(t, test.allowed, allowed)
		})
	}
}
