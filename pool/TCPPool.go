package pool

import (
	"DairoNPC/HeaderUtil"
	"DairoNPC/bridge"
	"DairoNPC/constant"
	"net"
	"strconv"
	"strings"
)

/**
 * 等待分配工作的Socket
 */
type TCPPool struct {
	npcTCP net.Conn
}

/**
 * 开始等待分配工作
 */
func (mine *TCPPool) start() {

	//发送客户端信息
	mine.sendClientInfoToServer()

	//等待分配工作
	mine.waitWork()
	removePool(mine)
}

/**
 * 关闭资源
 */
//func (mine *TCPPool) close() {
//    this.socket.close()
//}

/**
 * 发送客户端信息
 */
func (mine *TCPPool) sendClientInfoToServer() {

	//标记这是一个连接池  并且将客户端ID告诉服务器
	err := HeaderUtil.SendFlag(mine.npcTCP, HeaderUtil.SERVER_TCP_POOL_REQUEST, strconv.Itoa(constant.ClientId))
	if err != nil {
		mine.npcTCP.Close()
	}
}

/**
 * 等待分配工作
 */
func (mine *TCPPool) waitWork() {

	//加密类型及目标端口 格式:加密状态|端口  1|80   1|127.0.0.1:80
	//1:加密  0:不加密
	hearder, err := HeaderUtil.GetHeader(mine.npcTCP)
	if err != nil { //服务器端连接达到上限或者连接池被强制关闭
		mine.npcTCP.Close()
		return
	}

	hearders := strings.Split(hearder, "|")

	//加密状态  1:加密  0:不加密
	securityState := hearders[0]

	//目标服务器信息
	targetAddr := hearders[1]
	if !strings.Contains(targetAddr, ":") { //如果包含了ip地址
		targetAddr = "127.0.0.1:" + targetAddr
	}
	bridge.Start(securityState == "1", targetAddr, mine.npcTCP)
}
