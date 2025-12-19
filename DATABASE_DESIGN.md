# 数据库设计文档

## ER 图（实体关系图）

```
┌─────────────┐       ┌──────────────┐       ┌─────────────┐
│    User     │1    n │   Address    │       │    Cart     │
│─────────────│───────│──────────────│       │─────────────│
│ id          │       │ id           │       │ id          │
│ username    │       │ user_id      │  1  1 │ user_id     │
│ email       │       │ receiver_name│◄──────│             │
│ password    │       │ phone        │       └─────────────┘
│ role        │       │ detail       │              │
│ status      │       │ is_default   │              │ 1
└─────────────┘       └──────────────┘              │
      │ 1                                           │
      │                                             │ n
      │ n                                    ┌──────────────┐
┌─────────────┐                              │  Cart_Item   │
│    Order    │                              │──────────────│
│─────────────│                              │ id           │
│ id          │                              │ cart_id      │
│ order_no    │                         ┌────│ product_id   │
│ user_id     │                         │    │ quantity     │
│ total_amount│                         │    │ selected     │
│ status      │                         │    └──────────────┘
│ pay_status  │                         │
└─────────────┘                         │
      │ 1                               │
      │                                 │
      │ n                               │
┌─────────────┐                         │
│ Order_Item  │                         │
│─────────────│                         │
│ id          │                         │
│ order_id    │                         │
│ product_id  │◄────────────────────────┘
│ quantity    │
│ price       │
└─────────────┘
      │
      │
      ▼ n
┌─────────────┐       ┌──────────────┐
│   Product   │n    1 │   Category   │
│─────────────│───────│──────────────│
│ id          │       │ id           │
│ name        │       │ name         │
│ description │       │ description  │
│ price       │       │ parent_id    │
│ stock       │       │ sort         │
│ category_id │       │ status       │
│ sku         │       └──────────────┘
│ status      │
└─────────────┘
      │ 1
      │
      │ n
┌─────────────┐
│   Review    │
│─────────────│
│ id          │
│ user_id     │
│ product_id  │
│ rating      │
│ content     │
└─────────────┘
```

## 表结构详细设计

### 1. users（用户表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 用户ID |
| username | VARCHAR(50) | UNIQUE, NOT NULL | 用户名 |
| email | VARCHAR(100) | UNIQUE, NOT NULL | 邮箱 |
| password | VARCHAR(255) | NOT NULL | 密码（bcrypt加密） |
| phone | VARCHAR(20) | | 手机号 |
| avatar | VARCHAR(255) | | 头像URL |
| role | VARCHAR(20) | DEFAULT 'user' | 角色：user/admin |
| status | VARCHAR(20) | DEFAULT 'active' | 状态：active/inactive/banned |
| last_login | TIMESTAMP | | 最后登录时间 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- UNIQUE INDEX: username, email
- INDEX: status

### 2. products（商品表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 商品ID |
| name | VARCHAR(200) | NOT NULL | 商品名称 |
| description | TEXT | | 商品描述 |
| price | DECIMAL(10,2) | NOT NULL | 售价 |
| orig_price | DECIMAL(10,2) | | 原价 |
| stock | INTEGER | NOT NULL, DEFAULT 0 | 库存 |
| sku | VARCHAR(100) | UNIQUE | SKU编码 |
| images | TEXT | | 图片JSON数组 |
| category_id | INTEGER | FOREIGN KEY | 分类ID |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 |
| view_count | INTEGER | DEFAULT 0 | 浏览次数 |
| sale_count | INTEGER | DEFAULT 0 | 销量 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- UNIQUE INDEX: sku
- INDEX: category_id, status, name
- INDEX: sale_count DESC（热销排序）

### 3. categories（分类表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 分类ID |
| name | VARCHAR(100) | NOT NULL | 分类名称 |
| description | TEXT | | 分类描述 |
| icon | VARCHAR(255) | | 图标URL |
| parent_id | INTEGER | FOREIGN KEY | 父分类ID |
| sort | INTEGER | DEFAULT 0 | 排序 |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- INDEX: parent_id, status

