default: help
.PHONY: help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: mysql
mysql: ## 启动mysql容器（需要指定DB_PASSWORD与DB_NAME）
ifeq (${DB_PASSWORD}${DB_NAME},)
	@echo "missing parameter DB_PASSWORD and DB_NAME"
else
	docker run --name terryminal-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=${DB_PASSWORD} -e MYSQL_DATABASE=${DB_NAME} -d mysql:8.0
endif

.PHONY: migrate_init
migrate_init:	## 生成migrate模版文件
	rm -rf ./internal/db/migration/** && migrate create -ext sql -dir ./internal/db/migration -seq init_schema

.PHONY: sqlc
sqlc: ## 根据sql文件生成go代码
	sqlc generate

.PHONY: proto
# go-grpc_out对应protoc-gen-go-grpc插件；go_out对应protoc-gen-go插件
proto: ## 根据proto文件生成go代码
	rm -rf ./internal/pb/*
	protoc --proto_path=./internal/proto --go_out=./internal/pb --go-grpc_out=./internal/pb \
	--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
	./internal/proto/*.proto

.PHONY: remote_test
remote_test: ## github action执行的测试（通过-short flag区分）
	go clean -testcache && go test -v -cover -short ./...
	
.PHONY: local_test
local_test: ## 本地执行的测试
	go clean -testcache && go test -v -cover ./...

.PHONY: server
server: ## 启动服务
	go run main.go
