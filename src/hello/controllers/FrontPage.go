package controllers

import (
	"bytes"
	"encoding/gob"
	"hello/models"
	"math"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
)

type FrontPageController struct{
	beego.Controller
}

// 分离出来的一个用于从redis中取常用数据的函数
func GetTypeFromRedis(conn redis.Conn)(types[] models.ArticleType){
	//从redis数据库中取数据
		// conn,err := redis.Dial("tcp",":6379")
		// if err !=nil{
		// 	beego.Info("redis连接出错")
		// 	return
		// }
		rel , _ := redis.Bytes(conn.Do("get", "types"))
		// if err != nil{
		// 	beego.Info("从redis获取数据失败")
		// 	return
		// }
		//反序列化
		dec := gob.NewDecoder(bytes.NewReader(rel))
		dec.Decode(&types)//解码
		// beego.Info(types)
		return types
}


func (this *FrontPageController)ShowInfo(){
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
		// beego.Info(typeName)
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
		pageSize := count
		// start := pageSize*((int64)(pageIndex1 - 1))//分页起始页
		// qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)//放到articles里
		qs.RelatedSel("ArticleType").All(&articles)
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
		conn,err := redis.Dial("tcp",":6379")
		if err !=nil{
			beego.Info("redis连接出错")
			return
		}
		types = GetTypeFromRedis(conn)

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
			qs.RelatedSel("ArticleType").All(&articlewithtype)
		}else{
			qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articlewithtype)
		}
		// beego.Info(typeName)
		// beego.Info(articlewithtype)
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
	this.Layout = "FrontPage/layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["title"] = "FrontPage/Head/mainHead.html"
	this.LayoutSections["sidebar"] = "FrontPage/sidebar.html"
	this.TplName = "FrontPage/main.html"
}

// 展示更多的文章内容

func (this *FrontPageController)ShowMore(){
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

		// 多对多插入读者
		// 1.获取操作对象
		// article := models.Article{Id:id2} //上面已定义
		// 2.获取多对多操作对象
		// 联合哪个表
		// m2m := o.QueryM2M(&article, "User")
		// 3.获取插入对象
		// 因为登录，而且之前设置了session，所以可以通过session拿到用户名
		// userName := this.GetSession("userName")
		// user := models.User{}
		// user.UserName = userName.(string)//通过session拿到的是一个空接口，需要用类型断言
		// o.Read(&user, "UserName")
		// // 4.多对多插入
		// _ , err = m2m.Add(&user)
		// if err != nil{
		// 	beego.Info("插入失败。。")
		// 	return
		// }
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
	var types[] models.ArticleType
	conn,err := redis.Dial("tcp",":6379")
	if err !=nil{
		beego.Info("redis连接出错")
		return
	}
	types = GetTypeFromRedis(conn)
	// beego.Info(types)
	defer conn.Close()
	//3.传递数据给视图
	// beego.Info(article)
	this.Data["users"] = users
	this.Data["article"] = article
	this.Data["types"] = types
	// this.Layout = "FrontPage/layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["title"] = "head/Contenthead.html"
	this.LayoutSections["sidebar"] = "FrontPage/sidebarnostar.html"
	this.TplName = "FrontPage/ArticleContent.html"
}

// 展示图片页面
func(this *FrontPageController)ShowPic(){
	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var pictures []models.Article
	qs.All(&pictures, "Img")
	// beego.Info(pictures)
	// 从redis数据库中取出常用数据
	var types[] models.ArticleType
	conn,err := redis.Dial("tcp",":6379")
	if err !=nil{
		beego.Info("redis连接出错")
		return
	}
	types = GetTypeFromRedis(conn)
	// beego.Info(types)
	defer conn.Close()
	this.Data["pictures"] = pictures
	this.Data["types"] = types
	// this.Layout = "FrontPage/layout.html"
	this.LayoutSections = make(map[string]string)
	// this.LayoutSections["title"] = "FrontPage/Head/pictureHead.html"
	this.LayoutSections["sidebar"] = "FrontPage/sidebarnostar.html"
	// this.TplName = "FrontPage/Head/pictureHead.html"
	this.TplName = "FrontPage/picturePage.html"
}

// 展示关于页面
func(this *FrontPageController)ShowAbout(){
	// 从redis数据库中取出常用数据
	var types[] models.ArticleType
	conn,err := redis.Dial("tcp",":6379")
	if err !=nil{
		beego.Info("redis连接出错")
		return
	}
	types = GetTypeFromRedis(conn)
	// beego.Info(types)
	defer conn.Close()
	this.Data["types"] = types
	// this.Layout = "FrontPage/layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["sidebar"] = "FrontPage/sidebarnostar.html"
	this.TplName = "FrontPage/About.html"
}