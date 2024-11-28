package openfga

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/go-sdk/credentials"
)

type Client struct {
	fgaClient *client.OpenFgaClient

	apiUrl, storeId string
	credentials     credentials.Credentials
}

func NewClient(opts ...ClientOption) *Client {
	cli := &Client{
		credentials: credentials.Credentials{},
	}

	cli.init(opts...)

	return cli
}

func (c *Client) init(opts ...ClientOption) {
	for _, o := range opts {
		o(c)
	}

	if c.createApiClient() != nil {
		return
	}

	if c.ensureStore(context.Background()) != nil {
		return
	}
}

func (c *Client) ensureStore(ctx context.Context) error {
	stores, err := c.ListStore(context.Background())
	if err != nil {
		return err
	}

	if stores == nil || len(*stores) == 0 {
		_uuid := uuid.New()
		storeName := _uuid.String()
		err = c.CreateStore(ctx, storeName)
		if err != nil {
			return err
		}
	} else {
		_ = c.SetStoreId((*stores)[len(*stores)-1].GetId())
	}
	return nil
}

func (c *Client) createApiClient() error {
	cliConfig := &client.ClientConfiguration{
		ApiUrl:      c.apiUrl,
		StoreId:     c.storeId, // not needed when calling `CreateStore` or `ListStores`
		Credentials: &c.credentials,
	}

	fgaClient, err := client.NewSdkClient(cliConfig)
	if err != nil {
		log.Errorf("createApiClient error: [%s]", err.Error())
		return err
	}

	c.fgaClient = fgaClient

	return nil
}

func (c *Client) GetCheck(ctx context.Context, object, relation, subject string) (bool, error) {
	body := openfga.CheckRequest{
		TupleKey: openfga.CheckRequestTupleKey{
			User:     subject,
			Relation: relation,
			Object:   object,
		},
	}
	data, response, err := c.fgaClient.OpenFgaApi.
		Check(ctx, c.storeId).
		Body(body).
		Execute()
	if err != nil {
		log.Errorf("GetCheck error: [%s][%v]", err.Error(), response)
		return false, err
	}

	return *data.Allowed, nil
}

func (c *Client) ListStore(ctx context.Context) (*[]openfga.Store, error) {
	stores, response, err := c.fgaClient.OpenFgaApi.ListStores(ctx).Execute()
	if err != nil {
		log.Errorf("ListStore error: [%s][%v]", err.Error(), response)
		return nil, err
	}
	//log.Infof("%v", stores.Stores)
	return &stores.Stores, nil
}

func (c *Client) GetStore(ctx context.Context) string {
	store, response, err := c.fgaClient.OpenFgaApi.GetStore(ctx, c.storeId).Execute()
	if err != nil {
		log.Errorf("GetStore error [%s][%v]", err.Error(), response)
		return ""
	}
	return store.GetId()
}

func (c *Client) CreateStore(ctx context.Context, name string) error {
	store, response, err := c.fgaClient.OpenFgaApi.CreateStore(ctx).
		Body(openfga.CreateStoreRequest{
			Name: name,
		}).
		Execute()
	if err != nil {
		log.Errorf("CreateStore error: [%s][%v]", err.Error(), response)
		return err
	}

	_ = c.SetStoreId(store.GetId())

	return nil
}

func (c *Client) DeleteStore() error {
	body := openfga.ApiDeleteStoreRequest{}
	response, err := c.fgaClient.OpenFgaApi.DeleteStoreExecute(body)
	if err != nil {
		log.Errorf("DeleteStore error: [%s][%v]", err.Error(), response)
		return err
	}
	return nil
}

func (c *Client) SetStoreId(id string) error {
	return c.fgaClient.SetStoreId(id)
}

func (c *Client) CreateRelationTuple(ctx context.Context, object, relation, subject string) error {
	body := openfga.WriteRequest{
		Writes: &openfga.WriteRequestWrites{
			TupleKeys: []openfga.TupleKey{
				{
					User:     subject,
					Relation: relation,
					Object:   object,
				},
			},
		},
	}
	_, response, err := c.fgaClient.OpenFgaApi.
		Write(ctx, c.storeId).
		Body(body).
		Execute()
	if err != nil {
		log.Errorf("CreateRelationTuple error: [%s][%v]", err.Error(), response)
		return err
	}
	return nil
}

func (c *Client) DeleteRelationTuple(ctx context.Context, object, relation, subject string) error {
	body := openfga.WriteRequest{
		Deletes: &openfga.WriteRequestDeletes{
			TupleKeys: []openfga.TupleKeyWithoutCondition{
				{
					User:     subject,
					Relation: relation,
					Object:   object,
				},
			},
		},
	}
	_, response, err := c.fgaClient.OpenFgaApi.
		Write(ctx, c.storeId).
		Body(body).
		Execute()
	if err != nil {
		log.Errorf("DeleteRelationTuple error: [%s][%v]", err.Error(), response)
		return err
	}
	return nil
}

func (c *Client) ExpandRelationTuple(ctx context.Context, object, relation string) error {
	body := openfga.ExpandRequest{
		TupleKey: openfga.ExpandRequestTupleKey{
			Relation: relation,
			Object:   object,
		},
	}
	_, response, err := c.fgaClient.OpenFgaApi.
		Expand(ctx, c.storeId).
		Body(body).
		Execute()
	if err != nil {
		log.Errorf("ExpandRelationTuple error: [%s][%v]", err.Error(), response)
		return err
	}
	return nil
}

func (c *Client) CreateAuthorizationModel(ctx context.Context, writeAuthorizationModelRequestString string) (string, error) {
	var body openfga.WriteAuthorizationModelRequest
	if err := json.Unmarshal([]byte(writeAuthorizationModelRequestString), &body); err != nil {
		return "", err
	}

	data, response, err := c.fgaClient.OpenFgaApi.
		WriteAuthorizationModel(ctx, c.storeId).
		Body(body).
		Execute()
	if err != nil {
		log.Errorf("CreateAuthorizationModel error: [%s][%v]", err.Error(), response)
		return "", err
	}

	return data.GetAuthorizationModelId(), nil
}
