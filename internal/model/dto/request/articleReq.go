package request

// ArticleCreateReq 创建文章请求体
type ArticleCreateReq struct {
	Cover        string   `json:"cover" binding:"required"`      // 封面url
	Title        string   `json:"title" binding:"required"`      // 标题
	Category     string   `json:"category" binding:"required"`   // 列表/专栏
	Tags         []string `json:"tags" binding:"required"`       // 标记
	Abstract     string   `json:"abstract" binding:"required"`   // 摘要
	Content      string   `json:"content" binding:"required"`    // 文章内容
	VisibleRange uint     `json:"visible_range" bind:"required"` // 1-"全部可见"/2-"仅我可见"
}

// ArticleDeleteReq 删除文章请求体
type ArticleDeleteReq struct {
	IDs []string `json:"ids" binding:"required"` // elasticsearch中每篇文章对应的ID
}

// ArticleUpdateReq 更新文章请求体
type ArticleUpdateReq struct {
	ID           string   `json:"id" binding:"required"`         // elasticsearch中每篇文章对应的ID
	Title        string   `json:"title" binding:"required"`      // 标题
	Category     string   `json:"category" binding:"required"`   // 列表/专栏
	Tags         []string `json:"tags" binding:"required"`       // 标记
	Abstract     string   `json:"abstract" binding:"required"`   // 摘要
	Content      string   `json:"content" binding:"required"`    // 文章内容
	VisibleRange uint     `json:"visible_range" bind:"required"` // 1-"全部可见"/2-"仅我可见"
}

// ArticleListReq 查询文章请求体
type ArticleListReq struct {
	Title        *string `json:"title" form:"title"`            // 标题
	Category     *string `json:"category" form:"category"`      // 列表/专栏
	Abstract     *string `json:"abstract" form:"abstract"`      // 摘要
	VisibleRange uint    `json:"visible_range" bind:"required"` // 1-"全部可见"/2-"仅我可见"
	PageInfo             // 分页信息
}
