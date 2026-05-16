package model

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// BaseModel：所有表的公共字段
// 通过嵌入（embedding）的方式复用，避免每张表都写一遍
// ============================================================

type BaseModel struct {
	// 主键,类型用 uint64 而不是 int
	// uint64 范围更大(0 ~ 18.4 亿亿),适合高并发场景下自增 ID 的未来扩展
	ID uint64 `gorm:"primaryKey" json:"id"`

	// 自动填充创建时间, GORM 在 Insert 时自动写入当前时间
	// json 格式: 2024-01-15T10:30:00Z
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 自动填充更新时间, GORM 在每次 Save/Update 时自动写入时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 软删除标记
	// 不会真正执行 DELETE SQL,而是把 deleted_at 设为当前时间
	/// json:"-" 表示序列化时隐藏这个字段,前端不需要看到
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ============================================================
// 用户表（users）
// 存储注册用户的基本信息
// ============================================================

type User struct {
	BaseModel // 嵌入公共字段，等价于拥有 ID / CreatedAt / UpdatedAt / DeletedAt

	// 用户名
	// uniqueIndex: 创建唯一索引,保证用户名不重复
	// not null: 不能为空
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`

	// 邮箱,用于登录
	// uniqueIndex: 保证邮箱不重复
	Email string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`

	// 密码(bcrypt 加密后的密文,不是明文)
	// json:"-" 表示返回给前端是隐藏密码字段,即使加密后也不暴露
	Password string `gorm:"type:varchar(255);not null" json:"-"`

	// 头像 URL
	// default: '' 设置为空字符串,避免 NULL 值带来的麻烦
	Avatar string `gorm:"type:varchar(255);default:''" json:"avatar"`

	// 账户状态
	// 1 = 正常, 0 = 禁用
	// 用 int8 而不是 bool, 方便未来扩展更多的状态(如 2 = 待验证等)
	Status int8 `gorm:"default:1" json:"status"`
}

// TableName 指定表名
// GORM 默认会把结构体名转为蛇形复数 (User -> users),这里显示指定更清晰
func (User) TableName() string {
	return "users"
}

// ============================================================
// 商品分类表（categories）
// 支持二级分类：电子产品 → 手机 / 电脑
// 通过 parent_id 自关联实现树形结构
// ============================================================

type Category struct {
	BaseModel

	// 分类名称,如 "电子产品"、"手机"
	Name string `gorm:"type:varchar(50);not null" json:"name"`

	// 父分类 ID
	// 0 表示顶级分类 (没用父级)
	// 非 0 表示属于某个父分类,如 "手机" 的 parent_id = "电子产品" 的 ID
	ParentID uint64 `gorm:"default:0;index" json:"parent_id"`

	// 分类层级
	// 1 = 一级分类 (电子产品、服装)
	// 2 = 二级分类 (手机、电脑、男装、女装)
	Level int8 `gorm:"default:0" json:"level"`

	// 排序序号,数字越小越靠前
	// 如: 电子产品 sort_order=1,服装 sort_order=2
	SortOrder int `gorm:"default:0" json:"sort_order"`
}

func (Category) TableName() string {
	return "categories"
}

// ============================================================
// 商品表（products）
// 核心业务表，存储所有在售商品的信息
// ============================================================

type Product struct {
	BaseModel

	// 商品名称
	Name string `gorm:"type:varchar(200);not null" json:"name"`

	// 商品描述,用 text 类型支持长文本
	Description string `gorm:"type:text" json:"description"`

	// 价格，用 decimal(10,2) 保证精度
	// 10 位数字,其中 2 位小数,最大支持 99999999.99
	// 注意: Go 侧用 float64 接收,极端精度场景(如金融)应考虑用 shopspring/decimal
	Price float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	// 库存数量
	// uint32 无符号 32 位,最大约 42 亿, 足够用
	Stock uint32 `gorm:"not null;default:0" json:"stock"`

	// 所属分类 ID (外键)
	// index: 加索引,按分类查询是高频操作
	CategoryID uint64 `gorm:"not null;index" json:"category_id"`

	// 商品图片 URL
	ImageURL string `gorm:"type:varchar(500);default:''" json:"image_url"`

	// 商品状态
	// 1 = 上架(前台可见可买), 0 = 下架(前台不可见)
	Status int8 `gorm:"default:1" json:"status"`

	// 分类关联(不是数据库字段)
	// gorm:"foreignKey:CategoryID" 表示 CategoryID 字段做外键关联
	// json:"category,omitempty" 表示为空时不输出(避免 "category": null)
	// 使用 Preload("Category") 时 GORM 自动填充
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

// ============================================================
// 购物车表（cart_items）
// 每个用户的待购买商品列表
// ============================================================

type CartItem struct {
	BaseModel

	// 用户 ID
	UserID uint64 `gorm:"not null;uniqueIndex:uk_user_product" json:"user_id"`

	// 商品 ID
	ProductID uint64 `gorm:"not null;uniqueIndex:uk_user_product" json:"product_id"`

	// 购买数量,默认 1
	Quantity uint32 `gorm:"not null;default:1" json:"quantity"`

	// 商品关联 (用于查询购物车时带上商品详情)
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (CartItem) TableName() string {
	return "cart_items"
}

// ============================================================
// 关于联合唯一索引 uk_user_product 的说明：
//
// 在 UserID 和 ProductID 上都声明了同一个索引名 uk_user_product
// GORM 会创建一个联合唯一索引 (user_id, product_id)
// 效果：同一个用户不能把同一个商品重复加入购物车
// 如果想加数量，走"累加"逻辑而不是"新建"
// ============================================================

// ============================================================
// 订单表（orders）
// 用户下单后生成的订单记录
// ============================================================

type Order struct {
	BaseModel

	// 订单号,业务生成的唯一编号 (如: ORD20240115103000123456)
	// 不用自增 ID 暴露给用户,避免被遍历猜测
	OrderNo string `gorm:"type:varchar(64);uniqueIndex;not null" json:"order_no"`

	// 下单用户 ID
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// 订单总金额(所有商品价格 × 数量 的总和)
	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`

	// 订单状态流转：
	//   0 = 待付款（刚下单）
	//   1 = 已付款（用户完成支付）
	//   2 = 已发货（商家发货）
	//   3 = 已完成（用户确认收货）
	//   4 = 已取消（用户取消或超时）
	Status int8 `gorm:"default:0;index" json:"status"`

	// 收货地址
	Address string `gorm:"type:text;not null" json:"address"`

	// 订单备注
	Remark string `gorm:"type:varchar(255);default:''" json:"remark"`

	// 支付时间,指针类型 *time.Time
	// nil 表示还没支付,非 nil 表示已支付的具体时间
	PaidAt *time.Time `json:"paid_at,omitempty"`

	// 订单包含的商品列表（一对多关联）
	// 一个订单可以有多个订单项
	// json:"items,omitempty" 空时不输出
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

// ============================================================
// 订单商品项表（order_items）
// 记录每个订单里具体买了哪些商品
// ============================================================

type OrderItem struct {
	BaseModel

	// 所属订单 ID (外键,关联 orders 表)
	OrderID uint64 `gorm:"not null;index" json:"order_id"`

	// 商品 ID(指定 products 表)
	ProductID uint64 `gorm:"not null" json:"product_id"`

	// 商品名称快照
	// 下单时记录当时的商品名，之后即使商品改名，订单里显示的不变
	// 这是电商系统的核心设计思想：订单是"历史凭证"，不能被商品的修改影响
	ProductName string `gorm:"type:varchar(200);not null" json:"product_name"`

	// 商品价格快照
	// 下单时记录当时的价格，之后商品改价不影响已下单的金额
	ProductPrice float64 `gorm:"type:decimal(10,2);not null" json:"product_price"`

	// 购买数量
	Quantity uint32 `gorm:"not null" json:"quantity"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

// ============================================================
// 表关系总览（ER 图）：
//
// users (1) ──────── (N) orders
//    │                    │
//    │                    └── (1:N) order_items
//    │                              │
//    │                              └── product_id → products
//    │
//    └── (1:N) cart_items
//                 │
//                 └── product_id → products
//
// categories (1) ──── (N) products
// categories (1) ──── (N) categories (自关联，parent_id)
// ============================================================
