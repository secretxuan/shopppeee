# 本地运行指南（无 Docker）

由于你的环境无法运行 Docker，以下是本地直接运行的方案。

## 方案选择

### 🌟 方案 1：使用在线数据库（最简单，推荐）

使用免费的云数据库服务，无需本地安装：

#### PostgreSQL 免费服务
- **Supabase** (推荐): https://supabase.com - 500MB 免费
- **ElephantSQL**: https://www.elephantsql.com - 20MB 免费
- **Neon**: https://neon.tech - 无限免费层

#### Redis 免费服务
- **Upstash**: https://upstash.com - 10,000 命令/天免费
- **Redis Cloud**: https://redis.com/try-free - 30MB 免费

#### 配置步骤
1. 注册并创建数据库实例
2. 获取连接信息
3. 修改 `.env.local` 文件：
```bash
# PostgreSQL (替换为你的连接信息)
DB_HOST=your-db-host.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres

# Redis (替换为你的连接信息)
REDIS_HOST=your-redis.upstash.io
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
```

4. 运行应用：
```bash
./run-local.sh
```

---

### 🔧 方案 2：本地安装 PostgreSQL 和 Redis

如果你有 sudo 权限，可以安装到本地：

#### 安装 PostgreSQL
```bash
# CentOS/RHEL
sudo yum install -y postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 创建数据库
sudo -u postgres createdb shoppee
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';"
```

#### 安装 Redis
```bash
# CentOS/RHEL
sudo yum install -y redis
sudo systemctl start redis
sudo systemctl enable redis
```

#### 运行应用
```bash
./run-local.sh
```

---

### 🎯 方案 3：仅运行代码（Mock 数据）

如果只是想看代码运行，可以修改为 SQLite + 内存模式：

```bash
# 安装 SQLite driver
go get gorm.io/driver/sqlite

# 运行（我可以帮你修改代码支持 SQLite）
go run cmd/api/main.go
```

---

## 快速测试（推荐 Supabase + Upstash）

### 1. 创建 Supabase 数据库
```bash
# 访问 https://supabase.com/dashboard
# 1. 注册并登录
# 2. 创建新项目
# 3. 复制连接信息
```

### 2. 创建 Upstash Redis
```bash
# 访问 https://console.upstash.com
# 1. 注册并登录
# 2. 创建 Redis 数据库
# 3. 复制连接信息
```

### 3. 修改配置
编辑 `.env.local`，填入你的连接信息

### 4. 运行
```bash
./run-local.sh
```

---

## 验证运行

应用启动后，访问：
- **API 文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/api/v1/health

---

## 需要帮助？

告诉我你选择哪个方案，我可以：
1. 帮你配置云数据库连接
2. 修改代码支持 SQLite
3. 提供更详细的安装步骤
