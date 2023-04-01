DOCKER_USERNAME = tiangexiang
BASE_PTY_VERSION = 0.1.0


build_images:
	@echo "Building all the images..."
	docker build -f ./base-pty/Dockerfile -t ${DOCKER_USERNAME}/terryminal-base-pty:${BASE_PTY_VERSION} ./base-pty

remove_images:
	@echo "Removing all the images..."
	docker rmi ${DOCKER_USERNAME}/terryminal-base-pty:${BASE_PTY_VERSION}

# 构建支持多架构的镜像，并推到hub（镜像不保存到本地）
build_push_multi:
	@echo "Building and pushing all the images of multi-arch..."
	docker buildx build -t ${DOCKER_USERNAME}/terryminal-base-pty:${BASE_PTY_VERSION} --platform=linux/arm,linux/arm64,linux/amd64 --push ./base-pty

.PHONY: build_images remove_images build_push_multi
