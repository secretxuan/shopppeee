#!/bin/bash

# Shoppee 完整项目启动脚本

echo "🚀 启动 Shoppee 电商系统"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 检查是否在项目根目录
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ 错误：请在项目根目录运行此脚本"
    exit 1
fi

# 1. 启动后端服务
echo "1️⃣ 启动后端服务（PostgreSQL + Redis + Go API）..."
echo ""

# 检查 Docker 权限
if ! docker ps &> /dev/null; then
    echo "⚠️  需要 sudo 权限来运行 Docker"
    sudo docker compose up -d
else
    docker compose up -d
fi

echo ""
echo "⏳ 等待后端服务启动（10秒）..."
sleep 10

# 检查后端健康状态
echo ""
echo "🔍 检查后端服务状态..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ 后端服务运行正常！"
else
    echo "⚠️  后端服务可能还在启动中..."
    echo "   可以运行以下命令查看日志："
    echo "   sudo docker compose logs -f app"
fi

# 2. 启动前端服务
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2️⃣ 启动前端服务..."
echo ""

cd frontend

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "📦 首次运行，正在安装依赖..."
    npm install
    echo ""
fi

echo "🎨 启动前端开发服务器..."
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 服务启动完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📍 访问地址："
echo "   🎨 前端页面: http://localhost:3000"
echo "   🔧 后端 API: http://localhost:8080"
echo "   ❤️  健康检查: http://localhost:8080/health"
echo ""
echo "💡 提示："
echo "   - 前端开发服务器支持热重载"
echo "   - 按 Ctrl+C 可停止前端服务"
echo "   - 停止后端: sudo docker compose down"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 启动前端（在前台运行）
npm run dev
