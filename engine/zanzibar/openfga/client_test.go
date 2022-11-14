package openfga

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
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
			log.Infof("id: %s name:%s", store.GetId(), store.GetName())
		}

		cli.SetStoreId((*stores)[len(*stores)-1].GetId())
		//_ = cli.DeleteStore(ctx)
	}
}
