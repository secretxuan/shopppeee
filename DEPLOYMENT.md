# 部署指南

本文档提供详细的部署步骤和最佳实践。

## 🐳 Docker 部署（推荐）

### 1. 使用 Docker Compose 一键部署

```bash
# 克隆项目
git clone https://github.com/yourusername/shoppee.git
cd shoppee

# 配置环境变量
cp .env.example .env
vim .env  # 修改数据库密码、JWT密钥等

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 检查服务状态
docker-compose ps
```

### 2. 单独构建镜像

```bash
# 构建 Go 应用镜像
docker build -t shoppee:latest .

# 查看镜像大小
docker images shoppee

# 运行容器
docker run -d \
  --name shoppee-app \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e REDIS_HOST=your-redis-host \
  shoppee:latest
```

### 3. 镜像优化说明

本项目使用多阶段构建优化镜像大小：

- **编译阶段**：使用 `golang:1.21-alpine` 编译 Go 程序
- **运行阶段**：使用 `alpine:latest` 最小化镜像
- **静态编译**：`CGO_ENABLED=0` 避免依赖 C 库
- **编译优化**：`-ldflags="-w -s"` 去除调试信息

最终镜像大小约 **15-20MB**。

## 📦 传统部署

### 1. 编译可执行文件

```bash
# 本地编译
make build

# 交叉编译 Linux 版本（在 Mac/Windows 上）
make build-linux

# 输出文件位于 ./bin/shoppee
```

### 2. 系统要求

- Go 1.21+ （仅编译时需要）
- PostgreSQL 15+
- Redis 7+
- 系统：Linux/macOS/Windows

### 3. 手动部署步骤

#### 步骤 1：安装依赖

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql-15 redis-server

# CentOS/RHEL
sudo yum install postgresql15-server redis

# macOS
brew install postgresql@15 redis
```

#### 步骤 2：配置数据库

```bash
# 创建数据库
sudo -u postgres psql
CREATE DATABASE shoppee;
CREATE USER shoppee WITH PASSWORD 'your-password';
GRANT ALL PRIVILEGES ON DATABASE shoppee TO shoppee;
\q

# 导入初始化脚本
psql -U shoppee -d shoppee -f scripts/init.sql
```

#### 步骤 3：配置应用

```bash
# 创建配置文件
cp .env.example .env
vim .env

# 修改关键配置
APP_ENV=production
APP_DEBUG=false
DB_HOST=localhost
DB_PASSWORD=your-db-password
JWT_SECRET=your-super-secret-key-change-this
```

#### 步骤 4：启动应用

```bash
# 直接运行
./bin/shoppee

# 或使用 systemd（推荐）
sudo cp scripts/shoppee.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable shoppee
sudo systemctl start shoppee
sudo systemctl status shoppee
```

### 4. Systemd 服务配置

创建 `/etc/systemd/system/shoppee.service`：

```ini
[Unit]
Description=Shoppee E-Commerce Service
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=shoppee
WorkingDirectory=/opt/shoppee
ExecStart=/opt/shoppee/bin/shoppee
Restart=on-failure
RestartSec=5s

# 环境变量
Environment="APP_ENV=production"
Environment="APP_PORT=8080"

# 资源限制
LimitNOFILE=65535
LimitNPROC=65535

[Install]
WantedBy=multi-user.target
```

## 🚀 云平台部署

### AWS ECS 部署

```bash
# 1. 构建并推送镜像到 ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

docker tag shoppee:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/shoppee:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/shoppee:latest

# 2. 创建 ECS 任务定义（使用 AWS 控制台或 CLI）
# 3. 创建 ECS 服务并关联负载均衡器
```

### 阿里云 ACK（Kubernetes）部署

参考 Kubernetes 部署章节。

## ☸️ Kubernetes 部署

### 1. 准备配置文件

创建 `k8s/deployment.yaml`：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: shoppee
  labels:
    app: shoppee
spec:
  replicas: 3
  selector:
    matchLabels:
      app: shoppee
  template:
    metadata:
      labels:
        app: shoppee
    spec:
      containers:
      - name: shoppee
        image: shoppee:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: shoppee-secrets
              key: db-host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: shoppee-secrets
              key: db-password
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: shoppee-service
spec:
  selector:
    app: shoppee
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

### 2. 部署到 Kubernetes

```bash
# 创建命名空间
kubectl create namespace shoppee

# 创建 Secret
kubectl create secret generic shoppee-secrets \
  --from-literal=db-host=postgres-service \
  --from-literal=db-password=your-password \
  --from-literal=jwt-secret=your-jwt-secret \
  -n shoppee

# 部署应用
kubectl apply -f k8s/deployment.yaml -n shoppee

# 查看状态
kubectl get pods -n shoppee
kubectl get svc -n shoppee

# 查看日志
kubectl logs -f deployment/shoppee -n shoppee
```

## 🔧 反向代理配置

### Nginx 配置

创建 `/etc/nginx/sites-available/shoppee`：

```nginx
upstream shoppee_backend {
    server 127.0.0.1:8080;
    # 如果有多个实例，添加更多 server
    # server 127.0.0.1:8081;
    # server 127.0.0.1:8082;
}

