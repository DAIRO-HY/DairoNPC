package udp_bridge

import (
	"fmt"
	"net"
	"sync"
)

//桥接管理

/**
 * 当前桥接列表
 */
var bridgeMap = map[*UDPBridge]bool{}
var bridgeLock sync.Mutex

func Count() int {
	return len(bridgeMap)
}

// 开始UDP桥接
func Start(isEncodeData bool, targetAddr string, clientSocket *net.UDPConn) {

	// 创建一个 UDP 地址
	targetUDPAddr, _ := net.ResolveUDPAddr("udp", targetAddr)

	// 创建一个 UDP 连接
	targetUDP, err := net.DialUDP("udp", nil, targetUDPAddr)
	if err != nil { //目标服务器端口可能没有启动
		//@TODO: 这里应该通知服务器端关闭桥接,否则服务端收不到通知会一直保持连接
		fmt.Println("Error resolving address:", err)
		return
	}
	bridge := &UDPBridge{
		isEncodeData: isEncodeData,
		targetUDP:    targetUDP,
		npsUdp:       clientSocket,
	}
	bridgeLock.Lock()
	bridgeMap[bridge] = true
	bridgeLock.Unlock()
	go bridge.start()
}

// 从桥接列表移除
func removeBridgeList(bridge *UDPBridge) {
	bridgeLock.Lock()
	delete(bridgeMap, bridge)
	bridgeLock.Unlock()
}

// 清空连接
func ShutdownAll() {
	for bridge, _ := range bridgeMap { //@TODO: 这里调用close之后会不会执行removeBridgeList,待测试
		bridge.close()
	}
}
