# 🚀 Shoppee 电商系统 - 快速启动指南

## ⚡ 一键启动（最快）

```bash
cd /data/workspace/shopppeee

# 启动后端
sudo docker compose up -d

# 启动前端（新终端）
cd frontend
npm install
npm run dev
```

**访问地址**：
- 🎨 前端：http://localhost:3000
- 🔧 后端：http://localhost:8080
- ❤️ 健康检查：http://localhost:8080/health

---

## 📝 完整启动流程

### 第一步：启动后端服务

```bash
cd /data/workspace/shopppeee

# 启动数据库和后端
sudo docker compose up -d

# 查看启动日志
sudo docker compose logs -f app
```

**等待日志显示**：
```
INFO    api/main.go:82  服务器启动      {"port": 8080}
INFO    database/migrate.go:31  数据库迁移完成
```

### 第二步：启动前端（新终端窗口）

```bash
cd /data/workspace/shopppeee/frontend

# 安装依赖（首次运行）
npm install

# 启动开发服务器
npm run dev
```

**等待显示**：
```
➜  Local:   http://localhost:3000/
```

---

## 🎯 快速测试流程

### 1. 注册账号

浏览器访问：http://localhost:3000/register

填写信息：
- 用户名：admin
- 邮箱：admin@example.com
- 密码：password123
- 手机：13800138000

### 2. 创建管理员账号（可选）

**方式一：修改数据库**
```bash
# 连接数据库
sudo docker exec -it shoppee-postgres psql -U postgres -d shoppee

# 将用户设为管理员
UPDATE users SET role = 'admin' WHERE username = 'admin';

# 退出
\q
```

**方式二：注册时使用特殊用户名**
- 用户名包含 "admin" 的会自动设为管理员

### 3. 创建商品分类

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "电子产品",
    "description": "手机、电脑等电子产品",
    "sort": 1
  }'
```

### 4. 上架商品

**方式一：使用管理后台（推荐）**
- 访问：http://localhost:3000/admin/products
- 点击"添加商品"
- 填写商品信息并保存

**方式二：使用API**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "最新款苹果手机，性能强劲",
    "price": 7999.00,
    "orig_price": 8999.00,
    "stock": 50,
    "sku": "IPHONE15PRO-001",
    "category_id": 1,
    "status": "active"
  }'
```

**方式三：直接操作数据库**
```sql
-- 连接数据库
sudo docker exec -it shoppee-postgres psql -U postgres -d shoppee

-- 插入分类
INSERT INTO categories (name, description, sort, status, created_at, updated_at)
VALUES ('电子产品', '手机、电脑等', 1, 'active', NOW(), NOW());

-- 插入商品
INSERT INTO products (name, description, price, orig_price, stock, sku, category_id, status, created_at, updated_at)
VALUES 
  ('iPhone 15 Pro', '最新款苹果手机', 7999.00, 8999.00, 50, 'IPHONE15PRO-001', 1, 'active', NOW(), NOW()),
  ('MacBook Pro', '苹果笔记本电脑', 12999.00, 14999.00, 30, 'MACBOOK-001', 1, 'active', NOW(), NOW()),
  ('AirPods Pro', '苹果无线耳机', 1599.00, 1999.00, 100, 'AIRPODS-001', 1, 'active', NOW(), NOW());
```

### 5. 完整购物流程测试

1. **登录账号** → http://localhost:3000/login
2. **浏览商品** → http://localhost:3000/products
3. **查看详情** → 点击任意商品
4. **加入购物车** → 点击"加入购物车"按钮
5. **查看购物车** → 点击顶部购物车图标
6. **创建订单** → 点击"去结算"（需要先创建收货地址）
7. **支付订单** → 选择支付方式
8. **查看订单** → http://localhost:3000/orders

---

## 🔧 常见问题

### Q1: 后端启动失败？

**检查端口占用**：
```bash
sudo lsof -i:8080
sudo lsof -i:5432
sudo lsof -i:6379
```

**查看日志**：
```bash
sudo docker compose logs app
```

### Q2: 前端连接不上后端？

**检查 CORS 配置**：
确保后端启动成功，并且 CORS 中间件正常工作。

**检查环境变量**：
```bash
cat frontend/.env
```

