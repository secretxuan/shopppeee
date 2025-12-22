-- Shoppee 电商系统测试数据初始化脚本

-- 删除现有数据（可选）
TRUNCATE TABLE reviews, payments, order_items, orders, cart_items, carts, addresses, products, categories, users RESTART IDENTITY CASCADE;

-- 创建分类
INSERT INTO categories (name, description, icon, sort, status, created_at, updated_at) VALUES
  ('电子产品', '手机、电脑、数码配件等电子产品', '📱', 1, 'active', NOW(), NOW()),
  ('服装鞋包', '男装、女装、鞋子、箱包配饰', '👕', 2, 'active', NOW(), NOW()),
  ('食品饮料', '零食、饮料、生鲜水果', '🍎', 3, 'active', NOW(), NOW()),
  ('家居生活', '家具、家纺、日用百货', '🏠', 4, 'active', NOW(), NOW()),
  ('美妆个护', '化妆品、护肤品、个人护理', '💄', 5, 'active', NOW(), NOW());

-- 创建商品（电子产品）
INSERT INTO products (name, description, price, orig_price, stock, sku, category_id, status, view_count, sale_count, created_at, updated_at) VALUES
  ('iPhone 15 Pro 256GB', '苹果最新旗舰手机，A17仿生芯片，钛金属边框，超强性能', 7999.00, 8999.00, 50, 'IPHONE15PRO-256', 1, 'active', 1250, 87, NOW(), NOW()),
  ('MacBook Pro 14 M3', '苹果笔记本电脑，M3芯片，14.2英寸Liquid视网膜XDR显示屏', 12999.00, 14999.00, 30, 'MACBOOK14-M3', 1, 'active', 856, 43, NOW(), NOW()),
  ('AirPods Pro 2代', '苹果无线降噪耳机，主动降噪，空间音频', 1599.00, 1999.00, 100, 'AIRPODS2-001', 1, 'active', 2341, 156, NOW(), NOW()),
  ('iPad Air 10.9英寸', '平板电脑，M1芯片，支持Apple Pencil', 4599.00, 4999.00, 60, 'IPADAIR-M1', 1, 'active', 678, 34, NOW(), NOW()),
  ('Apple Watch Series 9', '智能手表，健康监测，全天候显示屏', 2999.00, 3299.00, 80, 'WATCH9-001', 1, 'active', 934, 67, NOW(), NOW()),
  ('小米14 Pro', '骁龙8 Gen3，徕卡光学镜头，120W快充', 4999.00, 5499.00, 120, 'MI14PRO-001', 1, 'active', 1567, 112, NOW(), NOW()),
  ('华为MatePad Pro', '12.2英寸平板，麒麟9000S，120Hz刷新率', 3999.00, 4499.00, 70, 'MATEPAD-PRO', 1, 'active', 543, 28, NOW(), NOW()),
  ('索尼WH-1000XM5', '无线降噪耳机，LDAC高清音质', 2299.00, 2999.00, 90, 'SONY-XM5', 1, 'active', 789, 56, NOW(), NOW());

-- 创建商品（服装鞋包）
INSERT INTO products (name, description, price, orig_price, stock, sku, category_id, status, view_count, sale_count, created_at, updated_at) VALUES
  ('优衣库男士T恤', '纯棉舒适，多色可选，基础百搭款', 99.00, 159.00, 500, 'UNIQLO-TSHIRT-M', 2, 'active', 3456, 567, NOW(), NOW()),
  ('ZARA女士连衣裙', '优雅时尚，适合春夏，收腰显瘦', 299.00, 499.00, 200, 'ZARA-DRESS-W', 2, 'active', 2134, 234, NOW(), NOW()),
  ('耐克Air Max运动鞋', '透气舒适，缓震设计，适合跑步', 799.00, 999.00, 150, 'NIKE-AIRMAX', 2, 'active', 1890, 178, NOW(), NOW()),
  ('阿迪达斯双肩包', '大容量，多功能口袋，防泼水', 399.00, 599.00, 180, 'ADIDAS-BAG', 2, 'active', 1234, 145, NOW(), NOW()),
  ('李维斯牛仔裤', '经典501款型，原色牛仔布', 499.00, 699.00, 220, 'LEVIS-501', 2, 'active', 2345, 289, NOW(), NOW());

