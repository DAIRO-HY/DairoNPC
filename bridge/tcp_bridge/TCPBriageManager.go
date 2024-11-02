package tcp_bridge

import (
	"net"
	"sync"
)

// 当前桥接Map
var bridgeMap map[*TCPBridge]bool
var lock sync.Mutex

// Start 开始TCP桥接
func Start(isEncodeData bool, targetAddr string, npcTCP net.Conn) {
	bridge := &TCPBridge{
		isEncodeData: isEncodeData,
		NpsTCP:       npcTCP,
	}
	lock.Lock()
	bridgeMap[bridge] = true
	lock.Unlock()
	go bridge.start(targetAddr)
}

// 移除桥接通信
func removeBridge(bridge *TCPBridge) {
	lock.Lock()
	delete(bridgeMap, bridge)
	lock.Unlock()
}

// ShutdownAll 关闭所有链接
func ShutdownAll() {
	var closeList []*TCPBridge
	lock.Lock()
	for bridge := range bridgeMap {
		closeList = append(closeList, bridge)
	}
	bridgeMap = map[*TCPBridge]bool{}
	lock.Unlock()

	//关闭操作需要通知服务端，可能会造成阻塞，所以最好新开一个协程单独处理
	go func() {
		for _, bridge := range closeList {
			bridge.shutdown()
		}
	}()
}
