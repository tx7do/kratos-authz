module github.com/tx7do/kratos-authz/engine/zanzibar

go 1.19

require (
	github.com/go-kratos/kratos/v2 v2.7.1
	github.com/google/uuid v1.3.1
	github.com/openfga/go-sdk v0.2.3
	github.com/ory/keto-client-go v0.11.0-alpha.0
	github.com/ory/keto/proto v0.11.1-alpha.0
	github.com/stretchr/testify v1.8.4
	github.com/tx7do/kratos-authz v1.0.0
	google.golang.org/grpc v1.59.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.13.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231016165738-49dd2c1f3d0b // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tx7do/kratos-authz => ../../
