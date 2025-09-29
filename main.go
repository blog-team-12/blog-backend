package main

import (
	"personal_blog/global"
	"personal_blog/initialize"
)

func main() {
	// 初始化，并开启项目
	initialize.Init()
	global.Log.Info("1234")
}
