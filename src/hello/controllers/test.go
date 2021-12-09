package controllers

import "github.com/astaxie/beego"

type TestController struct {
	beego.Controller
}

func (this *TestController) Show() {
	this.TplName = "text.html"
}

func (this *TestController) Handler() {
	f, h, _ := this.GetFile("upload")
	defer f.Close()

	err := this.SaveToFile("upload", "./static/img/"+h.Filename)
	if err != nil {
		beego.Info(err)
	}
}
