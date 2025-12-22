# 🎨 Shoppee 前端快速启动

## ⚡ 一键启动

```bash
# 1. 安装依赖（首次运行）
npm install

# 2. 启动开发服务器
npm run dev
```

访问: **http://localhost:3000**

## ✅ 前置条件

确保后端服务已启动:

```bash
# 在项目根目录
docker compose up -d

# 验证后端
curl http://localhost:8080/health
```

## 🎯 功能特性

✨ **现代化 UI**
- Ant Design 5 组件库
- 主色调: #336699
- 响应式设计

🛍️ **核心功能**
- 商品浏览和搜索
- 购物车管理
- 用户认证（注册/登录）
- 个人中心

📱 **响应式支持**
- 手机端优化
- 平板适配
- 桌面端完整体验

## 📝 快速测试流程

1. **访问首页** → http://localhost:3000
2. **注册账号** → 点击右上角"注册"
3. **浏览商品** → 点击"商品"菜单
4. **加入购物车** → 点击商品卡片上的购物车图标
5. **查看购物车** → 点击顶部购物车图标

## 🔧 开发命令

```bash
# 开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview

# 代码检查
npm run lint
```

## 📚 技术栈

- React 18
- TypeScript
- Ant Design 5
- Zustand (状态管理)
- React Router 6
- Axios
- Vite

## 🎨 设计亮点

- 优美的轮播图
- 卡片式商品展示
- 流畅的动画效果
- 购物车实时更新
- 本地数据持久化

## 🐛 遇到问题？

查看详细文档: [FRONTEND_GUIDE.md](README.md)

---

**祝开发愉快！** 🚀
