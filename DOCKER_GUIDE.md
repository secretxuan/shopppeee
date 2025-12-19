# Docker 运行指南

## 🎯 核心问题和解决方案

### 问题1: `docker-compose: command not found`

**原因：** 你的 Docker 版本太旧（1.13.1），不包含 docker-compose

**解决：** 运行安装脚本
```bash
./setup-docker.sh
```

或手动安装：
```bash
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" \
  -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version
```

---

### 问题2: `Cannot connect to Docker daemon`

**原因：** Docker 服务未运行

**解决方案（按优先级）：**

#### 方案A: 启动 Docker 服务（标准环境）
```bash
# CentOS/RHEL
sudo systemctl start docker
sudo systemctl enable docker

# 验证
docker ps
```

#### 方案B: 在容器内运行（Docker-in-Docker）
如果你在容器内，需要特殊配置：

```bash
# 1. 启动容器时挂载 Docker socket
docker run -v /var/run/docker.sock:/var/run/docker.sock ...

# 2. 或使用特权模式
docker run --privileged ...

# 3. 或使用 DinD（Docker in Docker）
docker run --privileged -d docker:dind
```

#### 方案C: 使用主机 Docker（推荐）
**如果你在云 IDE 或远程开发环境：**

你的环境可能本身就不支持嵌套 Docker。最简单的方案是：

1. **在本地机器上运行 Docker**
   ```bash
   # 在本地机器（不是云 IDE）
   git clone <your-repo>
   cd shoppee
   docker-compose up -d
   ```

2. **或使用云数据库，本地运行代码**
   ```bash
   # 在当前环境（云 IDE）
   ./run-local.sh
   ```

---

## ✅ 完整启动流程

### 1. 检查环境
```bash
# 检查 Docker
docker --version  # 期望: 任何版本

# 检查 docker-compose
docker-compose --version  # 期望: 1.29.2+

# 检查 Docker daemon
docker ps  # 期望: 无错误
```

### 2. 如果环境正常，启动服务
```bash
cd /data/workspace/shoppee
docker-compose up -d
```

### 3. 验证服务
```bash
# 查看容器状态
docker-compose ps

# 查看日志
docker-compose logs -f app

# 测试 API
curl http://localhost:8080/health
```

---

## 🔧 环境限制说明

### 你当前的环境：
- **OS**: Tencent tlinux 2.6
- **Docker**: 1.13.1 (非常旧)
- **环境类型**: 可能是容器内或云 IDE

### 推荐方案：

#### ✅ 如果只是开发测试
**使用本地运行模式**（无需 Docker）:
```bash
./run-local.sh
# 配合云数据库（Supabase + Upstash）
```
→ 详见 [LOCAL_SETUP.md](LOCAL_SETUP.md)

#### ✅ 如果需要完整 Docker 环境
**在本地机器运行**:
```bash
# 在你的笔记本/台式机
git clone <repo>
docker-compose up -d
```

#### ⚠️ 如果必须在当前环境用 Docker
需要联系系统管理员：
- 升级 Docker 到 20.10+
- 或提供 Docker socket 访问权限
- 或启用特权模式

---

## 📊 服务架构

docker-compose 会启动 3 个服务：

```
┌─────────────────────────────────────┐
│  App (Go 应用)                       │
│  Port: 8080                          │
└─────────────┬───────────────────────┘
              │
    ┌─────────┴─────────┐
    │                   │
┌───▼──────┐      ┌─────▼─────┐
│PostgreSQL│      │   Redis   │
│Port: 5432│      │Port: 6379 │
└──────────┘      └───────────┘
```

---

## 🎯 快速决策树

```
能运行 docker ps 吗？
│
├─ 是 → 运行 ./setup-docker.sh → docker-compose up -d ✅
│
└─ 否 → 选择：
       │
       ├─ 能安装/启动 Docker？
       │  └─ sudo systemctl start docker
       │
       ├─ 在云 IDE 容器内？
       │  └─ 使用本地运行: ./run-local.sh
       │
       └─ 有本地机器？
          └─ 本地运行 docker-compose
```

---

## 需要帮助？

告诉我你的情况：
1. `docker ps` 的输出结果
2. 你的使用场景（开发/测试/生产）
3. 是否可以访问主机 Docker

我会给你最合适的解决方案！
