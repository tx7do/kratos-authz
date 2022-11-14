package openfga

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/credentials"
)

type Client struct {
	apiClient *openfga.APIClient
}

func NewClient(scheme, host, storeId, token string) *Client {
	cli := &Client{}

	if cli.createApiClient(scheme, host, storeId, token) != nil {
		return nil
	}

	return cli
}

func (c *Client) createApiClient(scheme, host, storeId, token string) error {
	rawConfig := openfga.Configuration{
		ApiScheme: scheme,  // optional, defaults to "https"
		ApiHost:   host,    // required, define without the scheme (e.g. api.fga.example instead of https://api.fga.example)
		StoreId:   storeId, // not needed when calling `CreateStore` or `ListStores`
	}

	if token != "" {
		rawConfig.Credentials = &credentials.Credentials{
			Method: credentials.CredentialsMethodApiToken,
			Config: &credentials.Config{
				ApiToken: token, // will be passed as the "Authorization: Bearer ${ApiToken}" request header
			},
		}
	}

	configuration, err := openfga.NewConfiguration(rawConfig)
	if err != nil {
		return err
	}

	c.apiClient = openfga.NewAPIClient(configuration)

	return nil
}

func (c *Client) GetCheck(object, relation, subject string) (bool, error) {
	body := openfga.CheckRequest{
		TupleKey: &openfga.TupleKey{
			User:     openfga.PtrString(subject),
			Relation: openfga.PtrString(relation),
			Object:   openfga.PtrString(object),
		},
	}
	data, response, err := c.apiClient.OpenFgaApi.Check(context.Background()).Body(body).Execute()
	if err != nil {
		log.Errorf("GetCheck error: %v\n", response)
		return false, err
	}

	return *data.Allowed, nil
}

func (c *Client) ListStore(ctx context.Context) (*[]openfga.Store, error) {
	stores, response, err := c.apiClient.OpenFgaApi.ListStores(ctx).Execute()
	if err != nil {
		log.Errorf("ListStore error: [%s][%v]\n", err.Error(), response)
		return nil, err
	}
	//log.Infof("%v", stores.Stores)
	return stores.Stores, nil
}

func (c *Client) GetStore(ctx context.Context) string {
	store, response, err := c.apiClient.OpenFgaApi.GetStore(ctx).Execute()
	if err != nil {
		log.Errorf("GetStore error [%s][%v]\n", err.Error(), response)
		return ""
	}
	return store.GetId()
}

func (c *Client) CreateStore(ctx context.Context, name string) error {
	store, response, err := c.apiClient.OpenFgaApi.CreateStore(ctx).
		Body(openfga.CreateStoreRequest{
			Name: openfga.PtrString(name),
		}).
		Execute()
	if err != nil {
		log.Errorf("CreateStore error: [%s][%v]\n", err.Error(), response)
		return err
	}

	c.SetStoreId(store.GetId())

	return nil
}

func (c *Client) DeleteStore(ctx context.Context) error {
	body := openfga.ApiDeleteStoreRequest{}
	response, err := c.apiClient.OpenFgaApi.DeleteStoreExecute(body)
	if err != nil {
		log.Errorf("DeleteStore error: [%s][%v]\n", err.Error(), response)
		return err
	}
	return nil
}

func (c *Client) SetStoreId(id string) {
	c.apiClient.SetStoreId(id)
}
