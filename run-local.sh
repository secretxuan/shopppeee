#!/bin/bash

# Shoppee 本地运行脚本（不使用Docker）

set -e

echo "🚀 启动 Shoppee 电商系统（本地模式）"
echo "=========================================="

# 1. 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未安装 Go"
    exit 1
fi
echo "✅ Go 版本: $(go version)"

# 2. 加载环境变量
if [ -f .env.local ]; then
    export $(cat .env.local | grep -v '^#' | xargs)
    echo "✅ 已加载 .env.local 配置"
else
    echo "⚠️  警告: .env.local 不存在，使用默认配置"
fi

# 3. 检查PostgreSQL连接
echo ""
echo "📦 检查依赖服务..."
if command -v psql &> /dev/null; then
    if psql -h ${DB_HOST:-localhost} -U ${DB_USER:-postgres} -d postgres -c '\q' 2>/dev/null; then
        echo "✅ PostgreSQL 连接成功"
    else
        echo "⚠️  PostgreSQL 连接失败，请确保："
        echo "   1. PostgreSQL 已安装并运行"
        echo "   2. 创建数据库: createdb -U postgres shoppee"
        echo "   3. 或使用远程数据库并修改 .env.local"
    fi
else
    echo "⚠️  未安装 psql 客户端"
fi

# 4. 检查Redis连接
if command -v redis-cli &> /dev/null; then
    if redis-cli -h ${REDIS_HOST:-localhost} ping 2>/dev/null | grep -q PONG; then
        echo "✅ Redis 连接成功"
    else
        echo "⚠️  Redis 连接失败，请确保 Redis 已启动"
    fi
else
    echo "⚠️  未安装 redis-cli 客户端"
fi

# 5. 创建日志目录
mkdir -p logs
echo "✅ 日志目录已创建"

# 6. 下载依赖
echo ""
echo "📥 下载依赖..."
go mod download
echo "✅ 依赖下载完成"

# 7. 运行应用
echo ""
echo "=========================================="
echo "🎯 启动应用服务器..."
echo "=========================================="
echo ""

go run cmd/api/main.go
