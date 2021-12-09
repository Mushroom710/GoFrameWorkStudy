package controllers

//GetString()中的参数为html页面中的name属性
import (
	"bytes"
	"encoding/gob"
	"hello/models"
	"math"
	"os"
	"path"
	"strconv"

	// "strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
)

type ArticleController struct {
	beego.Controller
}

// 处理下拉框的改变的请求
func(this *ArticleController)HandlerSelect(){
	// 接收数据
	typeName := this.GetString("select")
	// 处理数据
	if typeName == ""{
		beego.Info("获取下拉框数据失败")
		return
	}
	// beego.Info(typeName)
	// 查询数据
	o := orm.NewOrm()
	var articles[] models.Article
	//这里不加 relatedsel 函数会出现惰性查询
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	// beego.Info(articles)
}


func (this *ArticleController) ShowArticleList() {
	/*
	这一步操作已经通过 过滤器函数处理
	// 先查询是否登录,只有登录成功才能访问后台
	// 返回的是一个引用
	// name := this.GetSession("userName")
	// if name == nil{
	// 	this.Redirect("/", 302)
	// 	return
	// }
	*/
	// 从数据库查数据展示在html页面
	// 1.查询
		// orm对象
		o := orm.NewOrm()
		qs := o.QueryTable("Article")
		var articles[] models.Article
		// qs.All(&articles)//等同于 select * from Article

		// pageIndex1 := 1
		pageIndex := this.GetString("pageIndex")
		pageIndex1 , err := strconv.Atoi(pageIndex)
		if err != nil{
			pageIndex1 = 1
		}	
		// 获取总记录数
		//多表查询，会有惰性查询的特性，也就是默认不关联表，需要指定
		// 用RelatedSel("关联的表")函数指定
		//获取传递过来的类型
		var count int64
		typeName := this.GetString("select")
		if typeName != ""{
			count , err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count() //加过滤器，根据类型筛选数据
		}else{
			count , err = qs.RelatedSel("ArticleType").Count() //加过滤器，根据类型筛选数据
		}
		if err != nil{
			beego.Info("count err : ",err)
			return
		}
		// beego.Info(count)
		// 获取总页数
		pageSize := 1
		start := pageSize*(pageIndex1 - 1)//分页起始页
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)//放到articles里

		pageCount := float64(count) / float64(pageSize) //避免 3/2 = 1
		pageCount1 := math.Ceil(pageCount) //向上取整

		//控制上一页不能超出一
		FirstPage := false
		if pageIndex1 == 1{
			FirstPage = true
		}
		LastPage := false
		if pageIndex1 == int(pageCount1){
			LastPage = true
		}

		//显示文章的所有分类
		// 有些数据是不需要时刻从数据库中查找
		// 使用redis作为缓存
		var types[] models.ArticleType

		//从redis数据库中取数据
		conn,err := redis.Dial("tcp",":6379")
		if err !=nil{
			beego.Info("redis连接出错")
			return
		}
		rel , _ := redis.Bytes(conn.Do("get", "types"))
		// if err != nil{
		// 	beego.Info("从redis获取数据失败")
		// 	return
		// }
		//反序列化
		dec := gob.NewDecoder(bytes.NewReader(rel))
		dec.Decode(&types)//解码
		// beego.Info(types)


		if len(types) == 0{
			o.QueryTable("ArticleType").All(&types)//从mysql取数据
			var buffer bytes.Buffer
			// 序列化
			enc := gob.NewEncoder(&buffer)
			enc.Encode(&types)//编码
			_,err = conn.Do("set", "types",buffer.Bytes())
			if err != nil{
				beego.Info("redis数据库操作错误")
				return
			}
			beego.Info("从mysql数据库中取数据")
		}

		//根据文章类型获取数据
		//1. 接收数据
		// typeName := this.GetString("select")
		// 2.处理数据
		var articlewithtype []models.Article
		// typeName := this.GetString("select")
		if typeName == ""{
			// beego.Info("获取下拉框数据失败")
			//默认是显示所有数据
			qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articlewithtype)
		}else{
			qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articlewithtype)
		}
		// beego.Info(typeName)
		//3. 查询数据
		//这里不加 relatedsel 函数会出现惰性查询
		// o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)

		//获取登录的用户名
		// layout.html也可以获取传递过来的数据
		userName := this.GetSession("userName")
		this.Data["userName"] = userName
		this.Data["typeData"] = typeName
		this.Data["types"] = types
		this.Data["FirstPage"] = FirstPage
		this.Data["LastPage"] = LastPage
		this.Data["count"] = count
		this.Data["pageCount"] = pageCount1
		this.Data["pageIndex"] = pageIndex1
		this.Data["articles"] = articlewithtype
	//2.数据在视图上显示 
	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["title"] = "head/Articlehead.html"
	// this.TplName = "article.html"
	this.TplName = "article.html"
}

