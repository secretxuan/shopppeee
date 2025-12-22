#!/bin/bash

# Shoppee 一键启动脚本

set -e

echo "🚀 启动 Shoppee 电商系统..."
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose
if ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "⚠️  Node.js 未安装，前端将无法启动"
    FRONTEND_AVAILABLE=false
else
    FRONTEND_AVAILABLE=true
fi

# 启动后端服务
echo "1️⃣ 启动后端服务（PostgreSQL + Redis + Go API）..."
docker compose up -d

echo ""
echo "⏳ 等待后端服务启动..."
sleep 5

# 检查后端服务状态
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ 后端服务启动成功！"
else
    echo "⚠️  后端服务可能还在启动中，请稍候..."
fi

# 启动前端服务
if [ "$FRONTEND_AVAILABLE" = true ]; then
    echo ""
    echo "2️⃣ 启动前端服务..."
    
    cd frontend
    
    # 检查是否已安装依赖
    if [ ! -d "node_modules" ]; then
        echo "📦 安装前端依赖..."
        npm install
    fi
    
    echo "🎨 启动前端开发服务器..."
    npm run dev &
    FRONTEND_PID=$!
    
    cd ..
    
    echo ""
    echo "✅ 前端服务启动成功！"
else
    echo ""
    echo "⚠️  跳过前端服务启动"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Shoppee 电商系统启动完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📍 服务地址："
echo "   后端 API:  http://localhost:8080"
echo "   健康检查:  http://localhost:8080/health"
if [ "$FRONTEND_AVAILABLE" = true ]; then
    echo "   前端页面:  http://localhost:3000"
fi
echo ""
echo "📚 快速开始："
echo "   1. 访问前端: http://localhost:3000"
echo "   2. 注册账号或登录"
echo "   3. 浏览商品并添加到购物车"
echo ""
echo "🛠️  管理命令："
echo "   查看日志: docker compose logs -f app"
echo "   停止服务: docker compose down"
if [ "$FRONTEND_AVAILABLE" = true ]; then
    echo "   停止前端: kill $FRONTEND_PID"
fi
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 保持脚本运行（可选）
if [ "$FRONTEND_AVAILABLE" = true ]; then
    echo "💡 按 Ctrl+C 停止所有服务"
    wait $FRONTEND_PID
fi
