# 🛍️ Shoppee 电商系统

<div align="center">

**现代化、全栈、企业级电商平台**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

</div>

---

## ✨ 项目亮点

- 🎯 **完整的电商业务流程** - 从商品浏览到订单完成的全链路
- 🔐 **企业级认证授权** - JWT + 角色权限管理
- 📦 **库存防超卖** - 悲观锁 + 事务保证数据一致性
- ⚡ **高性能缓存** - Redis 缓存热点数据
- 🎨 **现代化UI** - React + Ant Design 响应式设计
- 🐳 **容器化部署** - Docker Compose 一键启动
- 📊 **完善的后台管理** - 商品、订单、用户全方位管理

---

## 🏗️ 技术栈

### 后端
- **框架**: Gin (Go Web框架)
- **ORM**: GORM v2
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7
- **认证**: JWT
- **实时通信**: WebSocket
- **日志**: Zap
- **配置**: Viper

### 前端
- **框架**: React 18 + TypeScript
- **UI库**: Ant Design 5
- **路由**: React Router 6
- **状态管理**: Zustand
- **HTTP**: Axios
- **构建工具**: Vite

### 基础设施
- **容器化**: Docker + Docker Compose
- **代理**: Nginx (生产环境)

---

## 📋 核心功能

### 用户端功能
- ✅ 用户注册/登录
- ✅ 商品浏览/搜索
- ✅ 购物车管理
- ✅ 收货地址管理
- ✅ 下单支付
- ✅ 订单管理
- ✅ 商品评价
- ✅ 个人中心

### 管理端功能
- ✅ 商品管理（CRUD）
- ✅ 分类管理
- ✅ 订单管理
- ✅ 库存管理
- ✅ 评价管理
- ✅ 用户管理

---

## 🚀 快速开始

### 前置要求
- Docker & Docker Compose
- Node.js 18+ (前端开发)
- Go 1.21+ (后端开发，可选)

### 一键启动

```bash
# 克隆项目
cd /data/workspace/shopppeee

# 启动后端
sudo docker compose up -d

# 启动前端（新终端）
cd frontend
npm install
npm run dev
```

**访问地址**:
- 🎨 前端: http://localhost:3000
- 🔧 后端API: http://localhost:8080
- ❤️ 健康检查: http://localhost:8080/health

### 初始化数据（可选）

```bash
# 导入测试数据
sudo docker exec -i shoppee-postgres psql -U postgres -d shoppee < init_data.sql
```

---

## 📖 文档

- 📘 [快速启动指南](QUICK_START.md) - 5分钟启动项目
- 📗 [完成报告](COMPLETION_REPORT.md) - 详细功能列表和API文档
- 📙 [前端指南](frontend/README.md) - 前端开发说明

---

## 🎯 如何上架商品

### 方法一：使用管理后台（推荐）

1. 登录管理员账号
2. 访问 http://localhost:3000/admin/products
3. 点击"添加商品"
4. 填写商品信息并保存

### 方法二：使用API

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "最新款苹果手机",
    "price": 7999.00,
    "stock": 50,
    "sku": "IPHONE15PRO-001",
    "category_id": 1,
    "status": "active"
  }'
```

### 方法三：直接操作数据库

```sql
INSERT INTO products (name, description, price, stock, sku, category_id, status, created_at, updated_at)
VALUES ('测试商品', '商品描述', 99.99, 100, 'TEST-001', 1, 'active', NOW(), NOW());
```

---

## 🧪 API 测试

```bash
# 运行API测试脚本
./test_api.sh
```

### 核心API端点

#### 认证
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/auth/me` - 获取当前用户

#### 商品
- `GET /api/v1/products` - 商品列表
- `GET /api/v1/products/:id` - 商品详情
- `GET /api/v1/products/search` - 搜索商品
- `POST /api/v1/products` - 创建商品（管理员）

#### 购物车
- `GET /api/v1/cart` - 获取购物车
- `POST /api/v1/cart/items` - 添加商品
- `PUT /api/v1/cart/items/:id` - 更新数量
- `DELETE /api/v1/cart/items/:id` - 删除商品

#### 订单
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders` - 订单列表
- `GET /api/v1/orders/:id` - 订单详情
- `POST /api/v1/orders/:id/cancel` - 取消订单

[查看完整API文档](COMPLETION_REPORT.md#-api-文档总览)

---

## 📁 项目结构

```
shopppeee/
├── cmd/                    # 应用入口
├── internal/              # 内部代码
│   ├── handler/          # HTTP处理器
│   ├── service/          # 业务逻辑
│   ├── models/           # 数据模型
│   ├── router/           # 路由配置
│   ├── middleware/       # 中间件
│   └── database/         # 数据库
├── frontend/             # 前端项目
│   ├── src/
│   │   ├── api/         # API接口
│   │   ├── pages/       # 页面组件
│   │   ├── components/  # 通用组件
│   │   └── store/       # 状态管理
│   └── package.json
├── docker-compose.yml    # Docker编排
├── Dockerfile           # 后端镜像
└── README.md
```

---

## 🛠️ 开发

### 后端开发

```bash
# 安装依赖
go mod download

# 运行开发服务器
go run cmd/api/main.go

# 构建
go build -o shoppee cmd/api/main.go
```

### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

---

## 🐳 Docker 部署

### 开发环境

```bash
# 启动所有服务
docker compose up -d

# 查看日志
docker compose logs -f

# 停止服务
docker compose down
```

### 生产环境

```bash
# 包含前端的完整部署
docker compose -f docker-compose.frontend.yml up -d
```

---

## 📊 数据库

### 数据模型

- `users` - 用户表
- `products` - 商品表
- `categories` - 分类表
- `carts` - 购物车表
- `cart_items` - 购物车项表
- `addresses` - 收货地址表
- `orders` - 订单表
- `order_items` - 订单项表
- `payments` - 支付表
- `reviews` - 评价表

### 数据库操作

```bash
# 连接数据库
docker exec -it shoppee-postgres psql -U postgres -d shoppee

# 查看表
\dt

# 查询商品
SELECT * FROM products LIMIT 10;

# 退出
\q
```

---

## 🎨 UI截图

### 首页
- 轮播图展示
- 热门商品推荐
- 分类导航

### 商品列表
- 搜索筛选
- 排序功能
- 分页展示

### 购物车
- 商品管理
- 数量调整
- 实时总价

### 管理后台
- 商品管理
- 订单处理
- 数据统计

---

## 🔒 安全

- ✅ JWT 认证
- ✅ 密码加密（bcrypt）
- ✅ CORS 配置
- ✅ SQL 注入防护（GORM）
- ✅ XSS 防护
- ✅ 请求限流

---

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

---

## 📄 许可证

MIT License

---

## 📞 联系方式

- 项目地址: https://github.com/yourusername/shopppeee
- 问题反馈: Issues
- 邮箱: your@email.com

---

<div align="center">

**感谢使用 Shoppee 电商系统！**

⭐ 如果这个项目对你有帮助，请给个 Star！

Made with ❤️ by [Your Name]

</div>
