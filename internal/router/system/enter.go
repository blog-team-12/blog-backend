package system

type RouterGroup struct {
	RefreshTokenRouter
	BaseRouter
	UserRouter
	ImageRouter
	ArticleRouter
}
