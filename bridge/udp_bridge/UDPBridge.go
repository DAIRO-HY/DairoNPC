package udp_bridge

import (
	"DairoNPC/util/SecurityUtil"
	"net"
)

/**
 * UDP桥接
 */
type UDPBridge struct {
	isEncodeData bool //是否加密传输
	npsUdp       *net.UDPConn
	targetUDP    *net.UDPConn

	//标记从目标端口读取数据结束
	isTargetReadFinish bool

	//标记从服务端口读取数据结束
	isNpsReadFinish bool
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
		length, _, err := mine.npsUdp.ReadFromUDP(data)
		if err != nil { //服务器端关闭了
			break
		}

		//数据解密
		if mine.isEncodeData {
			SecurityUtil.Mapping(data, length)
		}

		//从服务端收到数据原样发送给目标服务器
		sendLen, err := mine.targetUDP.Write(data[:length])
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
		length, _, err := mine.targetUDP.ReadFromUDP(data)
		if err != nil { //目标端口关闭了
			break
		}

		//数据解密
		if mine.isEncodeData {
			SecurityUtil.Mapping(data, length)
		}

		//从目标端口收到数据原样发送给NPS服务器
		sendLen, err := mine.npsUdp.Write(data[:length])
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
		mine.npsUdp.Close()
		mine.targetUDP.Close()
		removeBridgeList(mine)
	}
}

/**
 * 关闭资源
 */
func (mine *UDPBridge) close() {
	mine.npsUdp.Close()
	mine.targetUDP.Close()
}
