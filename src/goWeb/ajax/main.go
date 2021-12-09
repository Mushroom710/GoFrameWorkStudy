package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type People struct{
	Name string `json:"name"`
	Age int 	`json:"age"`
}

func GetHandler(w http.ResponseWriter,r *http.Request){
	// 跨域请求的设置这个相应头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("这是get请求"))
}

func PostHandler(w http.ResponseWriter,r *http.Request){
	data := make(map[string]string)
	data["a"] = "hello"
	data["b"] = "world"
	jsondata,_ := json.Marshal(data)
	w.Header().Set("content-type", "text/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsondata)
}

func AjaxHandler(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	username := r.FormValue("admin")
	passwd := r.FormValue("passwd")
	if(username == "admin" && passwd =="123"){
		w.Header().Set("content-type", "text/html")
		str := "<div>hello world</div>"
		w.Write([]byte(str))
	}else{
		fmt.Fprintln(w, "用户名或密码错误")
	}
}

func main(){
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/get", GetHandler)
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/ajax",AjaxHandler)
	server.ListenAndServe()
}