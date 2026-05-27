core-protos:
	goctl rpc protoc desc/core.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --style go_zero

# 数据库表结构由 GORM entity + AutoMigrate 管理，见 internal/model/entity 与 internal/model/init.go
