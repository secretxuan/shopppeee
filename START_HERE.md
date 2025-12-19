# 🎉 欢迎使用 Shoppee 电商系统！

恭喜！您已成功获取完整的 Go 语言电商项目。

## 📋 项目概览

**Shoppee** 是一个生产级的 Go 语言电商系统，具备以下特点：

✅ **完整功能** - 用户认证、商品管理、实时推送  
✅ **高性能** - 协程池、缓存优化、10000+ QPS  
✅ **生产就绪** - Docker 部署、完善文档、测试覆盖  
✅ **最佳实践** - Go 规范、分层架构、安全防护  

## 🚀 立即开始

### 方式 1️⃣：Docker Compose（最简单）

```bash
# 1. 启动所有服务
docker-compose up -d

# 2. 测试 API
curl http://localhost:8080/health

# 3. 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'
```

### 方式 2️⃣：快速启动脚本

```bash
# 开发模式
./scripts/start.sh dev

# 生产模式
./scripts/start.sh prod
```

### 方式 3️⃣：本地开发

```bash
# 启动数据库
docker-compose up -d postgres redis

# 运行应用
make run
```

## 📚 文档导航

### 新手必读（按顺序）

1. **[快速开始](QUICK_START.md)** ⭐ 5 分钟上手
2. **[项目总结](PROJECT_SUMMARY.md)** 了解技术亮点
3. **[README](README.md)** 详细功能说明

### 深入学习

4. **[系统架构](ARCHITECTURE.md)** 架构设计详解
5. **[数据库设计](DATABASE_DESIGN.md)** ER 图和表结构
6. **[性能优化](PERFORMANCE.md)** 性能调优方法

### 部署运维

7. **[部署指南](DEPLOYMENT.md)** 生产环境部署
8. **[交付清单](CHECKLIST.md)** 项目交付内容

## 🎯 核心功能测试

### 1. 用户认证

```bash
# 注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "password123"
  }'

# 登录（获取 Token）
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123"
  }'
```

### 2. 商品管理

```bash
# 获取商品列表
curl http://localhost:8080/api/v1/products

# 搜索商品
curl "http://localhost:8080/api/v1/products/search?keyword=手机"
```

### 3. WebSocket 推送

```bash
# 安装 wscat
npm install -g wscat

# 连接 WebSocket
wscat -c "ws://localhost:8080/ws" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🛠️ 常用命令

### Docker Compose

```bash
docker-compose up -d      # 启动
docker-compose ps         # 状态
docker-compose logs -f    # 日志
docker-compose down       # 停止
```

### Makefile

```bash
make build       # 编译
make run         # 运行
make test        # 测试
make docker-up   # Docker 启动
make help        # 查看所有命令
```

### 测试脚本

```bash
./scripts/test_api.sh     # API 自动化测试
```

## 📊 项目统计

- **文件数**: 40+ 个
- **代码量**: 5000+ 行
- **文档**: 7 篇完整文档
- **测试**: 单元测试 + 基准测试

## 🏗️ 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.21+ |
| 框架 | Gin | v1.9+ |
| ORM | GORM | v2.0+ |
| 数据库 | PostgreSQL | 15 |
| 缓存 | Redis | 7 |

## 📁 项目结构

```
shoppee/
├── cmd/                # 应用入口
│   └── api/main.go    # 主程序
├── internal/          # 内部代码
│   ├── config/        # 配置管理
│   ├── handler/       # HTTP 处理器
│   ├── service/       # 业务逻辑
│   ├── models/        # 数据模型
│   ├── middleware/    # 中间件
│   ├── router/        # 路由
│   └── websocket/     # WebSocket
├── pkg/               # 公共库
│   ├── jwt/          # JWT 工具
│   ├── logger/       # 日志
│   └── response/     # 响应封装
├── scripts/          # 脚本
├── docs/             # 文档（.md 文件）
├── Dockerfile        # Docker 构建
├── docker-compose.yml # 服务编排
├── Makefile          # 自动化构建
└── go.mod            # 依赖管理
```

## 🎓 学习路径

### 初级（了解项目）
1. 阅读 `QUICK_START.md` 启动项目
2. 测试 API 接口
3. 查看代码结构

### 中级（深入理解）
1. 阅读 `ARCHITECTURE.md` 了解架构
2. 查看核心代码实现
3. 运行单元测试

### 高级（扩展优化）
1. 阅读 `PERFORMANCE.md` 学习优化
2. 尝试添加新功能
3. 部署到生产环境

## 💡 使用建议

### 学习用途
- Go Web 开发学习
- 微服务架构参考
- 高并发系统设计
- Docker 部署实践

### 生产使用
⚠️ **部署前必须修改**：
- [ ] 数据库密码（`.env` 中的 `DB_PASSWORD`）
- [ ] JWT 密钥（`.env` 中的 `JWT_SECRET`）
- [ ] 关闭 DEBUG 模式（`APP_DEBUG=false`）
- [ ] 配置 HTTPS 证书
- [ ] 设置防火墙规则

## 🆘 遇到问题？

### 常见问题

**Q: 端口被占用**
```bash
# 修改 .env 中的端口
APP_PORT=8081
```

**Q: 数据库连接失败**
```bash
# 检查数据库状态
docker-compose ps postgres
docker-compose logs postgres
```

**Q: Redis 连接失败**
```bash
# 测试 Redis
redis-cli ping
```

### 获取帮助

1. 查看项目文档
2. 阅读代码注释
3. 提交 GitHub Issue

## 🎯 下一步

1. ✅ 阅读 [QUICK_START.md](QUICK_START.md)
2. ✅ 启动项目并测试 API
3. ✅ 查看 [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
4. ✅ 深入学习核心代码

## 🌟 项目亮点

- 🚀 **高性能**：协程池、缓存优化
- 🔒 **高安全**：JWT、加密、限流
- 📦 **易部署**：Docker 一键启动
- 📚 **文档全**：7 篇完整文档
- 🎯 **可扩展**：模块化设计

---

## 📞 联系方式

- **项目主页**: https://github.com/yourusername/shoppee
- **问题反馈**: GitHub Issues
- **技术支持**: support@shoppee.com

---

**祝您使用愉快！** 🎉

如果觉得项目有帮助，请给个 Star ⭐
