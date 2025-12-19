-- 初始化数据库脚本

-- 设置时区
SET timezone = 'Asia/Shanghai';

-- 创建数据库（如果不存在）
SELECT 'CREATE DATABASE shoppee' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'shoppee')\gexec

-- 连接到shoppee数据库
\c shoppee;

-- 启用UUID扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 初始化一些基础数据（可选）
-- 注意：表结构由GORM自动迁移创建

-- 插入默认管理员用户（密码：admin123，需要启动应用后通过API注册）
-- INSERT INTO users (username, email, password, role, status, created_at, updated_at)
-- VALUES ('admin', 'admin@shoppee.com', '$2a$10$...', 'admin', 'active', NOW(), NOW())
-- ON CONFLICT (username) DO NOTHING;

-- 插入默认分类
-- INSERT INTO categories (name, description, status, sort, created_at, updated_at)
-- VALUES 
--   ('电子产品', '手机、电脑、数码产品', 'active', 1, NOW(), NOW()),
--   ('服装鞋包', '男装、女装、鞋类、箱包', 'active', 2, NOW(), NOW()),
--   ('食品饮料', '零食、饮料、生鲜', 'active', 3, NOW(), NOW()),
--   ('家居生活', '家具、家纺、厨具', 'active', 4, NOW(), NOW()),
--   ('图书音像', '图书、电子书、音乐', 'active', 5, NOW(), NOW())
-- ON CONFLICT DO NOTHING;

-- 创建索引优化查询性能
-- CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
-- CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
-- CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
-- CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

COMMIT;
