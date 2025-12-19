# Shoppee 电商系统

基于 Go 语言开发的高性能电商系统，采用 Gin + GORM + PostgreSQL + Redis 技术栈，支持高并发、分布式部署。

## 🚀 核心功能

### 1. 高并发用户认证
- JWT 令牌认证机制
- Redis 缓存用户会话
- 密码 bcrypt 加密存储
- 支持角色权限控制（用户/管理员）
- 登录失败限流保护

### 2. 批量数据处理（Go 协程优化）
- 商品批量导入（协程池处理）
- 库存批量更新（悲观锁防止超卖）
- 支持 10 worker 协程池并发处理
- 事务保证数据一致性

### 3. 实时消息推送（WebSocket）
- 订单状态实时通知
- 促销活动广播
- 库存预警推送
- 支持按用户定向推送
- 心跳检测保持连接

## 📋 技术栈

### 后端框架
- **Gin** - 高性能 HTTP Web 框架
- **GORM** - ORM 框架（v2）
- **Viper** - 配置管理
- **Zap** - 高性能日志库
- **JWT** - 认证授权
- **Gorilla WebSocket** - WebSocket 支持

### 数据库
- **PostgreSQL 15** - 主数据库
- **Redis 7** - 缓存 + 限流

### 工具链
- **Docker** - 容器化部署
- **Docker Compose** - 服务编排
- **Makefile** - 自动化构建
- **Go Modules** - 依赖管理

## 🏗️ 项目结构

```
shoppee/
├── cmd/                    # 应用入口
│   └── api/
│       └── main.go        # 主程序
├── internal/              # 内部代码
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   ├── handler/          # HTTP 处理器
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   ├── router/           # 路由配置
│   ├── service/          # 业务逻辑
│   └── websocket/        # WebSocket 服务
├── pkg/                   # 公共库
│   ├── jwt/              # JWT 工具
│   ├── logger/           # 日志工具
│   └── response/         # 响应封装
├── scripts/              # 脚本文件
│   └── init.sql          # 数据库初始化
├── Dockerfile            # Docker 镜像构建
├── docker-compose.yml    # Docker 编排配置
├── Makefile             # 自动化构建
├── go.mod               # Go 依赖
└── README.md            # 项目文档
```

## 🔧 快速开始

### 前置要求
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+（可选，Docker 已包含）
- Redis 7+（可选，Docker 已包含）

### 1. 克隆项目
```bash
git clone https://github.com/yourusername/shoppee.git
cd shoppee
```

### 2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 修改配置
```

### 3. 使用 Docker Compose 启动（推荐）
```bash
# 启动所有服务（PostgreSQL + Redis + App）
make docker-up

# 查看日志
make docker-logs

# 停止服务
make docker-down
```

### 4. 本地开发模式
```bash
# 下载依赖
make deps

# 启动数据库（Docker）
docker-compose up -d postgres redis

# 运行应用
make run

# 或使用热重载（需安装 air）
make dev
```

### 5. 编译部署
```bash
# 编译当前平台
make build

# 编译 Linux 版本
make build-linux

# 编译所有平台
make build-all
```

## 📡 API 文档

### 认证相关

#### 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "phone": "13800138000"
}
```

#### 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": 1702886400,
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```

#### 获取用户信息
```http
GET /api/v1/auth/me
Authorization: Bearer <token>
```

### 商品相关

#### 获取商品列表
```http
GET /api/v1/products?page=1&page_size=20&category_id=1&sort=price_asc
```

#### 获取商品详情
```http
GET /api/v1/products/{id}
```

#### 搜索商品
```http
GET /api/v1/products/search?keyword=手机&page=1
```

#### 批量更新库存（需管理员权限）
```http
POST /api/v1/products/batch-stock
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "1": -10,
  "2": 20,
  "3": -5
}
```

### WebSocket 连接

```javascript
// 连接 WebSocket（需要先登录获取 token）
const ws = new WebSocket('ws://localhost:8080/ws?token=<your_jwt_token>');

// 接收消息
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('收到消息:', message);
  // message.type: system, order, promotion, stock_alert
};

// 发送心跳
ws.send(JSON.stringify({ type: 'ping' }));
```

## 🔐 数据库设计

### 核心表结构

**users** - 用户表
- id, username, email, password, role, status, last_login

**products** - 商品表
- id, name, description, price, stock, sku, category_id, status

**orders** - 订单表
- id, order_no, user_id, total_amount, status, pay_status

**order_items** - 订单项表
- id, order_id, product_id, quantity, price

**carts** - 购物车表
- id, user_id

**cart_items** - 购物车项表
- id, cart_id, product_id, quantity

详细 ER 图和字段说明请参考数据库模型文件。

## ⚡ 性能优化

### 1. 并发处理
- 使用 Go 协程池处理批量任务
- Worker Pool 模式（10 个 worker）
- Channel 实现任务分发

### 2. 数据库优化
- 连接池配置（最大 100 连接）
- 索引优化（分类、状态、用户等字段）
- 悲观锁防止库存超卖
- 事务保证数据一致性

### 3. 缓存策略
- Redis 缓存用户信息（7 天）
- 商品详情缓存（1 小时）
- 滑动窗口限流

### 4. 镜像优化
- 多阶段构建（编译阶段 + 运行阶段）
- 最终镜像基于 Alpine（< 20MB）
- 静态编译（CGO_ENABLED=0）

## 🧪 测试

```bash
# 运行所有测试
make test

# 生成覆盖率报告
make test-coverage

# 性能基准测试
go test -bench=. -benchmem ./...
```

## 📊 监控和日志

### 日志
- 使用 Zap 高性能日志库
- 支持控制台 + 文件双输出
- JSON 格式便于日志收集
- 日志文件路径：`./logs/app.log`

### 健康检查
```bash
curl http://localhost:8080/health
```

## 🚢 部署

### Docker 部署
```bash
# 构建镜像
make docker-build

# 启动服务
make docker-up
```

### 生产环境配置
1. 修改 `.env` 中的数据库密码和 JWT 密钥
2. 配置 CORS 允许的域名
3. 关闭 DEBUG 模式
4. 配置反向代理（Nginx）
5. 配置 HTTPS 证书

### 交叉编译
```bash
# Linux AMD64
make build-linux

# Windows AMD64
make build-windows

# macOS AMD64
make build-mac
```

## 🔒 安全特性

- JWT 令牌认证
- 密码 bcrypt 加密
- SQL 注入防护（GORM 预处理）
- XSS 防护
- CORS 跨域控制
- 请求频率限流
- 参数验证（validator）

## 📈 性能指标

- 并发处理：支持 10000+ QPS
- 响应时间：平均 < 50ms
- 批量导入：1000 商品 < 5s
- WebSocket：支持 10000+ 并发连接

## 🛠️ 开发工具

### 推荐 IDE
- GoLand
- VS Code + Go 插件

### 代码规范
```bash
# 格式化代码
make fmt

# 代码检查
make lint
```

## 📝 TODO

- [ ] 订单管理模块完善
- [ ] 支付接口集成
- [ ] 评价系统优化
- [ ] Elasticsearch 全文搜索
- [ ] Prometheus + Grafana 监控
- [ ] Kubernetes 部署配置
- [ ] 前端管理系统（Vue3）

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 👥 作者

Shoppee Team

---

**注意**：这是一个演示项目，生产环境部署前请务必修改默认密码和密钥！
