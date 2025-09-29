package api

import "personal_blog/api/system"

// ApiGroup 控制器
type ApiGroup struct {
	SystemApiGroup system.Supplier
}

var ApiGroupApp *ApiGroup
