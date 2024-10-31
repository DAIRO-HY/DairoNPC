package tcp_bridge

import (
	"DairoNPC/util/SecurityUtil"
	"DairoNPC/util/TcpUtil"
	"net"
)

type TCPBridge struct {

	//数据是否加密
	isEncodeData bool

	//NPC客户端的TCP
	NpcTCP net.Conn

	//目标服务器连接
	TargetTCP net.Conn

	//Npc客户端的读操作关闭操作
	isNpcReadClosed bool

	// 目标服务连接入方向是否被关闭
	isTargetReadClosed bool
}

// 读取缓存大小(最好和服务器端保持一致)
const READ_CACHE_SIZE = 32 * 1024

// TCP桥接通信开始
func (mine *TCPBridge) start(targetAddr string) { // 连接到服务器

	//与目标端口建立连接
	tcp, err := net.Dial("tcp", targetAddr)
	if err != nil {
		mine.NpcTCP.Close()
		return
	}
	mine.TargetTCP = tcp
	go mine.receiveByNpsSendToTarget()
	mine.receiveByTargetSendToNps()
}

// 从内网穿透服务器接收数据,发送到目标端口
func (mine *TCPBridge) receiveByNpsSendToTarget() {
	data := make([]uint8, READ_CACHE_SIZE)
	for {
		length, err := mine.NpcTCP.Read(data)
		if err != nil {
			break
		}

		//数据解密
		if mine.isEncodeData {
			SecurityUtil.Mapping(data, length)
		}

		//从代理端读取到的数据立即发送目标端
		err = TcpUtil.WriteAll(mine.TargetTCP, data[:length])
		if err != nil {
			break
		}
	}

	//关闭代理端的读操作
	mine.NpcTCP.(*net.TCPConn).CloseRead()

	//关闭目标端的写操作
	mine.TargetTCP.(*net.TCPConn).CloseWrite()

	//标记代理端读操作已经关闭
	mine.isNpcReadClosed = true
	mine.recycle()
}

// 从目标端口接收到数据,发送到内网穿透服务器
func (mine *TCPBridge) receiveByTargetSendToNps() {
	data := make([]uint8, READ_CACHE_SIZE)
	for {
		length, err := mine.TargetTCP.Read(data)
		if err != nil {
			break
		}

		//数据解密
		if mine.isEncodeData {
			SecurityUtil.Mapping(data, length)
		}

		//往NPS服务器发送数据
		err = TcpUtil.WriteAll(mine.NpcTCP, data[:length])
		if err != nil {
			break
		}
	}

	//关闭目标端的读操作
	mine.TargetTCP.(*net.TCPConn).CloseRead()

	//关闭NPS服务端的写操作
	mine.NpcTCP.(*net.TCPConn).CloseWrite()

	//标记目标端读操作已经关闭
	mine.isTargetReadClosed = true
	mine.recycle()
}

// 回收连接
func (mine *TCPBridge) recycle() {
	if mine.isNpcReadClosed && mine.isTargetReadClosed {
		mine.NpcTCP.Close()
		mine.TargetTCP.Close()
		removeBridge(mine)
	}
}

// 关闭链接
func (mine *TCPBridge) shutdown() {
	mine.NpcTCP.Close()
	mine.TargetTCP.Close()
}
