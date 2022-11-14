package openfga

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	cli := NewClient("http", "127.0.0.1:8080", "", "")
	assert.NotNil(t, cli)

	stores, err := cli.ListStore(ctx)
	assert.Nil(t, err)
	if stores == nil || len(*stores) == 0 {
		_uuid := uuid.New()
		storeName := _uuid.String()
		err := cli.CreateStore(ctx, storeName)
		assert.Nil(t, err)
	} else {
		for _, store := range *stores {
			t.Logf("id: %s name:%s", store.GetId(), store.GetName())
		}

		cli.SetStoreId((*stores)[len(*stores)-1].GetId())
		//_ = cli.DeleteStore(ctx)
	}

	model := "{\"type_definitions\":[{\"type\":\"document\",\"relations\":{\"reader\":{\"this\":{}},\"writer\":{\"this\":{}},\"owner\":{\"this\":{}}}}]}"
	id, err := cli.CreateAuthorizationModel(ctx, model)
	assert.Nil(t, err)
	t.Logf("model id: %s", id)

	err = cli.CreateRelationTuple(ctx, "document:Z", "reader", "user:anne")
	assert.Nil(t, err)

	doTestData(t, cli)

	err = cli.DeleteRelationTuple(ctx, "document:Z", "reader", "user:anne")
	assert.Nil(t, err)
}

func doTestData(t *testing.T, cli *Client) {
	ctx := context.Background()

	testDatas := []struct {
		object   string
		relation string
		subject  string
		allowed  bool
	}{
		{
			object:   "document:Z",
			relation: "reader",
			subject:  "user:anne",
			allowed:  true,
		},
		{
			object:   "document:Z",
			relation: "reader",
			subject:  "user:kitty",
			allowed:  false,
		},
		{
			object:   "document:Z",
			relation: "writer",
			subject:  "user:anne",
			allowed:  false,
		},
		{
			object:   "document:Y",
			relation: "reader",
			subject:  "user:anne",
			allowed:  false,
		},
	}
	for _, test := range testDatas {
		t.Run(test.object, func(t *testing.T) {
			allowed, err := cli.GetCheck(ctx, test.object, test.relation, test.subject)
			assert.Nil(t, err)
			assert.Equal(t, test.allowed, allowed)
		})
	}
}
