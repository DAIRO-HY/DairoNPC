package main

import (
	"fmt"
	"net"
)

func startListener(port uint16) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("端口:%d监听失败\n", port)
		return
	}
	for {
		tcp, err := listener.Accept()
		if err != nil {
			fmt.Printf("端口:%d等待客户端连接失败\n", port)
			break
		}
		go receiveData(tcp, port)
	}
}

func receiveData(tcp net.Conn, port uint16) {
	buffer := make([]uint8, 64*1024)
	for {
		len, err := tcp.Read(buffer)
		if len == 0 || err != nil {
			//fmt.Printf("端口:%d读取数据失败,len=%d  err=%q\n", port, len, err)
			break
		}
		wLen, err := tcp.Write([]uint8("这条数据来至服务端，服务端的数据比客服端的多"))
		if err != nil {
			fmt.Printf("端口:%d写入数据失败,wLen=%d  err=%q\n", port, wLen, err)
			break
		}
	}
}
