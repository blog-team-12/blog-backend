package system

type controllerSupplier struct {
	refreshTokenCtrl *RefreshTokenCtrl
	baseCtrl         *BaseCtrl
	userCtrl         *UserCtrl
}

func (c *controllerSupplier) GetRefreshTokenCtrl() *RefreshTokenCtrl {
	return c.refreshTokenCtrl
}
func (c *controllerSupplier) GetBaseCtrl() *BaseCtrl {
	return c.baseCtrl
}
func (c *controllerSupplier) GetUserCtrl() *UserCtrl {
	return c.userCtrl
}
