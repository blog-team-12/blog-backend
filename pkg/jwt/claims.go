package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"personal_blog/global"
	"personal_blog/internal/model/consts"
	"personal_blog/internal/model/dto/request"
)

// GetAccessToken 从请求头获取Access Token
func GetAccessToken(c *gin.Context) string {
	// 获取x-access-token头部值
	token := c.Request.Header.Get("x-access-token")
	return token
}

// GetRefreshToken 从cookie获取Refresh Token
func GetRefreshToken(c *gin.Context) string {
	// 尝试从cookie中获取refresh-token
	token := c.Request.Header.Get("x-refresh-token")
	return token
}

// GetUserID 从Gin的Context中获取JWT解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	// 首先尝试从Context中获取"claims"
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析Access Token
		if cl, err := GetClaims(c); err != nil {
			// 如果解析失败，返回0
			return 0
		} else {
			// 返回解析出来的用户ID
			return cl.UserID
		}
	} else {
		// 如果已存在claims，则直接返回用户ID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.UserID
	}
}

// GetUUID 从Gin的Context中获取JWT解析出来的用户UUID
func GetUUID(c *gin.Context) uuid.UUID {
	// 首先尝试从Context中获取"claims"
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析Access Token
		if cl, err := GetClaims(c); err != nil {
			// 如果解析失败，返回一个空UUID
			return uuid.UUID{}
		} else {
			// 返回解析出来的UUID
			return cl.UUID
		}
	} else {
		// 如果已存在claims，则直接返回UUID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.UUID
	}
}

// GetRoleID 从Gin的Context中获取JWT解析出来的用户角色ID
func GetRoleID(c *gin.Context) consts.RoleID {
	// 首先尝试从Context中获取"claims"
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析Access Token
		if cl, err := GetClaims(c); err != nil {
			// 如果解析失败，返回0
			return 0
		} else {
			// 返回解析出来的角色ID
			return cl.RoleID
		}
	} else {
		// 如果已存在claims，则直接返回角色ID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.RoleID
	}
}

// GetClaims 从Gin的Context中解析并获取JWT的Claims
func GetClaims(c *gin.Context) (*request.JwtCustomClaims, error) {
	// 获取请求头中的Access Token
	token := GetAccessToken(c)
	// 创建JWT实例
	j := NewJWT()
	// 解析Access Token
	claims, err := j.ParseAccessToken(token)
	if err != nil {
		// 如果解析失败，记录错误日志
		global.Log.Error("Failed to retrieve JWT parsing information from Gin's Context. Please check if the request header contains 'x-access-token' and if the claims structure is correct.", zap.Error(err))
	}
	return claims, err
}

/*
	交互方式
	// 客户端需要在请求头中同时携带两个token
	headers: {
	  'x-access-token': 'your_access_token',
	  'x-refresh-token': 'your_refresh_token'  // 仅在刷新接口中需要
	}
*/

// 修改过之前的函数

