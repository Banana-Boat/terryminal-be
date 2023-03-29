# go-micro-template

基于 Go 的微服务架构后端模版。主服务使用 Gin 向外提供 API，各服务间通过 gRPC 进行通信。预置 Mysql 与 Redis 用作数据存储。

## CLI 工具

- [**Docker**](https://hub.docker.com/)
- [**golang-migrate**](https://github.com/golang-migrate/migrate)（数据库迁移）
- [**sqlc**](https://docs.sqlc.dev/en/stable/index.html)（sql->go 接口）
- [**Protocol Buffers**](https://protobuf.dev)（gRPC 数据定义）
- [**Evans**](https://github.com/ktr0731/evans)（gRPC 调试工具）

## 主要依赖

- [**gRPC**](https://grpc.io/)
- [**Protocol Buffers**](https://protobuf.dev)（gRPC 数据定义）
- [**Asynq**](https://github.com/hibiken/asynq)（任务队列异步处理框架）
- [**golang-migrate**](https://github.com/golang-migrate/migrate)（数据库迁移）
- [**Testify**](https://github.com/stretchr/testify)（测试框架）
- [**Viper**](https://github.com/spf13/viper)（配置项管理）
- [**Paseto**](https://github.com/o1egl/paseto)（用户鉴权）
- [**Zerolog**](https://github.com/rs/zerolog)（日志输出）

## 项目信息修改

- 修改 go.mod 中 module 名（全局替换）
- 向 .gitignore 文件中添加 app.env 与 compose.yaml

## 开发场景

#### 基本环境

- 安装 CLI 工具
- 安装 Go 依赖 `go mod tidy`

#### 数据库

- Mysql

  - 执行`migrate_init`生成 schema
  - 使用 [**dbdiagram**](https://dbdiagram.io/home) 工具设计数据库，将 sql 语句复制到上一步的 schema 中
  - 执行`make mysql DB_PASSWORD=? DB_NAME=?`，启动 mysql 容器

- Redis
  - 执行`make redis`，启动 mysql 容器

#### sqlc

- 在 internal/db/query/ 下创建 表名.sql 文件，根据官网编写 sql 语句
- 执行`make sqlc`生成.go 文件

#### protoc

- 修改 internal/proto 目录下的文件名以及文件内容
- 执行 `make proto`

#### 编译运行

- 执行`make server`，编译运行

## 部署场景

- 修改根目录下 compose.yaml 与 makefile 相关信息
- 应用容器化
  - 方法一：在服务器端拉取代码，执行`make build_images`，打包镜像
  - 方法二：本地执行`make build_push_multi`，打包多平台镜像并推至 hub

#### Swarm 方式部署（支持多节点集群）

> 参考 [**docker swarm**](https://docs.docker.com/engine/reference/commandline/swarm/)、[**docker service**](https://docs.docker.com/engine/reference/commandline/service/)、[**docker stack**](https://docs.docker.com/engine/reference/commandline/stack/)

注意！！！：Swarm 模式下数据不支持通过 Volume 挂载到 Host 进行持久化存储，此处待改进

- `docker swarm init` 创建集群
- `docker stack deploy -c compose.yaml ???`（??? 为项目名）部署 stack
- `docker service ls` 查看 service 列表
- `docker service logs SERVICE` 查看某个 service 的日志

- `docker stack rm ???` 结束 stack
- `docker swarm leave` 离开集群

#### Compose 方式部署（仅支持单节点）

> 参考 [**docker compose**](https://docs.docker.com/engine/reference/commandline/compose/)

- `docker compose up -d` 后台运行
- `docker compose logs` 查看日志
- `docker compose down` 结束运行并删除容器
