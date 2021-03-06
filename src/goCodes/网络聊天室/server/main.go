// server
package main


import (
	"fmt"
	"net"
)


var ConnMap map[string]*net.TCPConn


func checkErr(err error) int {
	if err != nil {
		if err.Error() == "EOF" {
			fmt.Println("用户退出")
			return 0
		}
		fmt.Println("发生错误")
		return 0
	}
	return 1
}
func say(tcpConn *net.TCPConn) {
	for {
		data := make([]byte, 256)
		total, err := tcpConn.Read(data)
		if err != nil {
			fmt.Println(string(data[:total]), err)
		} else {
			fmt.Println(string(data[:total]))
		}


		flag := checkErr(err)
		if flag == 0 {
			break
		}
		for _, conn := range ConnMap {
			if conn.RemoteAddr().String() == tcpConn.RemoteAddr().String() {
				continue
			}
			conn.Write(data[:total])
		}
	}
}
func main() {
	//var conn net.TCPConn
	//localAddr :=conn.LocalAddr().String()
	//fmt.Println(localAddr)
	//tcpAddr, _ := net.ResolveTCPAddr("tcp",localAddr)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	tcpListen, _ := net.ListenTCP("tcp", tcpAddr)
	ConnMap = make(map[string]*net.TCPConn)
	fmt.Println("服务器在 127.0.0.1:8080 启动服务")
	for {
		tcpConn, _ := tcpListen.AcceptTCP()
		defer tcpConn.Close()
		ConnMap[tcpConn.RemoteAddr().String()] = tcpConn
		fmt.Println("连接客户端信息：", tcpConn.RemoteAddr().String())


		go say(tcpConn)
	}
}