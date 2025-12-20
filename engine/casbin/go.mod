module github.com/tx7do/kratos-authz/engine/casbin

go 1.24.0

toolchain go1.24.3

replace github.com/tx7do/kratos-authz => ../../

require (
	github.com/casbin/casbin/v2 v2.135.0
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/stretchr/testify v1.11.1
	github.com/tx7do/kratos-authz v1.1.7
)

require (
	github.com/bmatcuk/doublestar/v4 v4.9.1 // indirect
	github.com/casbin/govaluate v1.10.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251213004720-97cd9d5aeac2 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
