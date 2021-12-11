package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

/**
这个文件与数据库操作有关
*/

// 定义用户结构体
type User struct {
	Id      int
	Account string
	Passwd  string
}

// 处理登录
func Login(w http.ResponseWriter, r *http.Request) {
	// 跨域请求得设置这个请求头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 拿到加密过的账户和密码
	account := r.PostFormValue("Account")
	password := r.PostFormValue("Password")
	// 解密回来的数据是字节切片
	accountplainText, passwdplainText := DePostForm(account, password)
	//建立数据库连接
	Db, err := sql.Open("mysql", "root:12345678@tcp(localhost:3306)/test")
	if err != nil {
		log.Default().Println("数据库连接出错：", err)
		return
	}
	defer Db.Close() //关闭数据库连接
	// 根据账户从数据库中查找对应密码
	row := Db.QueryRow("select passwd from login where account = ?", string(accountplainText))
	if err != nil {
		log.Default().Println("查无此人！", err)
		return
	}
	var passwd string
	// Scan需要传入一个指针值
	err = row.Scan(&passwd)
	if err != nil {
		log.Println("Scan出错：", err)
		return
	}
	// 解密
	Depasswd := DecryptionSQL(passwd)
	fmt.Println(string(Depasswd),string(passwdplainText))
	// 数据库的密码和form表单的密码进行比较
	if string(Depasswd) == string(passwdplainText) {
		log.Println("验证成功！")
		w.Write([]byte("登录成功！"))
	} else {
		log.Panicln("用户名或密码错误！")
		w.Write([]byte("用户名或密码错误！"))
	}
}

// 处理注册
// 账户可以用明文，但是密码要用RSA加密后的
func Register(w http.ResponseWriter, r *http.Request) {
	// 跨域请求得设置这个请求头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 拿到加密过的账户和密码
	account := r.PostFormValue("Account")
	password := r.PostFormValue("Password")
	// 解密回来的数据是字节切片
	accountplainText, _ := DePostForm(account, password)
	// 建立数据库连接
	Db, err := sql.Open("mysql", "root:12345678@tcp(localhost:3306)/test")
	if err != nil {
		log.Default().Println("数据库连接出错：", err)
		return
	}
	defer Db.Close() //关闭数据库连接
	// 账户不能重复
	// 通过查账户数量来确定账户是否已经被注册
	var count int
	row := Db.QueryRow("select count(*) from login where account = ?", string(accountplainText))
	row.Scan(&count)
	fmt.Println("账户数：",count)
	if count >= 1{
		w.Write([]byte("账户已存在！"))
		return
	}
	// 存入数据库。账户是明文，但密码是加密过的
	result, err := Db.Exec("insert into login (account,passwd) values(?,?)", accountplainText, password)
	if err != nil {
		log.Println("插入出错：", err)
		return
	}
	n1, _ := result.RowsAffected()
	log.Println("影响行数：", n1)
	w.Write([]byte("注册成功"))
}
