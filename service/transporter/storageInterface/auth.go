package storageInterface

import (
	"github.com/gin-gonic/gin"
)

// JSIAuthMiddleware  JSI 的认证中间件
func JSIAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("uid", "tester")
		c.Next() // 后续的处理函数可以用过c.Get("filePath")来获取当前请求的用户信息
	}
}
