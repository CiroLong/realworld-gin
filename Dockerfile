# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o realworld-server ./cmd/server

# 运行阶段
FROM alpine:3.18

WORKDIR /app

# 复制构建产物
COPY --from=builder /app/realworld-server .

# 复制配置文件目录
COPY --from=builder /app/config ./config

# 暴露端口
EXPOSE 8000

# 运行应用
CMD ["./realworld-server"]