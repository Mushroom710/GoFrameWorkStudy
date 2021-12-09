package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func headers(w http.ResponseWriter, r *http.Request) {
	h := r.Header["Accept-Encoding"]
	h1 := r.Header.Get("Accept-Encoding")
	fmt.Fprintln(w, h, h1)
}

func body(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Fprintln(w, string(body))
}

func process(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	// 第一种获取上传文件的方式
	// r.ParseMultipartForm(1024)
	// fileHeader := r.MultipartForm.File["file"][0]
	// file,err := fileHeader.Open()
	// if err == nil{
	// 	data,err := ioutil.ReadAll(file)
	// 	if err == nil{
	// 		fmt.Fprintln(w, string(data))
	// 	}
	// }
	// 第二种获取上传文件的方式
	// 并保存上传的文件
	// fileHeader.Header带有上传的文件名和后缀
	file, fileHeader, err := r.FormFile("file")
	if err == nil {
		fmt.Println(fileHeader.Filename, fileHeader.Header, fileHeader.Size)
		f, _ := os.OpenFile(fileHeader.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		defer f.Close()
		pic := make([]byte, 4096)
		picdata := ""
		// 循环读取
		for{
			n, err := file.Read(pic)
			if n == 0 {
				fmt.Fprintln(w, "<h1>上传成功！</h1>")
				break
			}else if err != io.EOF && err != nil{
				return
			}
			// f.Write(pic[:n])//直接写字节数组也可以
			picdata += string(pic[:n])//集中放在一个字符串中，后续一次性写入也可以
		}
		f.WriteString(picdata)
	}
	defer file.Close()
}

func main() {
	// 配置服务器
	server := http.Server{
		Addr: "localhost:8080",
	}

	http.HandleFunc("/headers", headers)
	http.HandleFunc("/body", body)
	http.HandleFunc("/process", process)

	server.ListenAndServe()
}
