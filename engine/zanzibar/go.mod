module github.com/tx7do/kratos-authz/engine/zanzibar

go 1.24.0

toolchain go1.24.3

require (
	github.com/go-kratos/kratos/v2 v2.8.4
	github.com/google/uuid v1.6.0
	github.com/openfga/go-sdk v0.7.1
	github.com/ory/keto-client-go v0.11.0-alpha.0
	github.com/ory/keto/proto v0.13.0-alpha.0
	github.com/stretchr/testify v1.11.1
	github.com/tx7do/kratos-authz v1.1.6
	google.golang.org/grpc v1.75.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250826171959-ef028d996bc1 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tx7do/kratos-authz => ../../
