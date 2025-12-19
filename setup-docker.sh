#!/bin/bash

# Docker 完整设置脚本（适用于容器环境）

set -e

echo "🐳 设置 Docker 环境"
echo "========================================"

# 1. 安装 docker-compose
if ! command -v docker-compose &> /dev/null; then
    echo "📦 安装 docker-compose..."
    sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" \
      -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo "✅ docker-compose 安装完成"
else
    echo "✅ docker-compose 已安装: $(docker-compose --version)"
fi

echo ""
echo "🔧 检查 Docker daemon..."

# 2. 检查 Docker daemon 是否运行
if docker info &> /dev/null; then
    echo "✅ Docker daemon 正在运行"
else
    echo "❌ Docker daemon 未运行"
    echo ""
    echo "⚠️  你的环境限制："
    echo "   - 在容器内运行 Docker (Docker-in-Docker)"
    echo "   - 需要特权模式或挂载 Docker socket"
    echo ""
    echo "🔧 可能的解决方案："
    echo ""
    echo "方案1: 如果主机有 Docker，挂载 socket"
    echo "   docker run -v /var/run/docker.sock:/var/run/docker.sock ..."
    echo ""
    echo "方案2: 使用特权模式启动容器"
    echo "   docker run --privileged ..."
    echo ""
    echo "方案3: 不使用 Docker，直接本地运行"
    echo "   ./run-local.sh"
    echo ""
    exit 1
fi

echo ""
echo "========================================"
echo "✅ Docker 环境设置完成！"
echo ""
echo "现在可以运行："
echo "  cd /data/workspace/shoppee"
echo "  docker-compose up -d"
echo "========================================"
