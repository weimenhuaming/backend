core-protos:
	goctl rpc protoc desc/core.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --style go_zero

core-mysql-user-login:
	goctl model mysql ddl -src desc/sql/user.sql -dir ./internal/model/user