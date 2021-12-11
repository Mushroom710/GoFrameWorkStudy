package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)
// 全局Db，简称小连接池
var (
	Db *sql.DB
	err error
)

type User struct{
	Id string
	Account string
	Passwd string
}

func main(){
	Db,err = sql.Open("mysql", "root:12345678@tcp(localhost:3306)/test")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer Db.Close()
	// 向数据库插入一条数据
	// result,err := Db.Exec("insert into login(account,passwd) values(?,?)", "zhangsan","123456")
	// if err != nil{
	// 	fmt.Println(err)
	// 	return
	// }
	// n1,_:=result.LastInsertId()
	// n2,_:=result.RowsAffected()
	// fmt.Println(n1,n2)
	
	user := User{}
	// 查询数据
	row := Db.QueryRow("select * from login where id = ?",1)
	// Scan要使用地址赋值
	err := row.Scan(&user.Id,&user.Account,&user.Passwd)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(user)

}