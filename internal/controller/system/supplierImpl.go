package system

type controllerSupplier struct {
	refreshTokenCtrl *RefreshTokenCtrl
	baseCtrl         *BaseCtrl
}

func (c *controllerSupplier) GetRefreshTokenCtrl() *RefreshTokenCtrl {
	return c.refreshTokenCtrl
}
func (c *controllerSupplier) GetBaseCtrl() *BaseCtrl {
	return c.baseCtrl
}
