module github.com/tx7do/kratos-authz/engine/casbin

go 1.23.0

toolchain go1.24.3

require (
	github.com/casbin/casbin/v2 v2.107.0
	github.com/stretchr/testify v1.10.0
	github.com/tx7do/kratos-authz v1.0.3
)

require (
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/casbin/govaluate v1.7.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/grpc v1.73.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tx7do/kratos-authz => ../../
