default: help
.PHONY: help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	
# =================================================测试相关=======================================================
.PHONY: remote_test
remote_test: ## github action执行的测试（通过-short flag区分）
	go clean -testcache && go test -v -cover -short ./...
	
.PHONY: local_test
local_test: ## 本地执行的测试
	go clean -testcache && go test -v -cover ./...

# =================================================服务运行（直接）=======================================================
.PHONY: redis
redis: ## 启动redis
	docker run --name terryminal-bot-redis -p 6379:6379 -d redis:7-alpine

.PHONY: run
run: ## 启动服务
	go run main.go

.PHONY: clean
clean: ## 清理
	docker stop terryminal-bot-redis
	docker rm terryminal-bot-redis

# =================================================服务运行（docker）=======================================================
.PHONY: build_image
build_image: ## 构建镜像
	docker build -t terryminal-bot:dev .

.PHONY: run_docker
run_docker: ## 启动服务（docker）
	docker run --name terryminal-bot-redis -p 6379:6379 -d redis:7-alpine
	docker run --name terryminal-bot --network terryminal -p 3230:3230 terryminal-bot:dev

.PHONY: clean_docker
clean_docker: ## 停止容器并删除容器与镜像
	docker stop terryminal-bot
	docker rm terryminal-bot
	docker rmi terryminal-bot:dev
	docker stop terryminal-bot-redis
	docker rm terryminal-bot-redis
