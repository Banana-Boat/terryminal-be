# go-grpc_out对应protoc-gen-go-grpc插件；go_out对应protoc-gen-go插件
proto:
	rm -rf ./internal/pb/*
	protoc --proto_path=./internal/proto --go_out=./internal/pb --go-grpc_out=./internal/pb \
	--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
	./internal/proto/*.proto

evans:
	evans --host localhost --port 8081 -r repl

redis:
	docker run --name terryminal-redis -p 6379:6379 -d redis:7-alpine

# github action执行的测试（通过-short flag区分）
remote_test:
	go clean -testcache && go test -v -cover -short ./...
	
# 本地执行的测试
local_test:
	go clean -testcache && go test -v -cover ./...

server:
	go run main.go

.PHONY: proto evans redis remote_test local_test server
