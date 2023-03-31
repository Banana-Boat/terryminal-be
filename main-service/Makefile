mysql:
ifeq (${DB_PASSWORD}${DB_NAME},)
	@echo "missing parameter DB_PASSWORD and DB_NAME"
else
	docker run --name terryminal-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=${DB_PASSWORD} -e MYSQL_DATABASE=${DB_NAME} -d mysql:8.0
endif

migrate_init:
	rm -rf ./internal/db/migration/** && migrate create -ext sql -dir ./internal/db/migration -seq init_schema

sqlc:
	sqlc generate

proto:
	rm -rf ./internal/pb/*.go
	protoc --proto_path=./internal/proto --go_out=./internal/pb --go_opt=paths=source_relative \
	--go-grpc_out=./internal/pb --go-grpc_opt=paths=source_relative \
	./internal/proto/*.proto

# github action执行的测试（通过-short flag区分）
remote_test:
	go clean -testcache && go test -v -cover -short ./...
	
# 本地执行的测试
local_test:
	go clean -testcache && go test -v -cover ./...

server:
	go run main.go


.PHONY: mysql migrate_init sqlc proto remote_test local_test server
