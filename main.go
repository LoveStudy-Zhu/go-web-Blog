package main

import (
	"BeenBlog/controllers"
	"BeenBlog/models"
	_ "BeenBlog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
)

func init()  {
	models.RegisterDB()
}
func main() {
	orm.Debug = true
	orm.RunSyncdb("default",false,true)//强制建表
	beego.Router("/",&controllers.MainController{})
	beego.Router("/login",&controllers.LoginController{})
	beego.Router("/category",&controllers.CategoryController{})
	beego.Router("/topic",&controllers.TopicController{})
	beego.Router("/reply/add",&controllers.ReplyController{},"post:Add")
	beego.Router("/reply/delete",&controllers.ReplyController{},"get:Delete")
	beego.AutoRouter(&controllers.TopicController{})		//自动路由

	//创建附件目录
	os.Mkdir("attachment",os.ModePerm)
	//作为静态文件
	//beego.SetStaticPath("/attachment","attachment")
	//作为一个单独的控制器处理

	beego.Router("/attachment/:all",&controllers.AttachController{})

	beego.Run()
}

