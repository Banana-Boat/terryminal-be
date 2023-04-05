# 使用 multi-stage 进行构建，进一步减小镜像大小

# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY ./app.env .

EXPOSE 3210 3211

CMD [ "/app/main" ]