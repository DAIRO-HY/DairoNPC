package udp_pool

import (
	"DairoNPC/bridge/udp_bridge"
	"DairoNPC/constant"
	"net"
	"strconv"
	"strings"
)

/**
 * 等待分配工作的Socket
 */
type UDPPool struct {
	Socket *net.UDPConn

	/**
	 * 关闭标记
	 */
	isCloseFlag bool
}

/**
 * 开始等待分配工作
 */
func (mine *UDPPool) start() {
	//发送客户端信息
	mine.sendClientInfo()

	//等待分噢诶工作
	mine.waitWork()
	removePoolList(mine)
}

/**
 * 关闭连接
 */
func (mine *UDPPool) close() {
	mine.isCloseFlag = true
	mine.Socket.Close()
}

/**
 * 发送客户端ID信息
 */
func (mine *UDPPool) sendClientInfo() {
	clientId := strconv.Itoa(constant.ClientId)
	mine.Socket.Write([]byte(clientId))
}

/**
 * 等待分配工作
 */
func (mine *UDPPool) waitWork() {
	headBuf := make([]byte, 1024)
	length, _, err := mine.Socket.ReadFromUDP(headBuf)
	if err != nil {
		mine.Socket.Close()
		return
	}
	//得到头部信息
	head := string(headBuf[length])

	//关闭链接池标识
	if head == "CLOSE" { //TODO:
		mine.Socket.Close()
		return
	}

	//头部信息数组
	headArr := strings.Split(head, "|")

	//加密类型及目标端口 格式:加密状态|端口  1|80   1|127.0.0.1:80
	//1:加密  0:不加密

	//加密状态  1:加密  0:不加密
	isEncodeData := headArr[0] == "1"

	//目标服务器信息
	targetAddr := headArr[1]
	if !strings.Contains(targetAddr, ":") { //如果没包含IP地址
		targetAddr = "127.0.0.1:" + targetAddr
	}

	//由于UDP不可靠协议,这里收到的数据可能不是我们想要的数据
	//println("-->由于UDP不可靠协议,这里收到的数据可能不是我们想要的数据")
	//println("-->:${hearder}")
	//this.socket.close()
	//return
	udp_bridge.Start(isEncodeData, targetAddr, mine.Socket)
}
