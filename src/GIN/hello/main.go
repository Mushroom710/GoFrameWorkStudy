package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetJson(c *gin.Context){
	c.JSON(200, gin.H{
		"hello":"world",
	})
}

/**
使用 AsciiJSON 生成具有转义的非 ASCII 字符的 ASCII-only JSON。
*/ 
func GetAsciiJson(c *gin.Context){
	data := map[string]interface{}{
		"lang":"GO语言",
		"tag":"<br/>",
	}

	// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
	c.AsciiJSON(http.StatusOK,data)
}

// HTML渲染
// 使用LoadHTMLGlob()或者LoadHTMLFiles()

func LoadHTML(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":"法外狂徒 张三",
	})
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("views/templates/*")

	r.GET("/ping", GetJson)
	r.GET("/somejson", GetAsciiJson)
	r.GET("/index", LoadHTML)

	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}