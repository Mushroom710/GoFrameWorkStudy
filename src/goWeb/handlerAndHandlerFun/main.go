package main

import (
	"fmt"
	"net/http"
)

type Myhandler struct{}

func(h *Myhandler) ServeHTTP(w http.ResponseWriter,r *http.Request){
	fmt.Fprintln(w, "hello,world")
}
func hand(w http.ResponseWriter,r *http.Request){
	fmt.Fprintln(w, "are you ok")
}


func main(){
	handler := Myhandler{}
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.Handle("/", &handler)
	http.HandleFunc("/are", hand)
	server.ListenAndServe()

}