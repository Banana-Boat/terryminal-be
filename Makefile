DOCKER_USERNAME = tiangexiang
MAIN_VERSION = 0.1.0
TERMINAL_VERSION = 0.1.0

create_network:
	docker network create terryminal

build_images:
	@echo "Building all the images..."
	docker build -f ./main-service/Dockerfile -t ${DOCKER_USERNAME}/terryminal-main:${MAIN_VERSION} ./main-service
	docker build -f ./terminal-service/Dockerfile -t ${DOCKER_USERNAME}/terryminal-terminal:${TERMINAL_VERSION} ./terminal-service

remove_images:
	@echo "Removing all the images..."
	docker rmi ${DOCKER_USERNAME}/terryminal-main:${MAIN_VERSION}
	docker rmi ${DOCKER_USERNAME}/terryminal-terminal:${TERMINAL_VERSION}

# 构建支持多架构的镜像，并推到hub（镜像不保存到本地）
build_push_multi:
	@echo "Building and pushing all the images of multi-arch..."
	docker buildx build -t ${DOCKER_USERNAME}/terryminal-main:${MAIN_VERSION} --platform=linux/arm,linux/arm64,linux/amd64 --push ./main-service
	docker buildx build -t ${DOCKER_USERNAME}/terryminal-terminal:${TERMINAL_VERSION} --platform=linux/arm,linux/arm64,linux/amd64 --push ./terminal-service

.PHONY: create_network build_images remove_images build_push_multi
