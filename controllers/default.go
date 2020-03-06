package controllers

import (
	"BeenBlog/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["IsHome"] = true
	c.TplName = "home.html"
	c.Data["IsLogin"] =checkAccount(c.Ctx)

	topics ,err:= models.GetAllTopics(c.Input().Get("cate"),c.Input().Get("label"),true)

	if err != nil{
		beego.Error(err)
	} else{
		c.Data["Topics"] =topics
	}
	categories,_:= models.GetAllCategories()
	if err != nil{
		beego.Error(err)
	}
	c.Data["Categories"]=categories
}
