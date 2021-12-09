package main

import (
	_ "hello/routers"
	"github.com/astaxie/beego"
	_ "hello/models"
)

func main() {
	//视图函数需要映射才能被调用
	//上一页 | 下一页
	beego.AddFuncMap("ShowPrePage", HandlerPrePage)
	beego.AddFuncMap("ShowNextPage",HandlerNextPage)


	beego.Run()
}

//视图函数，用于处理视图的数据
//显示上一页
func HandlerPrePage(data int)(int){
	pageIndex := data - 1
	return pageIndex
}
//显示下一页
func HandlerNextPage(data int)(int){
	pageIndex := data + 1
	return pageIndex
}