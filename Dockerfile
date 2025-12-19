# 多阶段构建 - 编译阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git make

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖（利用Docker层缓存）
RUN go mod download

# 复制源代码
COPY . .

# 编译Go程序（优化编译参数）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/shoppee \
    ./cmd/api/main.go

# 运行阶段 - 使用alpine最小化镜像
FROM alpine:latest

# 设置时区
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/shoppee /app/shoppee

# 复制配置文件（可选）
COPY --from=builder /app/.env.example /app/.env.example

# 创建日志目录
RUN mkdir -p /app/logs

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 运行程序
CMD ["/app/shoppee"]