-- 创建商品（食品饮料）
INSERT INTO products (name, description, price, orig_price, stock, sku, category_id, status, view_count, sale_count, created_at, updated_at) VALUES
  ('三只松鼠坚果礼盒', '每日坚果，混合装，健康零食', 89.00, 129.00, 300, 'SQUIRREL-NUT-BOX', 3, 'active', 4567, 678, NOW(), NOW()),
  ('可口可乐330ml*24罐', '经典可乐，整箱装，聚会必备', 49.00, 69.00, 500, 'COLA-330-24', 3, 'active', 3890, 456, NOW(), NOW()),
  ('奥利奥夹心饼干', '经典巧克力味，香浓可可', 15.90, 19.90, 800, 'OREO-COOKIE', 3, 'active', 5678, 890, NOW(), NOW()),
  ('农夫山泉矿泉水550ml*24瓶', '天然矿泉水，整箱装', 39.00, 59.00, 600, 'NONGFU-550-24', 3, 'active', 2345, 345, NOW(), NOW());

-- 创建商品（家居生活）
INSERT INTO products (name, description, price, orig_price, stock, sku, category_id, status, view_count, sale_count, created_at, updated_at) VALUES
  ('无印良品收纳盒', '简约设计，多尺寸组合，塑料材质', 59.00, 89.00, 400, 'MUJI-STORAGE', 4, 'active', 1678, 234, NOW(), NOW()),
  ('宜家懒人沙发', '舒适柔软，多色可选，客厅卧室适用', 599.00, 899.00, 80, 'IKEA-SOFA', 4, 'active', 2345, 123, NOW(), NOW()),
  ('飞利浦LED台灯', '护眼台灯，无频闪，多档调光', 199.00, 299.00, 150, 'PHILIPS-LAMP', 4, 'active', 1234, 156, NOW(), NOW()),
  ('戴森吸尘器V12', '无线手持，强力吸尘，轻巧便携', 2999.00, 3499.00, 50, 'DYSON-V12', 4, 'active', 890, 67, NOW(), NOW());

-- 创建商品（美妆个护）
INSERT INTO products (name, description, price, orig_price, stock, sku, category_id, status, view_count, sale_count, created_at, updated_ated) VALUES
  ('雅诗兰黛小棕瓶', '修护精华，抗老淡纹，经典明星产品', 799.00, 999.00, 120, 'ESTEE-SERUM', 5, 'active', 3456, 289, NOW(), NOW()),
  ('SK-II神仙水', '护肤精华水，改善肤质，焕发光彩', 1299.00, 1599.00, 80, 'SKII-WATER', 5, 'active', 2890, 234, NOW(), NOW()),
  ('兰蔻粉水', '柔肤水，补水保湿，温和舒缓', 399.00, 499.00, 200, 'LANCOME-TONER', 5, 'active', 2345, 345, NOW(), NOW()),
  ('欧莱雅洗发水', '修护受损发质，柔顺亮泽', 69.00, 89.00, 500, 'LOREAL-SHAMPOO', 5, 'active', 4567, 567, NOW(), NOW());

-- 创建测试用户（密码都是：password123）
-- 注意：这里的密码哈希是示例，实际使用时会在注册时自动生成
INSERT INTO users (username, email, password, phone, role, status, created_at, updated_at) VALUES
  ('admin', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye3wjd0w.2E/1PUQh3eeZ.3JJNsRAiS0K', '13800000000', 'admin', 'active', NOW(), NOW()),
  ('user001', 'user001@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye3wjd0w.2E/1PUQh3eeZ.3JJNsRAiS0K', '13800000001', 'user', 'active', NOW(), NOW()),
  ('user002', 'user002@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye3wjd0w.2E/1PUQh3eeZ.3JJNsRAiS0K', '13800000002', 'user', 'active', NOW(), NOW());

-- 创建测试地址
INSERT INTO addresses (user_id, name, phone, province, city, district, detail, is_default, created_at, updated_at) VALUES
  (2, '张三', '13800000001', '广东省', '深圳市', '南山区', '科技园南区xx路xx号xx室', true, NOW(), NOW()),
  (2, '李四', '13900000001', '广东省', '深圳市', '福田区', '华强北电子市场xx栋xx室', false, NOW(), NOW()),
  (3, '王五', '13800000002', '北京市', '北京市', '朝阳区', '三里屯soho xx号', true, NOW(), NOW());

-- 查看数据
SELECT '分类数据：' as info;
SELECT id, name, description, status FROM categories ORDER BY sort;

SELECT '商品数据：' as info;
SELECT id, name, price, stock, sku, category_id, status FROM products ORDER BY category_id, id LIMIT 10;

SELECT '用户数据：' as info;
SELECT id, username, email, role, status FROM users;

SELECT '地址数据：' as info;
SELECT id, user_id, name, phone, province, city, district, is_default FROM addresses;

SELECT '统计信息：' as info;
SELECT 
  (SELECT COUNT(*) FROM categories) as category_count,
  (SELECT COUNT(*) FROM products) as product_count,
  (SELECT COUNT(*) FROM users) as user_count,
  (SELECT COUNT(*) FROM addresses) as address_count;
