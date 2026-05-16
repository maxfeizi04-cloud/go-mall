package model

// ============================================================
// 请求结构体：用于接收和校验前端传来的 JSON 参数
// 每个字段的 binding tag 是 Gin 的参数校验规则
// ============================================================

// ========== 用户相关 ==========

// RegisterReq 注册请求参数
type RegisterReq struct {
	// 用户名
	// binding:"required" = 必填
	// min=3 最少 3 个字符,max=30 最多 30 个字符
	Username string `json:"username" binding:"required,min=3,max=30"`

	// 邮箱
	// email = 必须符合邮箱格式(如: a@b.com)
	Email string `json:"email" binding:"required,email"`

	// 密码
	// min=8 最少 8 个位,max=20 最多 20 位
	Password string `json:"password" binding:"required,min=8,max=20"`
}

// LoginReq 登录请求参数
type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required`
}

// LoginResp 登录成功后返回的数据
type LoginResp struct {
	Token string `json:"token"` // JWT Token，后续请求放在 Header 里
	User  User   `json:"user"`  // 用户基本信息
}

// ========== 商品相关 ==========

// CreateProductReq 创建商品请求参数
type CreateProductReq struct {
	Name        string  `json:"name" binding:"required,min=3,max=200"` // 商品名 (必填)
	Description string  `json:"description"`                           // 描述 (选填)
	Price       float64 `json:"price" binding:"required,gt=0"`         // 价格 (必填,且大于 0)
	Stock       uint32  `json:"stock" binding:"required,gte=0"`        // 库存 (大于等于 0,默认 0)
	CategoryID  uint64  `json:"category_id" binding:"required"`        // 分类 ID (必填)
	ImageURL    string  `json:"image-url"`                             // 图片 URL (选填)
}

// UpdateProductReq 更新商品请求参数
// 注意：所有字段都是指针类型（*string、*float64 等）
// 指针为 nil = 前端没传这个字段 = 不更新
// 指针非 nil = 前端传了这个字段 = 需要更新
// 这样可以实现"部分更新"，只改你想改的字段
type UpdateProductReq struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *uint32  `json:"stock"`
	CategoryID  *uint64  `json:"category_id"`
	ImageURL    *string  `json:"image-url"`
	Status      *int8    `json:"status"`
}

// ProductListReq 商品列表查询参数
// form tag 用于绑定 URL 查询参数（如 ?page=1&keyword=iPhone）
type ProductListReq struct {
	CategoryId uint64  `form:"category_id"`                                  // 按分类筛选
	keyword    string  `form:"keyword"`                                      // 按名称关键字搜索
	MinPrice   float64 `form:"min_price"`                                    // 最低价格
	MaxPrice   float64 `form:"max_price"`                                    // 最高价格
	Page       int     `form:"page,default=1" binding:"gte=1"`               // 页码,默认 1
	PageSize   int     `form:"page_size,default=10" binding:"gte=1,lte=100"` // 每页条数，默认 10，最大 100
}

// ========== 购物车相关 ==========

// AddCartReq 添加购物车请求参数
type AddCartReq struct {
	ProductId uint64 `json:"product_id" binding:"required"`     // 商品 ID（必填）
	Quantity  uint32 `json:"quantity" binding:"required,gte=1"` // 数量（必填，至少 1）
}

// UpdateCartReq 更新购物车数量
type UpdateCartReq struct {
	Quantity uint32 `json:"quantity" binding:"required,gte=1"` // 新的数量
}

// ========== 订单相关 ==========

// CreateOrderReq 创建订单请求参数
type CreateOrderReq struct {
	Address string `json:"address" binding:"required"` // 收货地址（必填）
	Remark  string `json:"remark"`                     // 备注（选填）
}
