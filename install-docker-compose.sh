#!/bin/bash

# 安装 docker-compose 脚本

set -e

echo "📦 开始安装 docker-compose..."

# 下载 docker-compose
echo "1️⃣ 下载 docker-compose 1.29.2..."
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" \
  -o /usr/local/bin/docker-compose

# 添加执行权限
echo "2️⃣ 添加执行权限..."
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
echo "3️⃣ 验证安装..."
docker-compose --version

echo ""
echo "✅ docker-compose 安装成功！"
echo ""
echo "现在可以运行："
echo "  docker-compose up -d"
