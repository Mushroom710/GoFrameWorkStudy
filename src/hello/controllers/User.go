package controllers

import (
	"fmt"
	"hello/models"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct{
	beego.Controller
}

// 注册业务
func(this *UserController)ShowReg(){
	this.TplName = "register.html"
}

func(this *UserController)HandlerReg(){
	//1.拿到浏览器传递的数据
	name := this.GetString("userName")
	passwd := this.GetString("password");
	// 2.数据处理
	if name == "" || passwd == ""{
		fmt.Println("用户名或密码不能为空！")
		this.TplName = "register.html"
		return
	}
	// 3.插入数据库
	// fmt.Println(name,passwd)

	// a.获取orm对象
	o := orm.NewOrm()
	// b.获取插入对象
	user := models.User{}
	// c.插入操作
	user.UserName = name
	user.Password = passwd

	_ ,err := o.Insert(&user)
	
	if err != nil{
		beego.Info("插入数据失败")
	} 
	// 4.返回登录页面
	// this.Ctx.WriteString("注册成功")
	// 这个时候需要重定向到登录页面
	// 而不是直接渲染登录页面
	this.Redirect("/Login", 302)
	// this.TplName = "login.html"
}

// 登录业务处理
type LoginController struct{
	beego.Controller
}

func (this *LoginController)ShowLogin(){
	name := this.Ctx.GetCookie("userName")
	if name != ""{
		this.Data["userName"] = name
		this.Data["check"] = "checked"
	}
	this.TplName = "login.html"
}

// 1.拿到浏览器的数据
//2.数据处理
// 3.查找数据库
// 4.返回视图
func(this *LoginController)HandlerLogin(){
	// 1.拿到浏览器的数据
	// name := this.GetString("userName")
	// passwd := this.GetString("password")
	// 通过ajax发出的post请求拿数据
	name := this.Ctx.Request.FormValue("UserName")
	passwd := this.Ctx.Request.FormValue("Password")
	// beego.Info(name,passwd)

	// 2.处理数据
	// if name == "" || passwd == ""{
	// 	beego.Info("用户名或密码不能为空")
	// 	this.TplName = "login.html"
	// 	return
	// }
	// 3.查找数据库
		// 1.获取orm对象
		o := orm.NewOrm();
		// 2.获取查询对象
		user := models.User{}
		// 3.查询
		user.UserName = name
		err := o.Read(&user, "UserName")
		if err != nil{
			beego.Info("用户名错误")
			this.TplName = "login.html"
			return
		}
		
		if user.Password != passwd{
			beego.Info("密码错误")
			this.TplName = "login.html"
			return
		}
	// 记住用户名的功能
	check := this.GetString("remeber")
	if check == "on"{
		//				input框的name属性 获取的值 设置时间
		this.Ctx.SetCookie("userName",name,time.Second*3600)
	}else{
		//不设置cookie
		this.Ctx.SetCookie("userName", "username", -1)
	}
	// beego.Info(check)
	// session还需要在conf文件中开启
	// 设置session，用于登录判断。
	this.SetSession("userName", name)

	// 4.返回视图
	// this.Ctx.WriteString("登录成功。。")
	this.Redirect("/Article/ShowArticle", 302)
}

// 退出登录，这时需要删除session
func(this *UserController)LogOut(){
	// 退出登录，删除session
	this.DelSession("userName")
	// 跳转到登录页面
	this.Redirect("/Login", 302)
}