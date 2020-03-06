package controllers

import (
	"BeenBlog/models"
	"github.com/astaxie/beego"
	"time"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	c.Data["time"] = time.Now()
	op := c.Input().Get("op")
	switch op {
	case "add":
		name := c.Input().Get("name")

		if len(name) ==0{
			break
		}
		err := models.AddCategory(name)
		if err !=nil{
			beego.Error(err)
		}
		c.Redirect("/category",301)
		return
	case "del":
		id := c.Input().Get("id")
		if len(id) ==0{
			break
		}
		err := models.DelCategory(id)
		if err !=nil{
			beego.Error(err)
		}
		c.Redirect("/category",301)
		return
	}
	c.Data["IsCategory"]=true
	var err error
	c.Data["Categories"],err =models.GetAllCategories()
	if err != nil{
		beego.Error(err)
	}
	c.TplName = "category.html"
}