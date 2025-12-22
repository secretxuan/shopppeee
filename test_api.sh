#!/bin/bash

# Shoppee 电商系统 API 测试脚本

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

echo "================================"
echo "Shoppee 电商系统 API 测试"
echo "================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local auth=$5

    echo -e "${BLUE}测试: ${name}${NC}"
    
    if [ -z "$auth" ]; then
        response=$(curl -s -X ${method} "${BASE_URL}${endpoint}" \
            -H "Content-Type: application/json" \
            -d "${data}")
    else
        response=$(curl -s -X ${method} "${BASE_URL}${endpoint}" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer ${TOKEN}" \
            -d "${data}")
    fi
    
    echo "$response" | jq '.'
    
    # 检查响应
    if echo "$response" | jq -e '.code' > /dev/null 2>&1; then
        code=$(echo "$response" | jq -r '.code')
        if [ "$code" == "0" ] || [ "$code" == "200" ]; then
            echo -e "${GREEN}✓ 成功${NC}"
        else
            echo -e "${RED}✗ 失败${NC}"
        fi
    fi
    
    echo ""
}

# 1. 健康检查
echo "================================"
echo "1. 健康检查"
echo "================================"
curl -s http://localhost:8080/health | jq '.'
echo ""

# 2. 用户注册
echo "================================"
echo "2. 用户注册"
echo "================================"
test_api "注册测试用户" "POST" "/auth/register" '{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "phone": "13800138000"
}'
echo ""

# 3. 用户登录
echo "================================"
echo "3. 用户登录"
echo "================================"
login_response=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
      "username": "testuser",
      "password": "password123"
    }')

echo "$login_response" | jq '.'

# 提取 token
TOKEN=$(echo "$login_response" | jq -r '.data.token')
echo -e "${GREEN}Token: ${TOKEN}${NC}"
echo ""

# 4. 获取商品列表
echo "================================"
echo "4. 获取商品列表"
echo "================================"
test_api "获取商品列表" "GET" "/products?page=1&page_size=5" ""
echo ""

# 5. 搜索商品
echo "================================"
echo "5. 搜索商品"
echo "================================"
test_api "搜索iPhone" "GET" "/products/search?keyword=iPhone" ""
echo ""

# 6. 获取分类列表
echo "================================"
echo "6. 获取分类列表"
echo "================================"
test_api "获取分类列表" "GET" "/categories" ""
echo ""

# 7. 添加到购物车（需要认证）
if [ -n "$TOKEN" ]; then
    echo "================================"
    echo "7. 添加到购物车"
    echo "================================"
    test_api "添加商品到购物车" "POST" "/cart/items" '{
      "product_id": 1,
      "quantity": 2
    }' "auth"
    echo ""

    # 8. 查看购物车
    echo "================================"
    echo "8. 查看购物车"
    echo "================================"
    test_api "查看购物车" "GET" "/cart" "" "auth"
    echo ""

    # 9. 获取用户信息
    echo "================================"
    echo "9. 获取用户信息"
    echo "================================"
    test_api "获取当前用户信息" "GET" "/auth/me" "" "auth"
    echo ""

    # 10. 创建收货地址
    echo "================================"
    echo "10. 创建收货地址"
    echo "================================"
    test_api "创建收货地址" "POST" "/addresses" '{
      "name": "张三",
      "phone": "13800138000",
      "province": "广东省",
      "city": "深圳市",
      "district": "南山区",
      "detail": "科技园南区xx路xx号",
      "is_default": true
    }' "auth"
    echo ""

    # 11. 获取地址列表
    echo "================================"
    echo "11. 获取地址列表"
    echo "================================"
    test_api "获取地址列表" "GET" "/addresses" "" "auth"
    echo ""
fi

echo "================================"
echo "测试完成！"
echo "================================"
echo ""
echo "提示："
echo "1. 如需测试管理员功能，请先将用户设为管理员"
echo "2. 可以使用 jq 工具美化 JSON 输出"
echo "3. 查看完整 API 文档：COMPLETION_REPORT.md"
echo ""
