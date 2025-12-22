# Shoppee 前端

基于 React 18 + TypeScript + Ant Design 5 构建的现代化电商前端应用。

## ✨ 特性

- 🎨 **现代化 UI** - 基于 Ant Design 5，优美的用户界面
- 📱 **响应式设计** - 完美适配移动端和桌面端
- ⚡ **高性能** - 使用 Vite 构建，快速的开发体验
- 🔐 **用户认证** - JWT 令牌认证，安全可靠
- 🛒 **购物车** - 本地持久化存储，数据不丢失
- 🎯 **TypeScript** - 类型安全，开发更高效
- 🌈 **主题定制** - 主色调 #336699，清新现代

## 📦 技术栈

- **React 18** - 最新的 React 版本
- **TypeScript** - 类型安全
- **Ant Design 5** - 企业级 UI 组件库
- **React Router 6** - 路由管理
- **Zustand** - 轻量级状态管理
- **Axios** - HTTP 请求
- **Vite** - 下一代前端构建工具
- **Dayjs** - 日期处理

## 🚀 快速开始

### 前置要求

- Node.js 16+
- npm/yarn/pnpm

### 安装依赖

```bash
cd frontend
npm install
# 或
yarn install
# 或
pnpm install
```

### 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 📁 项目结构

```
frontend/
├── src/
│   ├── api/              # API 接口
│   │   ├── axios.ts     # Axios 配置
│   │   ├── auth.ts      # 认证 API
│   │   └── product.ts   # 商品 API
│   ├── components/       # 通用组件
│   │   ├── Layout.tsx   # 布局组件
│   │   ├── ProductCard.tsx  # 商品卡片
│   │   └── PrivateRoute.tsx # 私有路由
│   ├── pages/           # 页面组件
│   │   ├── Home.tsx     # 首页
│   │   ├── Products.tsx # 商品列表
│   │   ├── ProductDetail.tsx # 商品详情
│   │   ├── Cart.tsx     # 购物车
│   │   ├── Login.tsx    # 登录
│   │   ├── Register.tsx # 注册
│   │   └── Profile.tsx  # 个人中心
│   ├── store/           # 状态管理
│   │   ├── useAuthStore.ts  # 认证状态
│   │   └── useCartStore.ts  # 购物车状态
│   ├── App.tsx          # 根组件
│   ├── main.tsx         # 入口文件
│   └── index.css        # 全局样式
├── index.html           # HTML 模板
├── package.json         # 依赖配置
├── tsconfig.json        # TypeScript 配置
├── vite.config.ts       # Vite 配置
└── README.md           # 说明文档
```

## 🎯 功能模块

### 1. 用户认证
- ✅ 用户注册
- ✅ 用户登录
- ✅ JWT 令牌管理
- ✅ 个人中心

### 2. 商品浏览
- ✅ 商品列表
- ✅ 商品详情
- ✅ 商品搜索
- ✅ 分类筛选
- ✅ 排序功能

### 3. 购物车
- ✅ 添加商品
- ✅ 修改数量
- ✅ 删除商品
- ✅ 本地持久化

### 4. 响应式设计
- ✅ 移动端适配
- ✅ 平板适配
- ✅ 桌面端优化

## 🔧 配置说明

### API 代理配置

在 `vite.config.ts` 中配置了 API 代理：

```typescript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

### 主题配置

在 `src/main.tsx` 中配置 Ant Design 主题：

```typescript
<ConfigProvider
  theme={{
    token: {
      colorPrimary: '#336699',
      borderRadius: 8,
      colorBgContainer: '#ffffff',
    },
  }}
>
```

## 🎨 设计规范

### 颜色规范

- **主色调**: #336699 (品牌蓝)
- **成功色**: #52c41a
- **警告色**: #faad14
- **错误色**: #ff4d4f
- **信息色**: #1890ff

### 间距规范

- **小间距**: 8px
- **中间距**: 16px
- **大间距**: 24px
- **超大间距**: 32px

### 圆角规范

- **默认圆角**: 8px
- **卡片圆角**: 8px
- **按钮圆角**: 6px

## 📱 响应式断点

- **手机**: < 768px
- **平板**: 768px - 1024px
- **桌面**: > 1024px

## 🔐 状态管理

使用 Zustand 进行状态管理：

- **useAuthStore**: 用户认证状态
- **useCartStore**: 购物车状态

数据持久化到 localStorage，刷新页面不丢失。

## 🚀 开发建议

1. **组件开发**: 遵循单一职责原则，组件尽量小而专注
2. **样式管理**: 使用 CSS Modules 或独立的 CSS 文件
3. **类型安全**: 充分利用 TypeScript 的类型检查
4. **性能优化**: 使用 React.memo、useMemo、useCallback 等优化性能
5. **代码规范**: 使用 ESLint 保持代码质量

## 📝 待办事项

- [ ] 订单管理模块
- [ ] 支付流程
- [ ] 商品评价
- [ ] 收货地址管理
- [ ] WebSocket 实时通知
- [ ] 多语言支持
- [ ] 暗黑模式

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
