package main

import (
	"fmt"
	"net"
	"time"
)

func startCli(port uint16) {

	// 连接到服务器
	tcp, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Printf("连接服务失败:%d\n", port)
		return
	}
	fmt.Printf("已连接到服务器:%d\n", port)
	go func() { //写入数据操作
		for {
			wLen, err := tcp.Write([]uint8("-->这是客服端发送的数据"))
			if err != nil {
				fmt.Printf("客户端端口:%d写入数据失败,wLen=%d  err=%q\n", port, wLen, err)
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		for {
			buffer := make([]uint8, 64*1024)
			len, err := tcp.Read(buffer)
			if len == 0 || err != nil {
				fmt.Printf("客户端口:%d读取数据失败,len=%d  err=%q\n", port, len, err)
				break
			}
			//fmt.Printf("-->得到服务端%d数据:%s\n", port, string(buffer[:len]))
		}
	}()
}
