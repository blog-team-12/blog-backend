package service

import (
	"personal_blog/internal/service/system"
)

type Group struct {
	SystemServiceSupplier system.Supplier
}

var GroupApp *Group
