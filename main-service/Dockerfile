# 使用 multi-stage 进行构建，进一步减小镜像大小

# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN GOPROXY=https://goproxy.cn go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY ./app.env .
COPY ./wait-for.sh .
# <src> 是一个目录，则将目录下的所有文件写入<dest>中
COPY ./internal/db/migration ./migration

EXPOSE 3200

CMD [ "/app/main" ]