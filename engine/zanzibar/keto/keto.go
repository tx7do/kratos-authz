package keto

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/go-kratos/kratos/v2/log"

	client "github.com/ory/keto-client-go"
	acl "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type Client struct {
	checkServiceClient  acl.CheckServiceClient
	readServiceClient   acl.ReadServiceClient
	writeServiceClient  acl.WriteServiceClient
	expandServiceClient acl.ExpandServiceClient

	readClient  *client.APIClient
	writeClient *client.APIClient

	useGRPC bool
}

func NewClient(readUrl, writeUrl string, useGRPC bool) *Client {
	cli := &Client{
		useGRPC: useGRPC,
	}

	if useGRPC {
		cli.createGrpcWriteClient(writeUrl)
		cli.createGrpcReadClient(readUrl)
	} else {
		cli.createRestWriteClient(writeUrl)
		cli.createRestReadClient(readUrl)
	}

	return cli
}

func (c *Client) GetCheck(ctx context.Context, namespace, object, relation, subject string) (bool, error) {
	if c.useGRPC {
		return c.grpcGetCheck(ctx, namespace, object, relation, subject)
	} else {
		return c.restCheckPermission(ctx, namespace, object, relation, subject)
	}
}

func (c *Client) CreateRelationTuple(ctx context.Context, namespace, object, relation, subject string) error {
	if c.useGRPC {
		return c.grpcCreateRelationTuple(ctx, namespace, object, relation, subject)
	} else {
		return c.restCreateRelationTuple(ctx, namespace, object, relation, subject)
	}
}

func (c *Client) createGrpcReadClient(uri string) {
	conn, err := grpc.NewClient(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Encountered error: " + err.Error())
	}

	c.checkServiceClient = acl.NewCheckServiceClient(conn)
	c.readServiceClient = acl.NewReadServiceClient(conn)
	c.expandServiceClient = acl.NewExpandServiceClient(conn)
}

func (c *Client) createGrpcWriteClient(uri string) {
	conn, err := grpc.NewClient(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Encountered error: " + err.Error())
	}

	c.writeServiceClient = acl.NewWriteServiceClient(conn)
}

func (c *Client) createRestReadClient(uri string) {
	configuration := client.NewConfiguration()
	configuration.Servers = []client.ServerConfiguration{
		{
			URL: uri,
		},
	}
	c.readClient = client.NewAPIClient(configuration)
}

func (c *Client) createRestWriteClient(uri string) {
	configuration := client.NewConfiguration()
	configuration.Servers = []client.ServerConfiguration{
		{
			URL: uri,
		},
	}
	c.writeClient = client.NewAPIClient(configuration)
}

func (c *Client) restCreateRelationTuple(ctx context.Context, namespace, object, relation, subject string) error {
	relationQuery := *client.NewCreateRelationshipBody()
	relationQuery.SetNamespace(namespace)
	relationQuery.SetObject(object)
	relationQuery.SetRelation(relation)
	relationQuery.SetSubjectId(subject)

	_, r, err := c.writeClient.RelationshipApi.CreateRelationship(ctx).
		CreateRelationshipBody(relationQuery).
		Execute()
	if err != nil {
		log.Errorf("restCreateRelationTuple error: [%s][%v]", err.Error(), r)
		return err
	}

	return nil
}

func (c *Client) restCheckPermission(ctx context.Context, namespace, object, relation, subject string) (bool, error) {
	check, r, err := c.readClient.PermissionApi.CheckPermission(ctx).
		Namespace(namespace).
		Object(object).
		Relation(relation).
		SubjectId(subject).
		Execute()
	if err != nil {
		log.Errorf("restCheckPermission error: [%s][%v]", err.Error(), r)
		return false, err
	}

	return check.Allowed, nil
}

func (c *Client) grpcCreateRelationTuple(ctx context.Context, namespace, object, relation, subject string) error {
	response, err := c.writeServiceClient.TransactRelationTuples(ctx, &acl.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*acl.RelationTupleDelta{
			{
				Action: acl.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &acl.RelationTuple{
					Namespace: namespace,
					Object:    object,
					Relation:  relation,
					Subject:   acl.NewSubjectID(subject),
				},
			},
		},
	})
	if err != nil {
		log.Errorf("grpcCreateRelationTuple error: [%s][%v]", err.Error(), response)
	}
	return err
}

func (c *Client) grpcGetCheck(ctx context.Context, namespace, object, relation, subject string) (bool, error) {
	response, err := c.checkServiceClient.Check(ctx, &acl.CheckRequest{
		Tuple: &acl.RelationTuple{
			Namespace: namespace,
			Object:    object,
			Relation:  relation,
			Subject:   acl.NewSubjectID(subject),
		},
	})
	if err != nil {
		// If namespace doesn't exist, we'll catch the Not Round error.
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		log.Errorf("grpcGetCheck error: [%s][%v]", err.Error(), response)
		return false, err
	}
	return response.Allowed, nil
}
