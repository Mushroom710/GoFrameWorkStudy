package routers

import (
	"hello/controllers"

	"github.com/astaxie/beego"
	// 这个包不能是默认的context包
	"github.com/astaxie/beego/context"
)

func init() {
	// 每一个页面都加session判断，会很难受，因此使用过滤器函数来实现这个功能
	// 这个函数的路由配置需要在最开始
	//正则路由匹配，满足条件就先匹配这个路由
	beego.InsertFilter("/Article/*", beego.BeforeRouter, FilterFun)
    // beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandlerReg")
	// 默认访问登录页面
	beego.Router("/Login", &controllers.LoginController{}, "get:ShowLogin;post:HandlerLogin")
	//文章页面路由
	beego.Router("/Article/ShowArticle", &controllers.ArticleController{}, "get:ShowArticleList;post:HandlerSelect")
	// 添加文章页面
	beego.Router("/Article/AddArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandlerAddArticle")
	//测试路由
	beego.Router("/test",&controllers.TestController{},"get:Show;post:Handler")
	// 显示文章详情的路由
	beego.Router("/Article/ArticleContent/:id",	&controllers.ArticleController{}, "get:ShowContent")
	// beego.Router("/ArticleContent", &controllers.ArticleController{}, "get:ShwoContent")
	//删除文章的路由
	beego.Router("/Article/DeleteArticle/:id",&controllers.ArticleController{},"get:HandlerDelete")
	// 编辑文章的路由
	beego.Router("/Article/UpdateArticle/:id",&controllers.ArticleController{},"get:ShowPage;post:HandlerUpdate")
	//添加文章分类的路由
	beego.Router("/Article/AddArticleType",&controllers.ArticleController{},"get:ShowAddType;post:HandlerAddType")
	// 退出登录的路由
	beego.Router("/Logout",&controllers.UserController{},"get:LogOut")

	// 展示前端页面
	beego.Router("/", &controllers.FrontPageController{}, "get:ShowInfo")
	// 展示文章内容
	beego.Router("/PageArticle/:id",&controllers.FrontPageController{},"get:ShowMore")
	// 展示图片页面
	beego.Router("/ShowPicture", &controllers.FrontPageController{}, "get:ShowPic")
	// 展示关于页面
	beego.Router("/ShowAbout", &controllers.FrontPageController{}, "get:ShowAbout")
}

var FilterFun = func(ctx *context.Context){
	// 获取session
	userName := ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302, "/")
	}
}
