module github.com/tx7do/kratos-authz/middleware

go 1.24.6

replace (
	github.com/tx7do/kratos-authz => ../
)

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/tx7do/kratos-authz v1.1.6
)

require (
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251213004720-97cd9d5aeac2 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
