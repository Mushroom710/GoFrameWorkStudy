package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 表设计
type User struct{
	Id int
	UserName string
	Password string
	Articles[] *Article `orm:"rel(m2m)"`//用户和文章是多对多关系
}

// 文章表和文章类型表是一对多关系
// 一对多关系 rel 和 reverse成对出现
type Article struct{
	Id int `orm:"pk;auto"`
	Title string `orm:"size(20)"`  //文章标题
	Content string `orm:"type(text)"` //内容
	Img string `orm:"size(50);null"` //图片，存放路径
	Time time.Time `orm:"type(datetime);auto_now_add"` //发布时间
	Count int `orm:"default(0)"` // 阅读量
	ModifyTime time.Time `orm:"type(datetime);auto_now"` //修改时间
	ArticleType *ArticleType `orm:"rel(fk)"`
	User[] *User `orm:"reverse(many)"`	//读者

}
//文章分类表
type ArticleType struct{
	Id int 
	TypeName string `orm:"size(20)"`
	Articles[] *Article `orm:"reverse(many)"` //阅读的文章
}

// 测试表
type Test struct{
	Id int `orm:"pk;auto"`
	Name string `orm:"size(10)"`
	ModifyTime time.Time `orm:"type(datetime);auto_now"`
	Updated time.Time 	`orm:"type(datetime);auto_now"`
}

func init(){
	//时间会相差8小时，在连接数据库是加上 &loc=Local即可
	orm.RegisterDataBase("default", "mysql", "root:12345678@tcp(127.0.0.1:3306)/newsWeb?charset=utf8&loc=Local")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	// 第二个参数会更新表结构，同时会删除之前的所有数据
	orm.RunSyncdb("default", false, true)
}