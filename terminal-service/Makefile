default: help
.PHONY: help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: proto
# go-grpc_out对应protoc-gen-go-grpc插件；go_out对应protoc-gen-go插件
proto: ## 根据proto文件生成go代码
	rm -rf ./internal/pb/*
	protoc --proto_path=./internal/proto --go_out=./internal/pb --go-grpc_out=./internal/pb \
	--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
	./internal/proto/*.proto

.PHONY: evans
evans: ##	测试gRPC服务
	evans --host localhost --port 8081 -r repl

.PHONY: redis
redis: ## 启动redis
	docker run --name terryminal-redis -p 6379:6379 -d redis:7-alpine

.PHONY: remote_test
remote_test: ## github action执行的测试（通过-short flag区分）
	go clean -testcache && go test -v -cover -short ./...
	
.PHONY: local_test
local_test: ## 本地执行的测试
	go clean -testcache && go test -v -cover ./...

.PHONY: server
server: ## 启动服务
	go run main.go