应该包含：
```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### Q3: 数据库连接失败？

```bash
# 查看PostgreSQL日志
sudo docker compose logs postgres

# 重启数据库
sudo docker compose restart postgres
```

### Q4: npm install 失败？

```bash
# 清除缓存
npm cache clean --force

# 删除 node_modules
rm -rf node_modules package-lock.json

# 重新安装
npm install
```

---

## 📦 数据库初始化

### 创建测试数据

```sql
-- 连接数据库
sudo docker exec -it shoppee-postgres psql -U postgres -d shoppee

-- 创建分类
INSERT INTO categories (name, description, sort, status, created_at, updated_at) VALUES
  ('电子产品', '手机、电脑等电子产品', 1, 'active', NOW(), NOW()),
  ('服装鞋包', '男装、女装、鞋子、包包', 2, 'active', NOW(), NOW()),
  ('食品饮料', '零食、饮料、生鲜', 3, 'active', NOW(), NOW()),
  ('家居生活', '家具、家纺、日用品', 4, 'active', NOW(), NOW());

-- 创建商品（电子产品）
INSERT INTO products (name, description, price, orig_price, stock, sku, category_id, status, created_at, updated_at) VALUES
  ('iPhone 15 Pro', '苹果最新旗舰手机，A17仿生芯片', 7999.00, 8999.00, 50, 'IPHONE15PRO-001', 1, 'active', NOW(), NOW()),
  ('MacBook Pro 14', '苹果笔记本电脑，M3芯片', 12999.00, 14999.00, 30, 'MACBOOK14-001', 1, 'active', NOW(), NOW()),
  ('AirPods Pro 2', '苹果无线降噪耳机', 1599.00, 1999.00, 100, 'AIRPODS2-001', 1, 'active', NOW(), NOW()),
  ('iPad Air', '10.9英寸平板电脑', 4599.00, 4999.00, 60, 'IPADAIR-001', 1, 'active', NOW(), NOW()),
  ('Apple Watch', '智能手表，健康监测', 2999.00, 3299.00, 80, 'WATCH-001', 1, 'active', NOW(), NOW());

-- 创建商品（服装）
INSERT INTO products (name, description, price, orig_price, stock, sku, category_id, status, created_at, updated_at) VALUES
  ('男士T恤', '纯棉舒适，多色可选', 99.00, 159.00, 200, 'TSHIRT-M-001', 2, 'active', NOW(), NOW()),
  ('女士连衣裙', '优雅时尚，适合春夏', 299.00, 499.00, 150, 'DRESS-W-001', 2, 'active', NOW(), NOW()),
  ('运动鞋', '透气舒适，适合跑步', 399.00, 599.00, 120, 'SHOES-001', 2, 'active', NOW(), NOW()),
  ('双肩包', '大容量，多功能口袋', 199.00, 299.00, 100, 'BAG-001', 2, 'active', NOW(), NOW());

-- 创建管理员用户
INSERT INTO users (username, email, password, phone, role, status, created_at, updated_at) VALUES
  ('admin', 'admin@example.com', '$2a$10$xxxxx', '13800000000', 'admin', 'active', NOW(), NOW());
```

---

## 🎯 测试API

### 获取商品列表
```bash
curl http://localhost:8080/api/v1/products
```

### 搜索商品
```bash
curl "http://localhost:8080/api/v1/products/search?keyword=iPhone"
```

### 健康检查
```bash
curl http://localhost:8080/health
```

---

## 🛑 停止服务

```bash
# 停止前端（在前端终端按 Ctrl+C）

# 停止后端
cd /data/workspace/shopppeee
sudo docker compose down

# 完全清理（包括数据）
sudo docker compose down -v
```

---

## 🎉 恭喜！

你已经成功启动了 Shoppee 电商系统！

**现在可以**：
- ✅ 浏览商品
- ✅ 搜索商品
- ✅ 加入购物车
- ✅ 创建订单
- ✅ 管理商品（管理员）
- ✅ 管理订单（管理员）

**下一步**：
- 阅读 `COMPLETION_REPORT.md` 了解完整功能
- 查看 API 文档了解所有接口
- 开始开发你自己的功能！

祝开发愉快！🚀
