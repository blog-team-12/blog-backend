package system

type controllerSupplier struct {
	refreshTokenApi *RefreshTokenApi
}

func (c *controllerSupplier) GetRefreshTokenApi() *RefreshTokenApi {
	return c.refreshTokenApi
}
