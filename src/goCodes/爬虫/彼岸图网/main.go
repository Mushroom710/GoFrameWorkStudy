package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

/**
用于爬取彼岸图网的图片
*/

var dirPath string

// 拿到单个网页的全部内容
func HTTPGet(url string)(data string){
	resp,err := http.Get(url)
	if err != nil{
		fmt.Println("打开网页失败：",resp.Status,err.Error())
		return
	}
	buf := make([]byte, 4096)
	for{
		n,err := resp.Body.Read(buf)
		if n==0{
			break;
		}else if err != nil && err != io.EOF{
			fmt.Println("读取出错：",err)
			return
		}
		// fmt.Println(string(buf[:n]))
		data += string(buf[:n])
	}
	return
}

// 从网页内容里提取出图片的超链接
func GetPicUrl(data string){
	// 正则匹配图片url
	picUrl := regexp.MustCompile(`<img src="(?s:(.*?))"`)
	picUrlSilce := picUrl.FindAllStringSubmatch(data, -1)
	n := len(picUrlSilce)

	for i:=0;i<n-1;i++{
		// fmt.Println(i+1,":",picUrlSilce[i][1])
		pic := "https://pic.netbian.com"+picUrlSilce[i][1]
		picData := GetPicData(pic)
		SavePic(picData, i+1)
	}
}

// 打开图片链接，获取图片内容
func GetPicData(picurl string)(picData string){
	resp,err := http.Get(picurl)
	if err != nil{
		fmt.Println("打开图片链接失败：",err.Error())
		return
	}
	buf := make([]byte, 4096)
	for{
		n,err := resp.Body.Read(buf)
		if n==0{
			break
		}else if err != nil && err != io.EOF{
			fmt.Println("获取图片数据出错：",err)
			return
		}
		picData += string(buf[:n])
	}
	return picData
}

// 保存图片
func SavePic(data string,index int){
	fileName := dirPath+"/"+ strconv.Itoa(index) +"张图.jpg"
	f,err := os.Create(fileName)
	if err !=nil{
		fmt.Println("创建文件失败：",err)
		return
	}
	f.WriteString(data)
	f.Close()
	fmt.Printf("第 %d 张图片下载完成！\n",index)
}

// 工作函数
func Working(url string,index int){
	fmt.Printf("开始爬取第 %d 页！\n",index)
	data := HTTPGet(url)
	// 创建一个文件夹用于保存图片
	dirName := time.Now().Format("2006-01-02-15-04-05")
	dirPath = "D:/爬虫文件/彼岸图网"+dirName
	err := os.MkdirAll(dirPath, 0777)
	if err != nil{
		fmt.Println("文件夹创建失败：",err)
		return
	}else{
		fmt.Println("文件夹创建成功>文件夹路径是：",dirPath)
	}
	GetPicUrl(data)
	fmt.Println("文件保存在：",dirPath)
	fmt.Printf("第 %d 页爬取完成！\n",index)
}

func main(){
	fmt.Print("请输入需要爬取的页面数(最小为 2 ,最大为1250，因为现在网站上只有1250个页面)：")
	var n int
	fmt.Scanln(&n)
	for i:=2;i<=n;i++{
		url := "https://pic.netbian.com/index_"+strconv.Itoa(i)+".html"
		Working(url,i)
	}
	fmt.Println("\n三秒后自动退出！")
	time.Sleep(3*time.Second)
}