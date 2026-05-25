core-protos:
	goctl rpc protoc desc/core.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --style go_zero

core-mysql-user-login:
	goctl model mysql ddl -src desc/sql/user.sql -dir ./internal/model/user

core-mysql-user-article:
	goctl model mysql ddl -src desc/sql/article.sql -dir ./internal/model/article

core-mysql-user-comment:
	goctl model mysql ddl -src desc/sql/interaction_comment.sql -dir ./internal/model/comment

core-mysql-user-like:
	goctl model mysql ddl -src desc/sql/interaction_like.sql -dir ./internal/model/like

core-mysql-user-favor:
	goctl model mysql ddl -src desc/sql/interaction_favor.sql -dir ./internal/model/favor

core-mysql-user-category:
	goctl model mysql ddl -src desc/sql/category.sql -dir ./internal/model/category