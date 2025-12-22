# 🚀 没有 Node.js 的启动方案

如果你的系统上没有安装 Node.js，有以下几种解决方案：

## 方案一：安装 Node.js（推荐用于开发）

### CentOS/RHEL 系统

```bash
# 安装 Node.js 18.x LTS
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install -y nodejs

# 验证安装
node --version
npm --version
```

### 使用 dnf（较新的系统）

```bash
# 安装 Node.js
sudo dnf module install nodejs:18

# 验证
node --version
npm --version
```

安装完成后：
```bash
cd frontend
npm install
npm run dev
```

---

## 方案二：使用 Docker 运行前端（推荐用于生产）

**无需安装 Node.js，直接使用 Docker！**

### 启动完整服务（后端 + 前端）

```bash
# 使用包含前端的 Docker Compose 配置
docker compose -f docker-compose.frontend.yml up -d

# 或者使用 profile 方式
docker compose --profile with-frontend up -d
```

### 访问地址

- 🎨 **前端**: http://localhost:3000
- 🔧 **后端 API**: http://localhost:8080

### 查看日志

```bash
# 查看所有服务
docker compose -f docker-compose.frontend.yml logs -f

# 只看前端
docker compose -f docker-compose.frontend.yml logs -f frontend

# 只看后端
docker compose -f docker-compose.frontend.yml logs -f app
```

### 停止服务

```bash
docker compose -f docker-compose.frontend.yml down
```

---

## 方案三：只使用后端，通过 API 测试

如果暂时不需要前端界面，可以直接使用后端 API：

```bash
# 启动后端（已完成）
docker compose up -d

# 测试 API
curl http://localhost:8080/health

# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 获取商品列表
curl http://localhost:8080/api/v1/products
```

---

## 推荐方案对比

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| **方案一: 安装 Node.js** | • 开发体验好<br>• 热重载快<br>• 调试方便 | • 需要安装 Node.js<br>• 占用系统资源 | 本地开发 |
| **方案二: Docker** | • 无需安装 Node.js<br>• 环境一致<br>• 部署简单 | • 构建时间较长<br>• 修改需重新构建 | 生产部署<br>演示环境 |
| **方案三: 只用后端** | • 最简单<br>• 资源占用少 | • 无用户界面<br>• 需要手动测试 API | API 开发<br>后端测试 |

---

## 🎯 快速决策

### 如果你想...

**🎨 看到完整的前端界面**
→ 选择**方案二**（Docker）最简单！

```bash
docker compose -f docker-compose.frontend.yml up -d
```

**🔧 进行前端开发**
→ 选择**方案一**（安装 Node.js）

```bash
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install -y nodejs
cd frontend && npm install && npm run dev
```

**⚡ 只测试后端 API**
→ 选择**方案三**（已完成）

```bash
# 后端已启动，直接测试
curl http://localhost:8080/api/v1/products
```

---

## 💡 我的建议

### 快速体验（推荐）

使用 Docker 启动完整服务：

```bash
# 1. 启动完整服务（后端 + 前端）
docker compose -f docker-compose.frontend.yml up -d

# 2. 等待构建完成（首次需要几分钟）
docker compose -f docker-compose.frontend.yml logs -f frontend

# 3. 访问前端
# http://localhost:3000
```

### 长期开发

安装 Node.js 进行开发：

```bash
# 1. 安装 Node.js
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install -y nodejs

# 2. 启动开发服务器
cd frontend
npm install
npm run dev

# 3. 访问
# http://localhost:3000
```

---

## 🐛 常见问题

### Q: Docker 构建前端很慢？

**A**: 首次构建需要下载依赖，可能需要 5-10 分钟。后续会使用缓存，很快。

### Q: 想修改前端代码怎么办？

**A**: 
- 如果用 Docker：修改后需要重新构建
  ```bash
  docker compose -f docker-compose.frontend.yml up -d --build frontend
  ```
- 如果用 Node.js：自动热重载，立即生效

### Q: 可以同时启动吗？

**A**: 不建议。选择一种方式即可：
- Docker 方式：前端在容器内，访问 http://localhost:3000
- Node.js 方式：前端在开发服务器，访问 http://localhost:3000

---

## ✅ 总结

**最简单的方式**：
```bash
docker compose -f docker-compose.frontend.yml up -d
```

访问 http://localhost:3000 即可看到完整的前端界面！

**开发推荐方式**：
```bash
# 安装 Node.js
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install -y nodejs

# 启动前端
cd frontend && npm install && npm run dev
```

选择适合你的方式开始吧！🚀
