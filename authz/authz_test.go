package authz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	"io"
	"os"
	"testing"

	jwtV4 "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"github.com/tx7do/kratos-authz/engine"
	"github.com/tx7do/kratos-authz/engine/casbin"
)

const (
	ClaimAuthorityId = "authorityId"
	ClaimDomain      = "domain"
)

type myTransport struct {
	transport.Transporter
	kind      transport.Kind
	endpoint  string
	operation string
	method    string
	reqHeader transport.Header
}

func (tr *myTransport) Kind() transport.Kind {
	return tr.kind
}

func (tr *myTransport) Endpoint() string {
	return tr.endpoint
}

func (tr *myTransport) Operation() string {
	return tr.operation
}

func (tr *myTransport) Method() string {
	return tr.method
}

func (tr *myTransport) RequestHeader() transport.Header {
	return tr.reqHeader
}

func (tr *myTransport) ReplyHeader() transport.Header {
	return nil
}

type mySecurityUser struct {
	Path        string
	Method      string
	AuthorityId string
	Domain      string
}

func NewSecurityUser() SecurityUser {
	return &mySecurityUser{}
}

func (su *mySecurityUser) ParseFromContext(ctx context.Context) error {
	if claims, ok := jwt.FromContext(ctx); ok {
		str, ok := claims.(jwtV4.MapClaims)[ClaimAuthorityId]
		if ok {
			su.AuthorityId = str.(string)
		}
		str, ok = claims.(jwtV4.MapClaims)[ClaimDomain]
		if ok {
			su.Domain = str.(string)
		}
		//str, ok = claims.(jwtV4.MapClaims)[ClaimMethod]
		//if ok {
		//	su.Method = str.(string)
		//}
	} else {
		return errors.New("jwt claim missing")
	}

	if header, ok := transport.FromServerContext(ctx); ok {
		myHeader := header.(*myTransport)
		su.Path = header.Operation()
		su.Method = myHeader.Method()
	} else {
		return errors.New("jwt claim missing")
	}

	return nil
}

func (su *mySecurityUser) GetSubject() string {
	return su.AuthorityId
}

func (su *mySecurityUser) GetObject() string {
	return su.Path
}

func (su *mySecurityUser) GetAction() string {
	return su.Method
}

func (su *mySecurityUser) GetDomain() string {
	return su.Domain
}

func createToken(authorityId, domain string) jwtV4.Claims {
	return jwtV4.MapClaims{
		ClaimAuthorityId: authorityId,
		ClaimDomain:      domain,
	}
}

func TestServer_Casbin(t *testing.T) {
	policies := map[string]interface{}{
		"policies": []casbin.PolicyRule{
			{PType: "p", V0: "bobo", V1: "/api/*", V2: "ANY", V3: "*"},
			{PType: "p", V0: "bobo01", V1: "/api/users", V2: "ANY", V3: "*"},
			{PType: "p", V0: "admin_role", V1: "/api/*", V2: "ANY", V3: "*"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "*"},
		},
		//"projects": allProjects,
	}

	tests := []struct {
		name        string
		authorityId string
		path        string
		exceptErr   error
	}{
		{
			authorityId: "admin",
			path:        "/api/login",
			exceptErr:   nil,
		},
		{
			authorityId: "admin",
			path:        "/api/logout",
			exceptErr:   nil,
		},
		{
			authorityId: "admin",
			path:        "/api/logout:hell",
			exceptErr:   nil,
		},
		{
			authorityId: "admin",
			path:        "/api/logout/login",
			exceptErr:   nil,
		},
		{
			authorityId: "admin",
			path:        "/api1/logout",
			exceptErr:   ErrUnauthorized,
		},
		{
			authorityId: "bobo",
			path:        "/api/login",
			exceptErr:   nil,
		},
		{
			authorityId: "bobo01",
			path:        "/api/users",
			exceptErr:   nil,
		},
		{
			authorityId: "bobo01",
			path:        "/api/dept",
			exceptErr:   ErrUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.authorityId, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				//t.Log(req)
				return "reply", nil
			}

			token := createToken(test.authorityId, "")
			ctx := transport.NewServerContext(context.Background(), &myTransport{operation: test.path, method: "ANY"})
			ctx = jwt.NewContext(ctx, token)

			var server middleware.Handler
			server = Server(
				WithPolicyEngine(ctx, engine.CasbinEngine),
				WithPolicies(ctx, policies, nil),
				WithSecurityUserCreator(NewSecurityUser),
			)(next)

			_, err := server(ctx, "request")
			assert.EqualValues(t, test.exceptErr, err)
		})
	}
}

