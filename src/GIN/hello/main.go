package main

import (
	"fmt"
	_ "html/template"
	"log"
	"net/http"
	_ "path"

	"github.com/gin-gonic/gin"
)

type LoginForm struct{
	User string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// protoexample结构体
type Test struct{
	Label *string
	Reps []int64
}

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

// 使用不同目录下名称相同的模板
func LoadPostsHTML(c *gin.Context){
	c.HTML(http.StatusOK, "posts/index.html", gin.H{
		"title":"Posts",
	})
}
func LoadUsersHTML(c *gin.Context){
	c.HTML(http.StatusOK, "users/index.html", gin.H{
		"title":"Users",
	})
}

// JSONP
// 使用 JSONP 向不同域的服务器请求数据。如果查询参数存在回调，则将回调添加到响应体中。
func GetJSONP(c *gin.Context){
	data := map[string]interface{}{
		"foo":"bar",
	}
	// /JSONP?callback=x
	// 将输出：x({\"foo\":\"bar\"})
	c.JSONP(http.StatusOK, data)
}
	
// Multipart/Urlencoded 绑定
func PostUrlencoded(c *gin.Context){
	// 可以使用显式绑定声明绑定 multipart form：
	// c.ShouldBindWith(&form, binding.Form)
	// 或者简单地使用 ShouldBind 方法自动绑定：
	var form LoginForm
	// 在这种情况下，将自动选择合适的绑定
	if c.ShouldBind(&form) == nil{
		if form.User == "user" && form.Password == "password"{
			c.JSON(200, gin.H{
				"status":"you are logged in",
			})
		}else{
			c.JSON(401, gin.H{
				"status":"unauthorized",
			})
		}
	}
}

// Multipart/Urlencoded 表单
func PostForm(c *gin.Context){
	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")

	c.JSON(200, gin.H{
		"status":"posted",
		"message":message,
		"nick":nick,
	})
}

// PureJSON
/**
 通常，JSON 使用 unicode 替换特殊 HTML 字符，例如 < 变为 \ u003c。
 如果要按字面对这些字符进行编码，则可以使用 PureJSON。Go 1.6 及更低版本无法使用此功能。
*/ 
func GetUnicode(c *gin.Context){
	c.JSON(200, gin.H{
		"html":"<b>hello,world!</b>",
	})
}
func GetPureJSON(c *gin.Context){
	c.PureJSON(200, gin.H{
		"html":"<b>hello,world!</b>",
	})
}

// Query 和 post form
func PostQuery(c *gin.Context){
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	name := c.PostForm("name")
	message := c.PostForm("message")

	fmt.Println(id,page,name,message)
}

// SecureJSON
/**
使用 SecureJSON 防止 json 劫持。
如果给定的结构是数组值，则默认预置 "while(1)," 到响应体。
*/ 
func GetSecureJSON(c *gin.Context){
	// 你也可以使用自己的 SecureJSON 前缀
	// r.SecureJsonPrefix(")]}',\n")
	names := []string{"lena","austin","zhangsan"}
	// 将输出：while(1);["lena","austin","foo"]
	c.SecureJSON(http.StatusOK, names)
}

// 结构体json渲染
func GetStructJson(c *gin.Context){
	// 使用一个结构体
	var msg struct{
		Name string `json:"user"`
		Message string
		Number int
	}
	msg.Name = "zhangsan"
	msg.Message = "hello world"
	msg.Number = 123
	// 注意 msg.Name在json中变成了 "user"
	c.JSON(http.StatusOK, msg)
}
// XML渲染
func GetXML(c *gin.Context){
	c.XML(http.StatusOK, gin.H{
		"message":"hello",
		"status":http.StatusOK,
	})
}
// YAML渲染
func GetYAML(c *gin.Context){
	c.YAML(http.StatusOK, gin.H{
		"message":"hey",
		"status":http.StatusAccepted,
	})
}
// ProtoBuf渲染
func GetProtobuf(c *gin.Context){
	reps := []int64{int64(1),int64(2)}
	label := "test"
	// 结构体定义在最上面
	data := &Test{
		Label: &label,
		Reps: reps,
	}
	// 请注意，数据在响应中变为二进制数据
	// 将输出被 Test protobuf
	// 序列化了的数据
	c.ProtoBuf(http.StatusOK, data)

}

// 文件上传
// 单文件
func UpLoadFile(c *gin.Context){
	// 单文件
	file,_ := c.FormFile("file")
	log.Println(file.Filename)
	
	// 上传文件至指定目录
	// 可以通过path包来辅助解析
	// dst := path.Join("./data/",file.Filename)
	// 也可以直接 目录/文件名.文件类型
	dst := "./data/"+file.Filename
	c.SaveUploadedFile(file,dst)

	c.String(http.StatusOK,fmt.Sprintf("'%s' uploaded!",file.Filename))
}

func UpLoadFiles(c *gin.Context){

	form,_ := c.MultipartForm()
	files := form.File["file"]

	for _,file := range files{
		log.Println(file.Filename)
		dst := "./data/"+file.Filename
		// 上传文件至指定目录
		c.SaveUploadedFile(file,dst)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d file uploaded!", len(files)))
}

func main() {
	r := gin.Default()
	// 普通get
	r.GET("/ping", GetJson)
	r.GET("/somejson", GetAsciiJson)
	r.GET("/index", LoadHTML)

	// r.LoadHTMLGlob("views/templates/*")
	// 相当于解析views/templates/下的所有模板文件
	r.LoadHTMLGlob("views/templates/**/*")
	r.GET("/posts/index", LoadPostsHTML)
	r.GET("/users/index", LoadUsersHTML)

	// 自定义模板渲染器
	// html := template.Must(template.ParseFiles("file1","file2"))
	// r.SetHTMLTemplate(html)

	// 自定义分隔符
	// r.Delims("{[{","}]}")
	// r.LoadHTMLGlob("/path/to/templates")

	// JSONP
	// /JSONP?callback=x
	r.GET("/jsonp", GetJSONP)

	// Multipart/Urlencoded 绑定
	r.POST("/login", PostUrlencoded)
	// Multipart/Urlencoded 表单
	r.POST("/form_post", PostForm)

	// PureJSON
	// 提供Unicode实体
	r.GET("/json", GetUnicode)
	// 提供字面字符
	r.GET("purejson", GetPureJSON)
	
	// Query 和 post form
	r.POST("/post", PostQuery)

	//SecureJSON
	r.GET("/securejson", GetSecureJSON)

	// XML/JSON/YAML/ProtoBuf 渲染
	r.GET("/morejson", GetStructJson)
	r.GET("/somexml", GetXML)
	r.GET("/someyaml", GetYAML)
	r.GET("/someprotobuf", GetProtobuf)

	// 上传文件
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	r.MaxMultipartMemory = 8 << 20 //8 Mib
	// 上传单个文件
	r.POST("upload", UpLoadFile)
	// 上传多个文件
	r.POST("uploads", UpLoadFiles)

	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}