# 🎯 Shoppee 项目功能补全计划

## 当前状态

**总体完成度**: 35% → 目标 95%

### ✅ 已完成
- 用户认证（注册/登录）
- 商品浏览（列表/详情/搜索）
- WebSocket 实时通信
- 分类管理（新增 ✨）
- 购物车管理（新增 ✨）

### 🔨 正在添加
1. ✅ 分类管理 (category_handler.go + category_service.go)
2. ✅ 购物车管理 (cart_handler.go + cart_service.go)
3. 🔄 订单管理
4. 🔄 收货地址管理
5. 🔄 商品完整CRUD
6. 🔄 支付管理
7. 🔄 商品评价
8. 🔄 用户个人中心

## 文件清单

### 新增 Handler 文件
- [x] `/internal/handler/category_handler.go` - 分类管理
- [x] `/internal/handler/cart_handler.go` - 购物车管理
- [ ] `/internal/handler/order_handler.go` - 订单管理
- [ ] `/internal/handler/address_handler.go` - 收货地址
- [ ] `/internal/handler/review_handler.go` - 商品评价
- [ ] `/internal/handler/payment_handler.go` - 支付管理

### 新增 Service 文件
- [x] `/internal/service/category_service.go` - 分类服务
- [x] `/internal/service/cart_service.go` - 购物车服务
- [ ] `/internal/service/order_service.go` - 订单服务
- [ ] `/internal/service/address_service.go` - 地址服务
- [ ] `/internal/service/review_service.go` - 评价服务
- [ ] `/internal/service/payment_service.go` - 支付服务

### 需要更新的文件
- [ ] `/internal/router/router.go` - 添加新路由
- [ ] `/internal/handler/product_handler.go` - 添加CRUD方法
- [ ] `/internal/service/product_service.go` - 添加CRUD方法

## API 接口清单

### 商品管理（管理员）
```
POST   /api/v1/products           # 创建商品
PUT    /api/v1/products/:id       # 更新商品
DELETE /api/v1/products/:id       # 删除商品
PUT    /api/v1/products/:id/status # 上下架商品
```

### 分类管理
```
GET    /api/v1/categories          # 分类列表/树
POST   /api/v1/categories          # 创建分类（管理员）
PUT    /api/v1/categories/:id      # 更新分类（管理员）
DELETE /api/v1/categories/:id      # 删除分类（管理员）
```

### 购物车管理
```
GET    /api/v1/cart                # 获取购物车
POST   /api/v1/cart/items          # 添加商品
PUT    /api/v1/cart/items/:id      # 更新数量
DELETE /api/v1/cart/items/:id      # 删除商品
PUT    /api/v1/cart/items/:id/select # 切换选中
POST   /api/v1/cart/clear          # 清空购物车
```

### 订单管理
```
POST   /api/v1/orders              # 创建订单
GET    /api/v1/orders              # 订单列表
GET    /api/v1/orders/:id          # 订单详情
PUT    /api/v1/orders/:id/cancel   # 取消订单
PUT    /api/v1/orders/:id/pay      # 支付订单
PUT    /api/v1/orders/:id/complete # 确认收货
PUT    /api/v1/orders/:id/ship     # 发货（管理员）
```

### 收货地址管理
```
GET    /api/v1/addresses           # 地址列表
POST   /api/v1/addresses           # 添加地址
PUT    /api/v1/addresses/:id       # 更新地址
DELETE /api/v1/addresses/:id       # 删除地址
PUT    /api/v1/addresses/:id/default # 设置默认
```

### 商品评价
```
POST   /api/v1/products/:id/reviews # 发表评价
GET    /api/v1/products/:id/reviews # 商品评价列表
GET    /api/v1/reviews              # 我的评价
PUT    /api/v1/reviews/:id          # 更新评价
DELETE /api/v1/reviews/:id          # 删除评价
```

### 支付管理
```
POST   /api/v1/payments             # 创建支付
POST   /api/v1/payments/callback    # 支付回调
GET    /api/v1/payments/:id         # 支付详情
```

### 用户个人中心
```
PUT    /api/v1/user/profile         # 更新个人信息
PUT    /api/v1/user/password        # 修改密码
PUT    /api/v1/user/avatar          # 更新头像
```

## 实施步骤

### 第一阶段：核心交易流程（当前）
1. [x] 分类管理 - 支持商品分类
2. [x] 购物车管理 - 用户加购商品
3. [ ] 收货地址管理 - 订单必需
4. [ ] 订单管理 - 核心业务流程
5. [ ] 商品CRUD - 管理员上架商品

### 第二阶段：增强功能
6. [ ] 支付管理 - 模拟支付流程
7. [ ] 商品评价 - 用户反馈
8. [ ] 个人中心 - 用户信息管理

### 第三阶段：前端对接
9. [ ] 前端商品管理页面
10. [ ] 前端订单管理页面
11. [ ] 前端个人中心页面

## 下一步行动

**立即执行**：
1. 创建订单相关文件（order_handler.go, order_service.go）
2. 创建地址相关文件（address_handler.go, address_service.go）
3. 扩展商品处理器（添加CRUD方法）
4. 更新路由配置

**预计完成时间**: 约 30 分钟

---

更新时间: 2025-12-19