func TestServer_CasbinWithDomain(t *testing.T) {
	policies := map[string]interface{}{
		"policies": []casbin.PolicyRule{
			{PType: "p", V0: "bobo", V1: "/api/*", V2: "ANY", V3: "*"},
			{PType: "p", V0: "bobo01", V1: "/api/users", V2: "ANY", V3: "project1"},
			{PType: "p", V0: "admin_role", V1: "/api/*", V2: "ANY", V3: "*"},
			{PType: "g", V0: "admin", V1: "admin_role", V2: "*"},
		},
		//"projects": allProjects,
	}

	tests := []struct {
		name        string
		authorityId string
		domain      string
		path        string
		exceptErr   error
	}{
		{
			authorityId: "admin",
			path:        "/api/login",
			domain:      "*",
			exceptErr:   nil,
		},
		{
			authorityId: "admin",
			path:        "/api/logout",
			domain:      "*",
			exceptErr:   nil,
		},
		{
			authorityId: "admin",
			path:        "/api/logout:hell",
			domain:      "*",
			exceptErr:   nil,
		},
		{
			authorityId: "admin",
			path:        "/api/logout/login",
			domain:      "*",
			exceptErr:   nil,
		},
		{
			authorityId: "admin",
			path:        "/api1/logout",
			domain:      "*",
			exceptErr:   ErrUnauthorized,
		},
		{
			authorityId: "bobo",
			path:        "/api/login",
			domain:      "*",
			exceptErr:   nil,
		},
		{
			authorityId: "bobo01",
			path:        "/api/users",
			domain:      "*",
			exceptErr:   ErrUnauthorized,
		},
		{
			authorityId: "bobo01",
			path:        "/api/users",
			domain:      "project1",
			exceptErr:   nil,
		},
		{
			authorityId: "bobo01",
			path:        "/api/users",
			domain:      "project2",
			exceptErr:   ErrUnauthorized,
		},
		{
			authorityId: "bobo01",
			path:        "/api/users1",
			domain:      "project1",
			exceptErr:   ErrUnauthorized,
		},
		{
			authorityId: "bobo01",
			path:        "/api/dept",
			domain:      "*",
			exceptErr:   ErrUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.authorityId, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				//t.Log(req)
				return "reply", nil
			}

			token := createToken(test.authorityId, test.domain)
			ctx := transport.NewServerContext(context.Background(), &myTransport{operation: test.path, method: "ANY"})
			ctx = jwt.NewContext(ctx, token)

			var server middleware.Handler
			server = Server(
				WithDomainSupport(),
				WithPolicyEngine(ctx, engine.CasbinEngine),
				WithPolicies(ctx, policies, nil),
				WithSecurityUserCreator(NewSecurityUser),
			)(next)

			_, err := server(ctx, "request")
			assert.EqualValues(t, test.exceptErr, err)
		})
	}
}

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

func baselinePoliciesAndRoles() (policies map[string]interface{}, roles map[string]interface{}) {
	// this file includes system, migrated legacy, and chef-managed policies
	// and chef-managed roles
	jsonFile, err := os.Open("../engine/opa/example/real_world_store.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var pr struct {
		Policies map[string]interface{} `json:"policies"`
		Roles    map[string]interface{} `json:"roles"`
	}
	_ = json.Unmarshal(byteValue, &pr)

	return pr.Policies, pr.Roles
}

func TestServer_OPA(t *testing.T) {
	policies, roles := baselinePoliciesAndRoles()

	tests := []struct {
		authorityId string
		path        string
		method      string
		domain      string

		exceptErr error
	}{
		{
			authorityId: "user:local:test",
			path:        "system:status",
			method:      "system:license:get",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   nil,
		},
		{
			authorityId: "user:local:test@example.com",
			path:        "system:status",
			method:      "system:license:get",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   nil,
		},
		{
			authorityId: "user:local:test@example.com",
			path:        "iam:users:test@example.com",
			method:      "iam:users:get",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.authorityId, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				//t.Log(req)
				return "reply", nil
			}

			token := createToken(test.authorityId, test.domain)
			ctx := transport.NewServerContext(context.Background(), &myTransport{operation: test.path, method: test.method})
			ctx = jwt.NewContext(ctx, token)

			var server middleware.Handler
			server = Server(
				WithPolicyEngine(ctx, engine.OpaEngine),
				WithPolicies(ctx, policies, roles),
				WithSecurityUserCreator(NewSecurityUser),
			)(next)

			_, err := server(ctx, "request")
			assert.EqualValues(t, test.exceptErr, err)
		})
	}
}

func TestServer_OPAWithDomain(t *testing.T) {
	//policyCount := 20
	//roleCount := 10
	//policies, roles := baselineAndRandomPoliciesAndRoles(policyCount, roleCount)

	policies, roles := baselinePoliciesAndRoles()

	//fmt.Println(policies, roles)

	tests := []struct {
		authorityId string
		path        string
		method      string
		domain      string

		exceptErr error
	}{
		{
			authorityId: "user:local:test",
			path:        "compliance:profiles",
			method:      "compliance:profiles:list",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   ErrUnauthorized,
		},
		{
			authorityId: "tls:service:automate-cs-nginx:test",
			path:        "compliance:profiles",
			method:      "compliance:profiles:list",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   nil,
		},
		{
			authorityId: "user:local:admin",
			path:        "system:status",
			method:      "system:license:get",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   nil,
		},
		{
			authorityId: "user:local:admin",
			path:        "system:license",
			method:      "system:license:get",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   ErrUnauthorized,
		},
		{
			authorityId: "user:local:admin",
			path:        "system:license",
			method:      "system:license:get",
			domain:      "project1",
			exceptErr:   ErrUnauthorized,
		},
		{
			authorityId: "user:local:admin",
			path:        "iam:users:admin",
			method:      "iam:users:get",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   nil,
		},
		{
			authorityId: "user:local:admin",
			path:        "iam:users:admin",
			method:      "iam:users:delete",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   ErrUnauthorized,
		},
		{
			authorityId: "user:local:admin",
			path:        "iam:introspect",
			method:      "iam:introspect:get",
			domain:      "~~ALL-PROJECTS~~",
			exceptErr:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.authorityId, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				//t.Log(req)
				return "reply", nil
			}

			token := createToken(test.authorityId, test.domain)
			ctx := transport.NewServerContext(context.Background(), &myTransport{operation: test.path, method: test.method})
			ctx = jwt.NewContext(ctx, token)

			var server middleware.Handler
			server = Server(
				WithDomainSupport(),
				WithPolicyEngine(ctx, engine.OpaEngine),
				WithPolicies(ctx, policies, roles),
				WithSecurityUserCreator(NewSecurityUser),
			)(next)

			_, err := server(ctx, "request")
			assert.EqualValues(t, test.exceptErr, err)
		})
	}
}
