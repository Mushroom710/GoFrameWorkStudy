package main

/**
这个文件专门处理加密和解密的数据
*/

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// 处理web：发送公钥给请求的页面
func GetPublic(w http.ResponseWriter, r *http.Request) {
	// 跨域请求得设置这个请求头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 读取公钥的数据
	data, err := ioutil.ReadFile("public.pem")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("请求公钥！")
	// 返回给发出请求的的页面
	w.Write([]byte(data))
}

// 处理web：用私钥对数据进行解密
func PostPrivate(w http.ResponseWriter, r *http.Request) {
	// 跨域请求得设置这个请求头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 拿到加密过的账户和密码
	account := r.PostFormValue("Account")
	passwd := r.PostFormValue("Password")
	// 对加密过的账户和密码先进行base64标准编码，然后再解码
	codeAccount, _ := base64.StdEncoding.DecodeString(account)
	codePasswd, _ := base64.StdEncoding.DecodeString(passwd)

	// 读取私钥文件内容
	file, _ := os.Open("private.pem")
	// 打开文件要记得关闭
	defer file.Close()
	//file.Stat()可以得到文件的一些基本信息。比如文件大小
	info, _ := file.Stat()
	// make一个字节数组，大小为 info.Size()
	buf := make([]byte, info.Size())
	// 把文件内容全部读进这个字节数组
	file.Read(buf)

	//pem解码
	block, _ := pem.Decode(buf)
	// x509解码
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	// 这个时候就可以用私钥解密数据
	accountplainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, codeAccount)
	passwdplainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, codePasswd)
	fmt.Println("账户：" + string(accountplainText))
	fmt.Println("密码：" + string(passwdplainText))
	w.Write([]byte("解析成功！"))
}

// 对来自web的数据进行解密
func DePostForm(account string, passwd string) (accountplainText, passwdplainText []byte) {
	// 对加密过的账户和密码先进行base64标准编码，然后再解码
	codeAccount, _ := base64.StdEncoding.DecodeString(account)
	codePasswd, _ := base64.StdEncoding.DecodeString(passwd)

	// 读取私钥文件内容
	file, _ := os.Open("private.pem")
	// 打开文件要记得关闭
	defer file.Close()
	//file.Stat()可以得到文件的一些基本信息。比如文件大小
	info, _ := file.Stat()
	// make一个字节数组，大小为 info.Size()
	buf := make([]byte, info.Size())
	// 把文件内容全部读进这个字节数组
	file.Read(buf)

	//pem解码
	block, _ := pem.Decode(buf)
	// x509解码
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	// 这个时候就可以用私钥解密数据
	accountplainText, _ = rsa.DecryptPKCS1v15(rand.Reader, privateKey, codeAccount)
	passwdplainText, _ = rsa.DecryptPKCS1v15(rand.Reader, privateKey, codePasswd)

	fmt.Println("账户：" + string(accountplainText))
	fmt.Println("密码：" + string(passwdplainText))

	return
}

// 对从数据库查出的密码进行解密
// 把解密的数据回传
func DecryptionSQL(passwd string) (passwdplainText []byte) {
	// 对密码进行base64编码解码
	codePasswd, _ := base64.StdEncoding.DecodeString(passwd)
	// 读取私钥文件内容
	file, _ := os.Open("private.pem")
	// 打开文件要记得关闭
	defer file.Close()
	//file.Stat()可以得到文件的一些基本信息。比如文件大小
	info, _ := file.Stat()
	// make一个字节数组，大小为 info.Size()
	buf := make([]byte, info.Size())
	// 把文件内容全部读进这个字节数组
	file.Read(buf)

	//pem解码
	block, _ := pem.Decode(buf)
	// x509解码
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	passwdplainText, _ = rsa.DecryptPKCS1v15(rand.Reader, privateKey, codePasswd)
	return passwdplainText
}
