#!/bin/bash

# API 测试脚本
# 用法: ./scripts/test_api.sh

set -e

BASE_URL="http://localhost:8080"
API_URL="${BASE_URL}/api/v1"

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "================================================"
echo "  Shoppee API 测试脚本"
echo "================================================"

# 测试健康检查
test_health() {
    echo -e "\n${YELLOW}[1/6] 测试健康检查...${NC}"
    
    response=$(curl -s "${BASE_URL}/health")
    echo "响应: $response"
    
    if [[ $response == *"ok"* ]]; then
        echo -e "${GREEN}✓ 健康检查通过${NC}"
    else
        echo -e "${RED}✗ 健康检查失败${NC}"
        exit 1
    fi
}

# 测试用户注册
test_register() {
    echo -e "\n${YELLOW}[2/6] 测试用户注册...${NC}"
    
    timestamp=$(date +%s)
    username="testuser${timestamp}"
    
    response=$(curl -s -X POST "${API_URL}/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"${username}\",
            \"email\": \"${username}@example.com\",
            \"password\": \"password123\"
        }")
    
    echo "响应: $response"
    
    if [[ $response == *"success"* ]] || [[ $response == *"注册成功"* ]]; then
        echo -e "${GREEN}✓ 用户注册成功${NC}"
    else
        echo -e "${RED}✗ 用户注册失败${NC}"
    fi
}

# 测试用户登录
test_login() {
    echo -e "\n${YELLOW}[3/6] 测试用户登录...${NC}"
    
    # 先注册一个用户
    timestamp=$(date +%s)
    username="logintest${timestamp}"
    
    curl -s -X POST "${API_URL}/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"${username}\",
            \"email\": \"${username}@example.com\",
            \"password\": \"password123\"
        }" > /dev/null
    
    # 登录
    response=$(curl -s -X POST "${API_URL}/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"${username}\",
            \"password\": \"password123\"
        }")
    
    echo "响应: $response"
    
    # 提取 token
    TOKEN=$(echo $response | grep -o '"token":"[^"]*' | grep -o '[^"]*$')
    
    if [[ ! -z "$TOKEN" ]]; then
        echo -e "${GREEN}✓ 用户登录成功${NC}"
        echo "Token: ${TOKEN:0:50}..."
        export AUTH_TOKEN="$TOKEN"
    else
        echo -e "${RED}✗ 用户登录失败${NC}"
        export AUTH_TOKEN=""
    fi
}

# 测试获取用户信息
test_user_info() {
    echo -e "\n${YELLOW}[4/6] 测试获取用户信息...${NC}"
    
    if [[ -z "$AUTH_TOKEN" ]]; then
        echo -e "${RED}✗ 未登录，跳过测试${NC}"
        return
    fi
    
    response=$(curl -s -X GET "${API_URL}/auth/me" \
        -H "Authorization: Bearer ${AUTH_TOKEN}")
    
    echo "响应: $response"
    
    if [[ $response == *"username"* ]]; then
        echo -e "${GREEN}✓ 获取用户信息成功${NC}"
    else
        echo -e "${RED}✗ 获取用户信息失败${NC}"
    fi
}

# 测试商品列表
test_products() {
    echo -e "\n${YELLOW}[5/6] 测试获取商品列表...${NC}"
    
    response=$(curl -s -X GET "${API_URL}/products?page=1&page_size=10")
    
    echo "响应: ${response:0:200}..."
    
    if [[ $response == *"list"* ]] || [[ $response == *"total"* ]]; then
        echo -e "${GREEN}✓ 获取商品列表成功${NC}"
    else
        echo -e "${YELLOW}⚠ 商品列表为空或格式异常${NC}"
    fi
}

# 测试 WebSocket
test_websocket() {
    echo -e "\n${YELLOW}[6/6] 测试 WebSocket 连接...${NC}"
    
    if [[ -z "$AUTH_TOKEN" ]]; then
        echo -e "${YELLOW}⚠ 未登录，跳过 WebSocket 测试${NC}"
        return
    fi
    
    # 检查 wscat 是否安装
    if ! command -v wscat &> /dev/null; then
        echo -e "${YELLOW}⚠ wscat 未安装，跳过 WebSocket 测试${NC}"
        echo "  安装方法: npm install -g wscat"
        return
    fi
    
    echo "连接 WebSocket: ws://localhost:8080/ws"
    echo "（3秒后自动断开）"
    
    timeout 3 wscat -c "ws://localhost:8080/ws" \
        -H "Authorization: Bearer ${AUTH_TOKEN}" 2>/dev/null || true
    
    echo -e "${GREEN}✓ WebSocket 连接测试完成${NC}"
}

# 主函数
main() {
    # 检查服务是否启动
    if ! curl -s "${BASE_URL}/health" > /dev/null 2>&1; then
        echo -e "${RED}错误: 服务未启动，请先运行 'make run' 或 'docker-compose up'${NC}"
        exit 1
    fi
    
    test_health
    test_register
    test_login
    test_user_info
    test_products
    test_websocket
    
    echo -e "\n================================================"
    echo -e "${GREEN}API 测试完成！${NC}"
    echo "================================================"
    echo -e "\n更多测试命令："
    echo "  查看 API 文档: http://localhost:8080/swagger/index.html"
    echo "  查看健康状态: curl http://localhost:8080/health"
    echo "  测试商品搜索: curl http://localhost:8080/api/v1/products/search?keyword=测试"
}

main
