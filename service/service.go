package service

import "personal_blog/service/system"

type Group struct {
	SystemServiceSupplier system.Supplier
}

var GroupApp *Group
