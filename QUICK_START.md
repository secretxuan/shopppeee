# 🚀 快速开始指南

这是最快速的启动方式，5 分钟内即可运行！

## 方式一：Docker Compose（最简单）

### 1. 设置 Docker 环境

```bash
# 检查 Docker 版本
docker --version

# 运行设置脚本（会自动安装 docker-compose 并检查环境）
./setup-docker.sh
```

**常见问题：**
- ❌ `docker-compose: command not found` → 运行 `./setup-docker.sh` 自动安装
- ❌ `Cannot connect to Docker daemon` → Docker 服务未启动，见下方解决方案

**如果 Docker daemon 未运行：**
```bash
# 方法1: 启动 Docker 服务（需要特权）
sudo systemctl start docker

# 方法2: 如果在容器内，需要挂载 Docker socket
# docker run -v /var/run/docker.sock:/var/run/docker.sock ...

# 方法3: 使用本地运行模式
./run-local.sh  # 见方式二
```

### 2. 一键启动

```bash
# 启动所有服务（PostgreSQL + Redis + App）
docker-compose up -d

# 查看日志
docker-compose logs -f app
```

### 3. 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**完成！** 应用已在 http://localhost:8080 运行。

---

## 方式二：本地开发（Go 环境）

### 1. 前置要求

- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### 2. 启动数据库

```bash
# 使用 Docker 启动数据库
docker-compose up -d postgres redis

# 或手动安装并启动
```

### 3. 配置环境变量

```bash
cp .env.example .env
# 根据需要修改 .env
```

### 4. 启动应用

```bash
# 下载依赖
go mod download

# 运行
go run cmd/api/main.go

# 或使用 Makefile
make run
```

---

## 方式三：使用快速启动脚本

```bash
# 开发模式（本地 Go + Docker 数据库）
./scripts/start.sh dev

# 生产模式（全部 Docker）
./scripts/start.sh prod
```

---

## 测试 API

### 自动化测试脚本

```bash
./scripts/test_api.sh
```

### 手动测试

#### 1. 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "password123",
    "phone": "13800138000"
  }'
```

#### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123"
  }'
```

响应示例：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": 1702886400,
    "user": {
      "id": 1,
      "username": "john",
      "email": "john@example.com",
      "role": "user"
    }
  }
}
```

#### 3. 获取用户信息（需要 Token）

```bash
TOKEN="your_jwt_token_here"

curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

#### 4. 获取商品列表

```bash
# 基础列表
curl http://localhost:8080/api/v1/products

# 分页 + 筛选
curl "http://localhost:8080/api/v1/products?page=1&page_size=20&category_id=1&sort=price_asc"

# 搜索商品
curl "http://localhost:8080/api/v1/products/search?keyword=手机"
```

---

## WebSocket 测试

### 使用 wscat

```bash
# 安装 wscat
npm install -g wscat

# 连接 WebSocket（需要 Token）
wscat -c "ws://localhost:8080/ws" \
  -H "Authorization: Bearer $TOKEN"

# 发送消息
> {"type": "ping"}

# 接收响应
< {"type": "pong", "content": "ok", "time": 1702886400}
```

### 使用浏览器

```javascript
// 在浏览器控制台运行
const token = 'your_jwt_token_here';
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

ws.onopen = () => {
  console.log('WebSocket 已连接');
  ws.send(JSON.stringify({ type: 'ping' }));
};

ws.onmessage = (event) => {
  console.log('收到消息:', JSON.parse(event.data));
};
```

---

## 常用命令

### Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down

# 重启服务
docker-compose restart app
```

### Makefile

```bash
# 编译
make build

# 运行
make run

# 测试
make test

# 代码检查
make lint

# 格式化
make fmt

# Docker 操作
make docker-build
make docker-up
make docker-down
```

---

## 停止服务

```bash
# Docker Compose
docker-compose down

# 本地运行（按 Ctrl+C）
```

---

## 问题排查

### 端口被占用

```bash
# 查看端口占用
lsof -i :8080
lsof -i :5432
lsof -i :6379

# 修改 .env 中的端口配置
APP_PORT=8081
DB_PORT=5433
REDIS_PORT=6380
```

### 数据库连接失败

```bash
# 检查 PostgreSQL 状态
docker-compose ps postgres

# 查看数据库日志
docker-compose logs postgres

# 手动连接测试
psql -h localhost -U postgres -d shoppee
```

### Redis 连接失败

```bash
# 检查 Redis 状态
docker-compose ps redis

# 测试连接
redis-cli ping
```

---

## 下一步

- 📖 阅读 [README.md](README.md) 了解详细功能
- 🗄️ 查看 [DATABASE_DESIGN.md](DATABASE_DESIGN.md) 了解数据库设计
- 🚀 阅读 [DEPLOYMENT.md](DEPLOYMENT.md) 了解部署方案
- ⚡ 查看 [PERFORMANCE.md](PERFORMANCE.md) 学习性能优化

---

## 需要帮助？

- GitHub Issues: https://github.com/yourusername/shoppee/issues
- 项目文档: [README.md](README.md)

祝你使用愉快！🎉