func (this *ArticleController) ShowAddArticle() {
	//查数据
	var types[] models.ArticleType
	o := orm.NewOrm()
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["title"] = "head/Addhead.html"
	this.TplName = "add.html"
}

//添加文章 post

func (this *ArticleController) HandlerAddArticle() {
	//获取文章名
	artiname := this.GetString("articleName")
	// 获取文章内容
	articontent := this.GetString("textcontent")
	//上传图片
	//字节流 文件的信息（大小等）
	file,head, _ := this.GetFile("uploadfile")
	defer file.Close()
	// 限定文件格式
	// 1、判断文件格式
	//也可不上传文件
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Info("文件格式错误")
		return
	}
	//2、限定文件大小
	// 单位 B
	if head.Size > 500000{
		beego.Info("文件太大，不允许上传")
		return
	}
	//存入项目中的static/img文件夹下
	// 3、不能重名
	// Format参数是规定好的，不能改
	//文件名不能有一些特殊字符：比如：英文 : | 等
	fileName := time.Now().Format("2006-01-02-15-04-05")
	// 					name参数				文件名需要加上后缀
	err := this.SaveToFile("uploadfile","./static/img/"+fileName+ext)
	if err != nil{
		beego.Info(err)
		return
	}
	// beego.Info(fileName)
	// beego.Info(artiname, articontent)
	// 3、插入数据
		//1.获取orm对象
		o := orm.NewOrm()
		// 2.创建一个插入对象
		article := models.Article{}
		// 3.赋值
		article.Title = artiname
		article.Content = articontent
		article.Img = "./static/img/"+fileName+ext

		//给article对象赋值
		// 获取到下拉框传递过来的类型数据
		typeName := this.GetString("select")
		// 类型判断
		if typeName == ""{
			beego.Info("下拉框数据错误")
			return
		}
		// 获取articleType对象
		var artiType models.ArticleType
		artiType.TypeName = typeName
		err = o.Read(&artiType,"TypeName")
		if err != nil{
			beego.Info("获取类型错误")
			return
		}
		article.ArticleType = &artiType

		// 4.插入
		_ , err = o.Insert(&article)
		if err != nil{
			beego.Info("插入失败！")
			return
		}
	// 4.返回视图,文章列表页
	this.Redirect("/Article/ShowArticle", 302)
}

//显示文章详情
func(this *ArticleController)ShowContent(){
	// 1.获取id,这个需要在路由里配置 :id
	pageid := this.Ctx.Input.Param(":id")
	// beego.Info(pageid)
	// 2.查询数据
		// 获取orm对象
		o := orm.NewOrm()
		//获取查询对象
		// 从数据库中查询相应id的文章
		id2,_ := strconv.Atoi(pageid)
		// beego.Info(id2)
		article := models.Article{Id: id2}
		err := o.Read(&article)
		if err != nil{
			beego.Info("查询数据为空")
			return
		}
		// 每点击一次，阅读量加一
		article.Count += 1
		// 根据文章的类型id去文章类型表查相应的类型
		o2 := orm.NewOrm()
		qs := o2.QueryTable("ArticleType")
		var result models.ArticleType
		qs.Filter("Id", article.ArticleType.Id).All(&result)
		// beego.Info(result)
		// 多对多插入读者
		// 1.获取操作对象
		// article := models.Article{Id:id2} //上面已定义
		// 2.获取多对多操作对象
		// 联合哪个表
		m2m := o.QueryM2M(&article, "User")
		// 3.获取插入对象
		// 因为登录，而且之前设置了session，所以可以通过session拿到用户名
		userName := this.GetSession("userName")
		user := models.User{}
		user.UserName = userName.(string)//通过session拿到的是一个空接口，需要用类型断言
		o.Read(&user, "UserName")
		// 4.多对多插入
		_ , err = m2m.Add(&user)
		if err != nil{
			beego.Info("插入失败。。")
			return
		}
		// 更新数据库
		o.Update(&article) //没有指定哪一页更新，会自动查找更新的那一页

		// 多对多查询,查询最近浏览用户，并显示在视图上
		// 这个查询会有很多重复的信息
		// o.LoadRelated(&article, "User")
		var users []models.User
		// o.QueryTable("Article").Filter("Id",id2).Filter("User__User__Id",user.Id).Distinct().One(&article)
		//这个查询需要反向查询,这是beego的多对多查询的设计问题
		// 关键：查询和插入是反着来的
		o.QueryTable("User").Filter("Articles__Article__Id", id2).Distinct().All(&users)

	//3.传递数据给视图
	this.Data["users"] = users
	this.Data["article"] = article
	this.Data["isoType"] = result
	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["title"] = "head/Contenthead.html"
	this.TplName = "ArticleContent.html"
}

