default: help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# =================================================gRPC相关=======================================================
.PHONY: proto
proto: ## 根据proto文件生成python代码
	rm -rf ./*_pb2.py ./*_pb2_grpc.py ./*.pyi
	python -m grpc_tools.protoc --proto_path=./proto --python_out=. --pyi_out=. --grpc_python_out=. ./proto/*.proto
