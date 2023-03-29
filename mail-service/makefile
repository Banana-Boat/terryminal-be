proto:
	rm -rf ./internal/pb/*.go
	protoc --proto_path=./internal/proto --go_out=./internal/pb --go_opt=paths=source_relative \
	--go-grpc_out=./internal/pb --go-grpc_opt=paths=source_relative \
	./internal/proto/*.proto

evans:
	evans --host localhost --port 8081 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

test:
	go test -v -cover -short ./...

server:
	go run main.go

.PHONY: proto evans redis test server
