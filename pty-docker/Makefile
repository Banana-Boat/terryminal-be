default: help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

DOCKER_USERNAME = tiangexiang
BASE_PTY_VERSION = 0.1.0

.PHONY: evans
evans: ## 测试gRPC服务
	evans --proto ./base-pty/proto/base_pty.proto --host localhost --port 3220 repl

.PHONY: build_images
build_images: ## 构建镜像
	@echo "Building all the images..."
# node-base-pty需要打开注释
# npm --prefix ./base-pty/ run build
	docker build -f ./base-pty/Dockerfile -t ${DOCKER_USERNAME}/terryminal-base-pty:${BASE_PTY_VERSION} ./base-pty

.PHONY: remove_images
remove_images: ## 删除镜像
	@echo "Removing all the images..."
	docker rmi ${DOCKER_USERNAME}/terryminal-base-pty:${BASE_PTY_VERSION}

.PHONY: build_push_multi
build_push_multi: ## 构建支持多架构的镜像，并推到hub（镜像不保存到本地）
	@echo "Building and pushing all the images of multi-arch..."
	docker buildx build -t ${DOCKER_USERNAME}/terryminal-base-pty:${BASE_PTY_VERSION} --platform=linux/arm,linux/arm64,linux/amd64 --push ./base-pty