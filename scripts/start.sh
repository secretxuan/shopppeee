#!/bin/bash

# Shoppee 快速启动脚本
# 用法: ./scripts/start.sh [dev|prod]

set -e

ENV=${1:-dev}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "================================================"
echo "  Shoppee 电商系统启动脚本"
echo "  环境: $ENV"
echo "================================================"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查依赖
check_dependencies() {
    echo -e "\n${YELLOW}[1/5]${NC} 检查依赖..."
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}错误: Docker 未安装${NC}"
        exit 1
    fi
    
    # 检查 Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}错误: Docker Compose 未安装${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 依赖检查通过${NC}"
}

# 配置环境变量
setup_env() {
    echo -e "\n${YELLOW}[2/5]${NC} 配置环境变量..."
    
    if [ ! -f .env ]; then
        echo -e "${YELLOW}未找到 .env 文件，从 .env.example 复制...${NC}"
        cp .env.example .env
        echo -e "${GREEN}✓ .env 文件已创建${NC}"
        echo -e "${RED}警告: 请修改 .env 中的密码和密钥！${NC}"
    else
        echo -e "${GREEN}✓ .env 文件已存在${NC}"
    fi
}

# 启动开发环境
start_dev() {
    echo -e "\n${YELLOW}[3/5]${NC} 启动开发环境..."
    
    # 启动数据库和 Redis
    docker-compose up -d postgres redis
    
    echo -e "${GREEN}✓ PostgreSQL 和 Redis 已启动${NC}"
    echo -e "\n${YELLOW}等待数据库就绪...${NC}"
    sleep 5
    
    # 下载依赖
    echo -e "\n${YELLOW}[4/5]${NC} 下载 Go 依赖..."
    go mod download
    echo -e "${GREEN}✓ 依赖下载完成${NC}"
    
    # 运行应用
    echo -e "\n${YELLOW}[5/5]${NC} 启动应用..."
    echo -e "${GREEN}应用将在 http://localhost:8080 启动${NC}"
    echo -e "${YELLOW}按 Ctrl+C 停止${NC}\n"
    
    go run cmd/api/main.go
}

# 启动生产环境
start_prod() {
    echo -e "\n${YELLOW}[3/5]${NC} 构建 Docker 镜像..."
    docker-compose build
    
    echo -e "\n${YELLOW}[4/5]${NC} 启动所有服务..."
    docker-compose up -d
    
    echo -e "\n${YELLOW}[5/5]${NC} 检查服务状态..."
    sleep 3
    docker-compose ps
    
    echo -e "\n${GREEN}✓ 所有服务已启动${NC}"
    echo -e "\n访问地址："
    echo -e "  应用: ${GREEN}http://localhost:8080${NC}"
    echo -e "  健康检查: ${GREEN}http://localhost:8080/health${NC}"
    echo -e "\n查看日志："
    echo -e "  ${YELLOW}docker-compose logs -f app${NC}"
    echo -e "\n停止服务："
    echo -e "  ${YELLOW}docker-compose down${NC}"
}

# 主函数
main() {
    check_dependencies
    setup_env
    
    if [ "$ENV" = "dev" ]; then
        start_dev
    elif [ "$ENV" = "prod" ]; then
        start_prod
    else
        echo -e "${RED}错误: 未知环境 '$ENV'${NC}"
        echo "用法: $0 [dev|prod]"
        exit 1
    fi
}

main
