package main

import (
	// "bytes"
	//"encoding/gob"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main(){
	// 连接redis
	conn,err := redis.Dial("tcp", ":6379")
	if err !=nil{
		fmt.Println(err)
		return
	}
	defer conn.Close()
	// 直接执行一条命令
	conn.Do("set", "c1","hello")
	// 发送一条命令到缓冲区，此时还没有执行
	conn.Send("get", "c1")//发送一条获取命令
	conn.Send("set","a","abcd")

	// 执行缓冲区的命令
	conn.Flush()
	// 获取数据
	//conn.Receive()返回的rel是一串字节码
	// 需要使用强转函数
	//rel ,err := redis.String(conn.Receive())

	//直接发送命令到redis
	// 获取单个用Sring
	// rel,err := redis.String(conn.Do("get", "c1"))
	// 获取多个用Stings，返回一个string切片
	// rel,err := redis.Strings(conn.Do("mget","name","age","class"))
	// 不想返回一个string切片，而是得到对应的数据类型
	// redis.Values返回一个interface{}切片
	// interface{}可以转为任意类型
	rel,err := redis.Values(conn.Do("mget","name","age","class"))
	var name string
	var age int
	var class string
	// scan扫描器可以可以把相应的数据转为其他数据类型
	// 但是scan只能获取常见类型，不能获取自定义类型
	// 要解析自定义类型，需要通过序列化和反序列
	// 例如：string类型的数字可以转为int类型
	redis.Scan(rel,&name,&age,&class)
	if err !=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(name,age,class)
	fmt.Printf("age type is:%T\n", age)
	conn.Do("set","a","张三")
	// 序列化和反序列化，关键五步代码
	//在我的hello项目中有使用
	// 序列化（字节化）
	// var buffer bytes.Buffer//容器
	// enc := gob.NewDecoder(&buffer)//编码器
	// err = enc.Encode(dest)//编码
	// // 反序列化（反字节化）
	// dec := gob.NewDecoder(bytes.NewReader(buffer byte[]))//解码器
	// dec.Decode(&src)//解码
}