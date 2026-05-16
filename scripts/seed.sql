-- 切换到 go_mall 数据库
USE go_mall;

-- ============================================================
-- 插入分类数据（两级结构）
-- 第 1 行：电子产品（顶级，parent_id=0）
-- 第 3 行：手机（二级，parent_id=1 指向电子产品）
-- ============================================================

INSERT INTO categories (name, parent_id, level, sort_order) VALUES
('电子产品', 0, 1, 1),   -- id=1
('服装',     0, 1, 2),   -- id=2
('手机',     1, 2, 1),   -- id=3, 属于电子产品
('电脑',     1, 2, 2),   -- id=4, 属于电子产品
('男装',     2, 2, 1),   -- id=5, 属于服装
('女装',     2, 2, 2);   -- id=6, 属于服装

-- ============================================================
-- 插入商品数据
-- 每个商品关联一个分类（category_id 对应上面的分类 ID）
-- ============================================================

INSERT INTO products (name, description, price, stock, category_id, image_url, status)
VALUES
('iPhone 16 Pro',  'Apple 最新旗舰手机，A18 Pro 芯片',8999.00,  100, 3, '/images/iphone16.jpg', 1),
('MacBook Pro 14', 'M4 Pro 芯片，16GB 内存',14999.00, 50,  4, '/images/macbook.jpg',  1),
('经典白T恤',      '100% 纯棉，舒适透气',99.00,    500, 5, '/images/tshirt.jpg',   1),
('连衣裙',         '夏季新款，轻薄飘逸',299.00,   200, 6, '/images/dress.jpg',    1);

-- ============================================================
-- 验证数据
-- ============================================================

-- 查看分类（树形结构）
-- SELECT id, name, parent_id, level FROM categories ORDER BY level, sort_order;

-- 查看商品（带分类名）
-- SELECT p.id, p.name, p.price, p.stock, c.name AS category
-- FROM products p
-- LEFT JOIN categories c ON p.category_id = c.id;