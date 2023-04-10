default: help
.PHONY: help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

DOCKER_USERNAME = tiangexiang
MAIN_VERSION = 0.1.0
TERMINAL_VERSION = 0.1.0

.PHONY: create_network
create_network: ## 创建docker网络
	docker network create terryminal

.PHONY: build_images
build_images: ## 构建镜像
	@echo "Building all the images..."
	docker build -f ./main-service/Dockerfile -t ${DOCKER_USERNAME}/terryminal-main:${MAIN_VERSION} ./main-service
	docker build -f ./terminal-service/Dockerfile -t ${DOCKER_USERNAME}/terryminal-terminal:${TERMINAL_VERSION} ./terminal-service

.PHONY: remove_images
remove_images: ##	删除镜像
	@echo "Removing all the images..."
	docker rmi ${DOCKER_USERNAME}/terryminal-main:${MAIN_VERSION}
	docker rmi ${DOCKER_USERNAME}/terryminal-terminal:${TERMINAL_VERSION}

.PHONY: build_push_multi
build_push_multi: ## 构建支持多架构的镜像，并推到hub（镜像不保存到本地）
	@echo "Building and pushing all the images of multi-arch..."
	docker buildx build -t ${DOCKER_USERNAME}/terryminal-main:${MAIN_VERSION} --platform=linux/arm,linux/arm64,linux/amd64 --push ./main-service
	docker buildx build -t ${DOCKER_USERNAME}/terryminal-terminal:${TERMINAL_VERSION} --platform=linux/arm,linux/arm64,linux/amd64 --push ./terminal-service
