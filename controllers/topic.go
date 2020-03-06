package controllers

import (
	"BeenBlog/models"
	"github.com/astaxie/beego"
	"path"
	"strings"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get(){
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["IsTopic"]=true
	c.TplName = "topic.html"
	topics ,err:= models.GetAllTopics("","",false)
	if err != nil{
		beego.Error(err)
	} else{
		c.Data["Topics"] =topics
	}
}
func (c *TopicController)Post(){
	if !checkAccount(c.Ctx){
		c.Redirect("/login",302)
		return
	}
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	tid := c.Input().Get("tid")
	category := c.Input().Get("category")
	label := c.Input().Get("label")


	//获得附件
	_,fh,err := c.GetFile("attachment")
	if err !=nil{
		beego.Error(err)
	}
	var attachment string
	if fh != nil{
		//保存附件
		attachment= fh.Filename
		beego.Info(attachment)
		err = c.SaveToFile("attachment",path.Join("attachment",attachment))
		if err !=nil{
			beego.Error(err)
		}
	}
	if len(tid)==0{
		err =models.AddTopic(title,label,content,category,attachment)
	} else{
		err =models.ModifyTopic(tid,label,title,category,content,attachment)
	}
	if err != nil{
		beego.Error(err)
	}
	c.Redirect("/topic",302)
}
func (c *TopicController) Add(){
	c.TplName = "topic_add.html"

}
func (c *TopicController)View(){
	c.TplName = "topic_view.html"
	topic,err :=models.GetTopic(c.Ctx.Input.Param("0"))
	if err != nil{
		beego.Error(err)
		c.Redirect("/",302)
	}
	c.Data["Topic"] = topic
	c.Data["Tid"]= c.Ctx.Input.Param("0")
	c.Data["Labels"] = strings.Split(topic.Labels," ")
	replies,err := models.GetAllReplies(c.Ctx.Input.Param("0"))
	if err !=nil{
		beego.Error(err)
		return
	}
	c.Data["Replies"] = replies
	c.Data["IsLogin"] = checkAccount(c.Ctx)
}
func (c *TopicController)Modify(){
	c.TplName="topic_modify.html"
	tid := c.Input().Get("tid")
	topic ,err := models.GetTopic(tid)
	if err !=nil{
		beego.Error(err)
		c.Redirect("/",302)
		return
	}
	c.Data["Topic"] =topic
	c.Data["Tid"]=tid

}
func (c *TopicController)Delete(){
	if !checkAccount(c.Ctx){
		c.Redirect("/login",302)
		return
	}
	err := models.DeleteTopic(c.Input().Get("tid"))
	if err !=nil{
		beego.Error(err)
	}
	c.Redirect("/topic",302)
}