### 4. orders（订单表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 订单ID |
| order_no | VARCHAR(50) | UNIQUE, NOT NULL | 订单号 |
| user_id | INTEGER | FOREIGN KEY, NOT NULL | 用户ID |
| total_amount | DECIMAL(10,2) | NOT NULL | 总金额 |
| pay_amount | DECIMAL(10,2) | NOT NULL | 实付金额 |
| status | VARCHAR(20) | DEFAULT 'pending' | 订单状态 |
| pay_status | VARCHAR(20) | DEFAULT 'unpaid' | 支付状态 |
| pay_method | VARCHAR(20) | | 支付方式 |
| receiver_name | VARCHAR(50) | | 收货人 |
| receiver_phone | VARCHAR(20) | | 收货电话 |
| receiver_address | VARCHAR(255) | | 收货地址 |
| remark | TEXT | | 备注 |
| paid_at | TIMESTAMP | | 支付时间 |
| shipped_at | TIMESTAMP | | 发货时间 |
| completed_at | TIMESTAMP | | 完成时间 |
| cancelled_at | TIMESTAMP | | 取消时间 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- UNIQUE INDEX: order_no
- INDEX: user_id, status, created_at DESC

### 5. order_items（订单项表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 订单项ID |
| order_id | INTEGER | FOREIGN KEY, NOT NULL | 订单ID |
| product_id | INTEGER | FOREIGN KEY, NOT NULL | 商品ID |
| quantity | INTEGER | NOT NULL | 数量 |
| price | DECIMAL(10,2) | NOT NULL | 单价 |
| total_price | DECIMAL(10,2) | NOT NULL | 小计 |
| product_name | VARCHAR(200) | | 商品名称快照 |
| product_image | VARCHAR(255) | | 商品图片快照 |
| product_sku | VARCHAR(100) | | SKU快照 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- INDEX: order_id, product_id

### 6. carts（购物车表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 购物车ID |
| user_id | INTEGER | UNIQUE, FOREIGN KEY, NOT NULL | 用户ID |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- UNIQUE INDEX: user_id

### 7. cart_items（购物车项表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 购物车项ID |
| cart_id | INTEGER | FOREIGN KEY, NOT NULL | 购物车ID |
| product_id | INTEGER | FOREIGN KEY, NOT NULL | 商品ID |
| quantity | INTEGER | NOT NULL, DEFAULT 1 | 数量 |
| selected | BOOLEAN | DEFAULT true | 是否选中 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- INDEX: cart_id, product_id

### 8. addresses（收货地址表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 地址ID |
| user_id | INTEGER | FOREIGN KEY, NOT NULL | 用户ID |
| receiver_name | VARCHAR(50) | NOT NULL | 收货人 |
| phone | VARCHAR(20) | NOT NULL | 电话 |
| province | VARCHAR(50) | | 省份 |
| city | VARCHAR(50) | | 城市 |
| district | VARCHAR(50) | | 区县 |
| detail | VARCHAR(255) | NOT NULL | 详细地址 |
| is_default | BOOLEAN | DEFAULT false | 是否默认 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- INDEX: user_id

### 9. payments（支付记录表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 支付ID |
| order_id | INTEGER | UNIQUE, FOREIGN KEY, NOT NULL | 订单ID |
| transaction_no | VARCHAR(100) | | 第三方交易号 |
| pay_method | VARCHAR(20) | NOT NULL | 支付方式 |
| pay_amount | DECIMAL(10,2) | NOT NULL | 支付金额 |
| status | VARCHAR(20) | DEFAULT 'pending' | 支付状态 |
| paid_at | TIMESTAMP | | 支付时间 |
| refunded_at | TIMESTAMP | | 退款时间 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- UNIQUE INDEX: order_id

### 10. reviews（评价表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 评价ID |
| user_id | INTEGER | FOREIGN KEY, NOT NULL | 用户ID |
| product_id | INTEGER | FOREIGN KEY, NOT NULL | 商品ID |
| order_id | INTEGER | FOREIGN KEY | 订单ID |
| rating | INTEGER | NOT NULL | 评分（1-5） |
| content | TEXT | | 评价内容 |
| images | TEXT | | 图片JSON数组 |
| status | VARCHAR(20) | DEFAULT 'pending' | 审核状态 |
| created_at | TIMESTAMP | NOT NULL | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引：**
- PRIMARY KEY: id
- INDEX: user_id, product_id, status

## 性能优化建议

1. **索引优化**
   - 所有外键添加索引
   - 高频查询字段添加组合索引
   - 避免过多索引影响写入性能

2. **分区表**
   - 订单表按月份分区
   - 日志表按日期分区

3. **读写分离**
   - 主库负责写入
   - 从库负责查询

4. **缓存策略**
   - 热门商品缓存
   - 用户会话缓存
   - 分类树缓存

5. **数据归档**
   - 历史订单归档
   - 已删除数据定期清理
