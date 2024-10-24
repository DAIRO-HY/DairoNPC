package main

import (
	"DairoNPC/constant"
	"DairoNPC/session"
	"os"
	"strings"
)

func main() {
	var args = os.Args
	if len(args) < 2 { //测试用
		args = []string{"-h:127.0.0.1", "-k:njeHds*fs4tfsd", "-t:1781", "-u:1682"}
	}
	for _, it := range args {
		paramArr := strings.Split(it, ":")
		switch paramArr[0] {
		case "-h":
			constant.Host = paramArr[1]
		case "-k":
			constant.Key = paramArr[1]
		case "-t":
			constant.TcpPort = paramArr[1]
		case "-u":
			constant.UdpPort = paramArr[1]
		}
	}
	session.Open()
}