// 删除文章
func(this *ArticleController)HandlerDelete(){
	//获取删除文章的id
	pageid := this.Ctx.Input.Param(":id")
	// orm对象
	o := orm.NewOrm()
	id2,_ := strconv.Atoi(pageid)
	article := models.Article{Id:id2}
	//删除文章时，还需要删除对应的图片
	// 先从数据库读出这篇文章的信息
	o.Read(&article)
	// beego.Info(article.Img)
	os.Remove(article.Img)
	o.Delete(&article)
	this.Redirect("/Article/ShowArticle", 302)
}

// 更新文章
//先显示文章内容
func(this *ArticleController)ShowPage(){
	//拿到文章的id
	pageId := this.Ctx.Input.Param(":id")
	//orm对象
	o := orm.NewOrm()
	id2 ,_ := strconv.Atoi(pageId)
	//显示文章的所有分类
	// 有些数据是不需要时刻从数据库中查找
	// 使用redis作为缓存
	var types[] models.ArticleType
	conn,err := redis.Dial("tcp",":6379")
	if err !=nil{
		beego.Info("redis连接出错")
		return
	}
	types = GetTypeFromRedis(conn)
	//查数据库
	article := models.Article{Id:id2}
	o.Read(&article)
	//在页面上显示
	// beego.Info(article.Img)
	this.Data["article"] = article
	this.Data["types"] = types
	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["title"] = "head/Updatehead.html"
	this.TplName = "update.html"
}
//更新文章
//先从数据库读取，再更新
func(this *ArticleController)HandlerUpdate(){
	//拿到文章的id
	//如果使用/xxxx/:id的路由，那么获取id就该这样获取
	pageId ,err := this.GetInt(":id")
	// beego.Info(pageId)
	if err != nil{
		beego.Info("id出错...")
		return
	}
	//获取文章名
	artiname := this.GetString("articleName")
	// 获取文章内容
	articontent := this.GetString("textcontent")
	//orm对象
	//上传图片
	//字节流 文件的信息（大小等）
	file,head, err := this.GetFile("uploadfile")
	if err != nil{
		beego.Info(err)
		return
	}
	defer file.Close()
	// 限定文件格式
	// 1、判断文件格式
	//也可不上传文件
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Info("文件格式错误")
		return
	}
	//2、限定文件大小
	// 单位 B
	if head.Size > 500000{
		beego.Info("文件太大，不允许上传")
		return
	}
	//存入项目中的static/img文件夹下
	// 3、不能重名
	// Format参数是规定好的，不能改
	//文件名不能有一些特殊字符：比如：英文 : | 等
	fileName := time.Now().Format("2006-01-02-15-04-05")
	// 					name参数				文件名需要加上后缀
	err = this.SaveToFile("uploadfile","./static/img/"+fileName+ext)
	if err != nil{
		beego.Info(err)
		return
	}
	// beego.Info(fileName)
	// beego.Info(artiname, articontent)
	// 3、读取数据
		//1.获取orm对象
		o := orm.NewOrm()
		// 2.更新对象
		article := models.Article{Id:pageId}
		// 3.读取
		err = o.Read(&article)
		if err != nil{
			beego.Info("没有这篇文章")
			return
		}
		// beego.Info(article.Img)
		// 删除原来的图片
		os.Remove(article.Img)
		// 3.赋值
		article.Title = artiname
		article.Content = articontent
		article.Img = "./static/img/"+fileName+ext
		// 4.更新
		_ , err = o.Update(&article)
		if err != nil{
			beego.Info("更新失败！")
			return
		}
	// 4.返回视图,文章列表页
	this.Redirect("/Article/ShowArticle", 302)
}

//添加文章分类
func(this *ArticleController)ShowAddType(){
	//查数据库，显示在视图上
	var articles[] models.ArticleType
	o := orm.NewOrm()
	//查询所有
	_,err := o.QueryTable("ArticleType").All(&articles)
	if err != nil{
		beego.Info("查询分类失败")
	}

	this.Data["articles"] = articles
	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["title"] = "head/Typehead.html"
	this.TplName = "addType.html"
}
//添加分类
func(this *ArticleController)HandlerAddType(){
	//获取数据
	typeName := this.GetString("typeName")
	//判断数据
	if typeName == ""{
		beego.Info("类型名不能为空")
		return
	}
	//执行插入操作
	o := orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName = typeName
	_ , err := o.Insert(&artiType)
	if err != nil{
		beego.Info("插入失败")
		return
	}
	this.Redirect("/Article/AddArticleType", 302)
}