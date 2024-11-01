package udp_bridge

import (
	"DairoNPC/util/SecurityUtil"
	"fmt"
	"net"
	"sync"
)

/**
 * UDP桥接
 */
type UDPBridge struct {
	isEncodeData  bool
	targetAddr    string
	mClientSocket *net.UDPConn

	mTargetSocket *net.UDPConn

	/**
	 * 关闭标识
	 */
	isCloseFlag bool

	//标记从目标端口读取数据结束
	isTargetReadFinish bool

	//标记从服务端口读取数据结束
	isNpsReadFinish bool

	mCloseCheckLock sync.Mutex
}

/**
 * UDP会话
 */
func mTargetSocket() {
	//var ex: Exception? = null
	//repeat(100) {
	//    try {
	//
	//        //因为端口是随机占用,所以这里有端口被占用的可能
	//        return@lazy DatagramSocket()
	//    } catch (e: BindException) {
	//        ex = e
	//    }
	//}
	//println("-->创建UDP失败")
	//throw RuntimeException(ex)
}

/**
 * 目标端IP地址
 */
//private val mTargetInetAddress: InetAddress by lazy {
//    InetManager.get(this.targetIp)
//}

/**
 * 服务器端IP地址
 */
//private val mClsInetAddress: InetAddress by lazy {
//    InetManager.get(CLClient.host)
//}

/**
 * 开始传输数据
 */
func (mine *UDPBridge) start() {

	// 创建一个 UDP 地址
	serverAddr, err := net.ResolveUDPAddr("udp", mine.targetAddr)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// 创建一个 UDP 连接
	mine.mTargetSocket, err = net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	go mine.receiveByCLSServerSendToTarget()
	mine.receiveByTargetSendToCLSServer()
}

/**
 * 从内网穿透服务器接收数据,发送到目标端口
 */
func (mine *UDPBridge) receiveByCLSServerSendToTarget() {

	//经过(MAC)实际测试,UDP每次数据最大在9KB左右(MAC为9126B),所以缓存设置10KB缓存已经足够
	data := make([]byte, 10*1024)
	for {

		//从NPS服务器端接收数据
		length, _, err := mine.mClientSocket.ReadFromUDP(data)
		if err != nil { //服务器端关闭了
			break
		}

		//数据解密
		if mine.isEncodeData {
			SecurityUtil.Mapping(data, length)
		}

		//从服务端收到数据原样发送给目标服务器
		sendLen, err := mine.mTargetSocket.Write(data[:length])
		if err != nil { //目标端口可能已经关闭
			break
		}
		if sendLen < length { //这种情况,应该是出问题,目前还没有遇到过这种情况
			break
		}
	}
	mine.isNpsReadFinish = true
	mine.checkClose()
}

/**
 * 从目标端口接收到数据,发送到内网穿透服务器
 */
func (mine *UDPBridge) receiveByTargetSendToCLSServer() {

	//经过(MAC)实际测试,UDP每次数据最大在9KB左右(MAC为9126B),所以缓存设置10KB缓存已经足够
	data := make([]byte, 10*1024)
	for {

		//从目标端口接收数据
		length, _, err := mine.mTargetSocket.ReadFromUDP(data)
		if err != nil { //目标端口关闭了
			break
		}

		//数据解密
		if mine.isEncodeData {
			SecurityUtil.Mapping(data, length)
		}

		//从目标端口收到数据原样发送给NPS服务器
		sendLen, err := mine.mClientSocket.Write(data[:length])
		if err != nil { //服务端可能已经关闭
			break
		}
		if sendLen < length { //这种情况,应该是出问题,目前还没有遇到过这种情况
			break
		}
	}
	mine.isTargetReadFinish = true
	mine.checkClose()
}

/**
 * 检查资源是否都关闭了
 */
func (mine *UDPBridge) checkClose() {
	if mine.isTargetReadFinish && mine.isNpsReadFinish {
		mine.isCloseFlag = true
		mine.mClientSocket.Close()
		mine.mTargetSocket.Close()
	}
}

/**
 * 关闭资源
 */
func (mine *UDPBridge) close() {
	mine.isCloseFlag = true
	mine.mClientSocket.Close()
	mine.mTargetSocket.Close()
}
