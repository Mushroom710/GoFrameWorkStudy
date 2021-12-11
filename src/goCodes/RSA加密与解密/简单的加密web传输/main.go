package main

import (
	"net/http"

)

/**
	简单的使用了RSA加密技术
	RSA加密：公钥加密，用私钥解密。
*/ 
func main() {
	server := http.Server{
		Addr: "localhost:8080",
	}

	http.HandleFunc("/getPublic", GetPublic)
	http.HandleFunc("/postPrivate", PostPrivate)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)

	server.ListenAndServe()
}