/*
SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool)
以下是每个参数的含义：
name（字符串）：
	含义：Cookie 的名称，客户端会用这个名称来标识 Cookie。
	例子：在这里，name 是 "x-refresh-token"，表示这是一个存储 Refresh Token 的 Cookie。
	作用：浏览器会将 Cookie 的名称和值一起存储，并在后续请求中发送匹配的 Cookie（如果路径和域名匹配）。
value（字符串）：
	含义：Cookie 的值，也就是存储的具体数据。
	例子：value 是 token（传入的 Refresh Token 字符串），可能是类似 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9... 的 JWT。
	作用：这是 Cookie 的实际内容，比如 Refresh Token，用于在客户端和服务器之间传递认证信息。
maxAge（整数）：
	含义：Cookie 的有效期（以秒为单位）。
	maxAge > 0：Cookie 在指定秒数后过期（例如 maxAge = 3600 表示 1 小时后过期）。
	maxAge = 0：Cookie 没有设置过期时间，默认为会话 Cookie（浏览器关闭后失效）。
	maxAge < 0：立即删除 Cookie。
path（字符串）：
	含义：Cookie 的生效路径，表示 Cookie 适用于服务器的哪些 URL 路径。
	例子：path = "/" 表示 Cookie 对整个网站（所有路径）生效。
	如果设置为 /controller，则 Cookie 只对 /controller 开头的路径生效.
	作用：限制 Cookie 的使用范围。/ 是最宽松的设置，适用于整个域名下的所有请求。
domain（字符串）：
	含义：Cookie 的生效域名，表示 Cookie 适用于哪些域名。
	例子：domain = host（从 net.SplitHostPort 提取的主机名，比如 example.com）。
	如果 domain = "example.com"，Cookie 适用于 example.com 及其子域名（如 controller.example.com）。
	如果不设置 domain（空字符串），Cookie 仅适用于当前请求的域名（不包括子域名）。
	作用：控制 Cookie 的共享范围。设置正确的 domain 确保 Cookie 在相关域名下可用，比如在生产环境中支持子域名共享。
secure（布尔值）：
	含义：是否只通过 HTTPS 传输 Cookie。
	例子：secure = false 表示 Cookie 可以通过 HTTP 或 HTTPS 传输。
	如果 secure = true，浏览器只会在 HTTPS 连接中发送 Cookie。
	作用：提高安全性，防止 Cookie 在不安全的 HTTP 连接中被窃听。设置为 false 可能用于开发环境（HTTP），但生产环境通常应设为 true。
httpOnly（布尔值）：
	含义：是否禁止客户端 JavaScript 访问 Cookie。
	例子：httpOnly = true 表示 Cookie 是 HttpOnly 的，浏览器中的 JavaScript（比如 document.cookie）无法访问该 Cookie。
	作用：增强安全性，防止跨站脚本攻击（XSS）。HttpOnly 确保 Cookie 只能通过 HTTP 请求发送到服务器，客户端脚本无法读取，适合存储敏感数据如 Refresh Token。
*/

//// SetRefreshToken 设置Refresh Token的cookie
//func SetRefreshToken(c *gin.Context, token string, maxAge int) {
//	// 获取请求的host，如果失败则取原始请求host
//	// 为了正确设置domain属性
//	host, _, err := net.SplitHostPort(c.Request.Host)
//	if err != nil {
//		host = c.Request.Host
//	}
//	/*
//		- 1.开发环境 ： localhost:8080 → 提取出 localhost
//		- 2.生产环境 ： controller.example.com:443 → 提取出 controller.example.com
//		- 3.IP访问 ：  192.168.1.100:3000 → 提取出 192.168.1.100
//		- 4.标准端口 ： example.com → 直接使用 example.com
//	*/
//	/*
//		net.SplitHostPort(c.Request.Host) 是一个标准库函数（net 包），用于将一个形如 host:port 的字符串拆分为主机名（host）和端口号（port）。
//		例如：
//			输入 localhost:8080 → 输出 host = "localhost", port = "8080"
//			输入 controller.example.com:443 → 输出 host = "controller.example.com", port = "443"
//		为什么需要拆分？ Cookie 的 Domain 属性只需要主机名（比如 example.com），不需要端口号。
//	*/
//	// 调用setCookie设置refresh-token
//	setCookie(c, "x-refresh-token", token, maxAge, host)
//
//}
//
//// ClearRefreshToken 清除Refresh Token的cookie
//func ClearRefreshToken(c *gin.Context) {
//	// 获取请求的host，如果失败则取原始请求host
//	host, _, err := net.SplitHostPort(c.Request.Host)
//	if err != nil {
//		host = c.Request.Host
//	}
//	// 调用setCookie设置cookie值为空并过期，删除refresh-token
//	setCookie(c, "x-refresh-token", "", -1, host)
//}
//
//// setCookie 设置指定名称和值的cookie
//func setCookie(c *gin.Context, name, value string, maxAge int, host string) {
//	// 判断host是否是IP地址
//	if net.ParseIP(host) != nil {
//		// 如果是IP地址，设置cookie的domain为“/”
//		c.SetCookie(name, value, maxAge, "/", "", false, true)
//	} else {
//		// 如果是域名，设置cookie的domain为域名
//		c.SetCookie(name, value, maxAge, "/", host, false, true)
//	}
//}
