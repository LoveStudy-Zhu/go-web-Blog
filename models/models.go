package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/unknwon/com"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)
const (
	_DB_NAME = "data/beeblog.db"
	_SQLITE_DRIVER ="sqlite3"
)

func RegisterDB(){
	if !com.IsExist(_DB_NAME){
		os.MkdirAll(path.Dir(_DB_NAME),os.ModePerm)
		os.Create(_DB_NAME)
	}
	orm.RegisterModel(new(Category),new(Topic),new(Comment))
	orm.RegisterDriver(_SQLITE_DRIVER,orm.DRSqlite)
	orm.RegisterDataBase("default",_SQLITE_DRIVER,_DB_NAME,10)
}

type Category struct {
	Id int64
	Title string	`orm:"null"`
	Created time.Time `orm:"index;null"`
	Views int64 `orm:"index"`
	TopicTime time.Time `orm:"index;null"`
	TopicCount int64	`orm:"null"`
	TopicLastUserId int64	`orm:"null"`
}

type  Topic struct {
	Id int64
	Uid int64			//谁写的
	Title string
	Content string `orm:"size(5000)"1`		//内容
	Attachment string						//附件
	Created time.Time	`orm:"index"`	//创建时间
	Updated time.Time	`orm:"index"`	//更新时间
	Views int64		`orm:"index"`	//浏览次数
	Author string	`orm:"index"`	//作者
	ReplyTime time.Time		//回复时间
	ReplyCount	int64		//回复次数
	ReplyLastUserId int64 	//最后回复着用户id
	Labels string			//标签
	Category string
}

//评论
type Comment struct {
	Id int64
	Tid int64
	Name string
	Content string `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func AddTopic(title , label,content, category,attachment string) error {

	//处理标签
	label = "$" +strings.Join(
		strings.Split(label," "),"#$") + "#"
	//空格作为多个标签的分割符


	o := orm.NewOrm()
	topic := &Topic{
		Title:    title,
		Attachment:attachment,
		Content:  content,
		Created:  time.Now(),
		Updated:  time.Now(),
		Category: category,
		Labels:   label,
	}
	_ ,err := o.Insert(topic)
	if err != nil {
		return err
	}
	//更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title",category).One(cate)
	if err ==nil{
		//如果不存在，忽略更新操作
		cate.TopicCount++
		_,err =o.Update(cate)
	}
	return err
}

func AddCategory(name string) error{
	o := orm.NewOrm()
	cate := &Category{Title:name}
	qs := o.QueryTable("category")
	err := qs.Filter("title",name).One(cate)
	if err == nil{
		return err
	}
	_,err =o.Insert(cate)
	if err !=nil{
		return err
	}
	return nil
}
func GetAllTopics(cate ,label string,isDesc bool)([]*Topic,error){
	o :=orm.NewOrm()
	topic := make([]*Topic,0)
	qs := o.QueryTable("topic")
	var err error
	if isDesc{
		if len(cate)>0{
			qs =qs.Filter("category",cate)
		}
		if len(label)>0{
			qs = qs.Filter("labels__contains","$"+label+"#")
		}
		_,err = qs.OrderBy("-Created").All(&topic)
	}else {
		_,err = qs.All(&topic)
	}
	return topic,err
}
func GetAllCategories()([]*Category ,error){
	o := orm.NewOrm()
	cates := make([]*Category,0)
	qs := o.QueryTable("category")
	_,err := qs.All(&cates)
	return cates,err
}

func GetAllReplies(tid string)(replies []*Comment,err error){
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil{
		return nil,err
	}
	replies = make([]*Comment,0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_,err = qs.Filter("tid",tidNum).All(&replies)
	return replies,err
}


func DelCategory(id string)error{
	cid,err := strconv.ParseInt(id,10,64)
	if err != nil{
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id:cid}
	_,err = o.Delete(cate)
	return err
}
func GetTopic(tid string)(*Topic,error) {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil{
		return  nil,err
	}
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id",tidNum).One(topic)
	if err != nil{
		return nil,err
	}
	topic.Views++
	_,err = o.Update(topic)
	topic.Labels=strings.Replace(strings.Replace(
		topic.Labels,"#"," ",-1),"$","",-1)
	return topic,err
}
func ModifyTopic(tid ,label , title , category , content,attachment string) error {
	//处理标签
	label= "$" +strings.Join(
		strings.Split(label," "),"#$") + "#"
	//空格作为多个标签的分割符


	tidNum ,err := strconv.ParseInt(tid,10,6)
	if err != nil{
		return err
	}
	var oldCate,oldAttach string
	o :=orm.NewOrm()
	topic := &Topic{Id:tidNum}
	if o.Read(topic) ==nil{
		oldCate = topic.Category
		topic.Category = category
		topic.Labels =label
		topic.Title =title
		topic.Content =content
		topic.Updated =time.Now()
		topic.Attachment=attachment
		_,err = o.Update(topic)
		if err != nil{
			return err
		}
	}
	//更新分类统计
	if len(oldCate) >0{
		cate:=new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title",oldCate).One(cate)
		if err != nil{
			cate.TopicCount--
			_,err = o.Update(cate)
		}
	}
	//删除旧的附件
	if len(oldAttach) >0{
		//os.Remove(path.Join("attachment",oldAttach))
	}
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title",category).One(cate)
	if err == nil{
		cate.TopicCount++
		_,err =o.Update(cate)
	}


	return err
}
func DeleteTopic(tid string) error {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err !=nil{
		return err
	}
	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id:tidNum}
	if o.Read(topic) ==nil{
		oldCate = topic.Category
		_,err = o.Delete(topic)
		if err !=nil{
			return err
		}
	}
	if len(oldCate)>0{
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title",oldCate).One(cate)
		if err ==nil{
			cate.TopicCount--
			_,err = o.Update(cate)
		}
	}
	return  err
}
func DeleteReply(rid string)error{
	ridNum,err := strconv.ParseInt(rid,10,64)
	if err  !=nil{
		return err
	}
	o := orm.NewOrm()
	var tidNum int64
	reply := &Comment{Id:ridNum}
	if o.Read(reply)==nil{
		tidNum=reply.Tid
		_,err = o.Delete(reply)
		if err != nil{
			return err
		}
	}
	replies := make([]*Comment,0)
	qs := o.QueryTable("comment")
	_,err =qs.Filter("tid",tidNum).OrderBy("-created").All(&replies)
	if err != nil{
		return err
	}
	topic := &Topic{Id:tidNum}
	if o.Read(topic)==nil{
		if len(replies) == 0{
			topic.ReplyTime=time.Time{}
			topic.ReplyCount=0
			_,err = o.Update(topic)
			return nil
		}
		topic.ReplyTime =replies[0].Created
		topic.ReplyCount = int64(len(replies))
		_,err = o.Update(topic)
	}
	return err
}
func AddReply(tid,nickname,content string) error{
	tidNum ,err := strconv.ParseInt(tid,10,64)
	if err != nil{
		return err
	}
	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	o := orm.NewOrm()
	_,err = o.Insert(reply)
	if err !=nil {
		return err
	}
	topic := &Topic{Id:tidNum}
	if o.Read(topic)==nil{
		topic.ReplyTime =time.Now()
		topic.ReplyCount++
		_,err = o.Update(topic)
	}
	return err
}