package main

import (
	"DairoNPC/DairoNPCMain"
	"net"
	"os"
)

func main() {
	net.DialUDP("udp", nil, raddr)
	DairoNPCMain.Open(os.Args)
}
