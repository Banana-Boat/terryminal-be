default: help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# =================================================测试相关=======================================================
.PHONY: remote_test
remote_test: ## github action执行的测试（通过-short flag区分）
	go clean -testcache && go test -v -cover -short ./...
	
.PHONY: local_test
local_test: ## 本地执行的测试
	go clean -testcache && go test -v -cover ./...

# =================================================服务运行（直接）=======================================================
.PHONY: run
run: ## 启动服务
	go run main.go

# =================================================服务运行（docker）=======================================================
.PHONY: build_image
build_image: ## 构建镜像
	docker build -t terryminal-chatbot:dev .

.PHONY: run_docker
run_docker: ## 启动服务（docker）
	docker run --name terryminal-chatbot --network terryminal -p 3230:3230 terryminal-chatbot:dev

.PHONY: clean_docker
clean_docker: ## 停止容器并删除容器与镜像
	docker stop terryminal-chatbot
	docker rm terryminal-chatbot
	docker rmi terryminal-chatbot:dev
