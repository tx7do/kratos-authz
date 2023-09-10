module github.com/tx7do/kratos-authz/engine/zanzibar

go 1.19

replace github.com/tx7do/kratos-authz => ../../

require (
	github.com/go-kratos/kratos/v2 v2.7.0
	github.com/google/uuid v1.3.1
	github.com/openfga/go-sdk v0.2.2
	github.com/ory/keto-client-go v0.11.0-alpha.0
	github.com/ory/keto/proto v0.11.1-alpha.0
	github.com/stretchr/testify v1.8.4
	github.com/tx7do/kratos-authz v0.0.4
	google.golang.org/grpc v1.58.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/oauth2 v0.12.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
