package main

// 这个可以完成身份验证的功能

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123456"},
	"austin": gin.H{"email": "austin@example.com", "phone": "88888"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "666666"},
}

func main() {
	r := gin.Default()

	// 路由组使用 gin.BasicAuth() 中间件
	// gin.Accounts 是map[string]string的一种快捷方式
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello",
		"manu":   "world",
	}))

	//  /admin/secrets 端点
	// 触发 "localhost:8080/admin/secrets"
	authorized.GET("/secrets", func(c *gin.Context) {
		// 获取用户，它是由 BasicAuth 中间件设置的
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{
				"user": user, "secret": secret,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"user": user, "secret": secret,
			})
		}
	})

	// 监听8080端口
	r.Run(":8080")
}
