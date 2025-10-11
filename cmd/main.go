package main

import (
	"personal_blog/global"
	Init "personal_blog/internal/init"
)

func main() {
	// 初始化，并开启项目
	Init.Init()
	global.Log.Info("1234")
}