server {
    listen 80;
    server_name shoppee.com www.shoppee.com;

    # HTTPS 重定向
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name shoppee.com www.shoppee.com;

    # SSL 证书配置
    ssl_certificate /etc/letsencrypt/live/shoppee.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/shoppee.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 日志
    access_log /var/log/nginx/shoppee_access.log;
    error_log /var/log/nginx/shoppee_error.log;

    # 代理配置
    location / {
        proxy_pass http://shoppee_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # WebSocket 支持
    location /ws {
        proxy_pass http://shoppee_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        # WebSocket 超时
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }

    # 静态文件缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
        proxy_pass http://shoppee_backend;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/shoppee /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 🔐 生产环境安全配置

### 1. 环境变量配置

**强制修改的配置：**
- `JWT_SECRET`：使用强随机字符串（至少 32 字符）
- `DB_PASSWORD`：数据库密码
- `APP_DEBUG`：设置为 `false`

**推荐修改的配置：**
- `CORS_ALLOWED_ORIGINS`：限制允许的域名
- `LOG_LEVEL`：设置为 `info` 或 `warn`

### 2. 防火墙配置

```bash
# Ubuntu/Debian (UFW)
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# 禁止直接访问应用端口
sudo ufw deny 8080/tcp
```

### 3. SSL/TLS 证书

使用 Let's Encrypt 免费证书：

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d shoppee.com -d www.shoppee.com

# 自动续期
sudo certbot renew --dry-run
```

## 📊 监控和日志

### 1. 日志管理

```bash
# 查看应用日志
tail -f ./logs/app.log

# 使用 journalctl 查看 systemd 日志
sudo journalctl -u shoppee -f

# 日志轮转配置 /etc/logrotate.d/shoppee
/opt/shoppee/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    missingok
    create 0644 shoppee shoppee
}
```

### 2. 性能监控

推荐工具：
- **Prometheus + Grafana**：指标监控
- **ELK Stack**：日志聚合分析
- **APM**：如 New Relic、Datadog

## 🔄 更新和回滚

### 更新应用

```bash
# Docker Compose
docker-compose pull app
docker-compose up -d app

# Kubernetes
kubectl set image deployment/shoppee shoppee=shoppee:v2.0 -n shoppee
kubectl rollout status deployment/shoppee -n shoppee
```

### 回滚版本

```bash
# Kubernetes
kubectl rollout undo deployment/shoppee -n shoppee

# Docker Compose
docker-compose down
docker-compose up -d
```

## 🧪 健康检查

```bash
# 检查服务状态
curl http://localhost:8080/health

# 预期响应
{"status":"ok","app":"Shoppee"}
```

## 📞 故障排查

### 1. 应用无法启动

```bash
# 检查日志
docker-compose logs app
# 或
sudo journalctl -u shoppee -n 100

# 常见问题：
# - 数据库连接失败：检查 DB_HOST, DB_PASSWORD
# - 端口占用：lsof -i :8080
# - 权限问题：检查文件权限
```

### 2. 数据库连接问题

```bash
# 测试数据库连接
psql -h localhost -U postgres -d shoppee

# 检查 PostgreSQL 状态
sudo systemctl status postgresql
```

### 3. Redis 连接问题

```bash
# 测试 Redis 连接
redis-cli ping

# 检查 Redis 状态
sudo systemctl status redis
```

## 🎯 性能调优

### 1. 数据库优化

```sql
-- 创建必要索引
CREATE INDEX idx_products_category_status ON products(category_id, status);
CREATE INDEX idx_orders_user_created ON orders(user_id, created_at DESC);

-- 分析查询性能
EXPLAIN ANALYZE SELECT * FROM products WHERE status = 'active';
```

### 2. Go 应用优化

```bash
# 编译时优化
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app

# 运行时优化（环境变量）
GOMAXPROCS=4  # CPU 核心数
GOGC=100      # GC 触发百分比
```

### 3. 连接池配置

修改代码中的连接池参数：

```go
sqlDB.SetMaxIdleConns(50)    // 最大空闲连接
sqlDB.SetMaxOpenConns(200)   // 最大打开连接
sqlDB.SetConnMaxLifetime(time.Hour)
```

## 📝 检查清单

部署前请确认：

- [ ] 修改了所有默认密码和密钥
- [ ] 配置了 HTTPS 证书
- [ ] 设置了防火墙规则
- [ ] 配置了日志轮转
- [ ] 设置了数据库备份
- [ ] 配置了监控告警
- [ ] 测试了健康检查接口
- [ ] 验证了 WebSocket 连接
- [ ] 执行了压力测试
- [ ] 准备了回滚方案

## 🆘 技术支持

遇到问题请查看：
- GitHub Issues
- 项目文档
- 社区论坛

---

祝部署顺利！🎉
