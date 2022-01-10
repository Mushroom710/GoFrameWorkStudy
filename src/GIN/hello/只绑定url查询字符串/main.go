package main

/**
 ShouldBindQuery 函数只绑定 url 查询参数而忽略 post 数据。
*/ 

import (
	"log"

	"github.com/gin-gonic/gin"
)


type Person struct{
	Name string `form:"name"`
	Address string `form:"address"`
}

func startPage(c *gin.Context){
	var person Person
	if c.ShouldBindQuery(&person) == nil{
		log.Println("----只能绑定url查询参数-----")
		log.Println(person.Name,person.Address)
	}
	c.String(200,"success")
}

func main(){
	r := gin.Default()

	r.Any("/testing",startPage)
	r.Run(":8080")
}