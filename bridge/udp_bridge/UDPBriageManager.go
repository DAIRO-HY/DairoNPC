package udp_bridge

import (
	"net"
	"sync"
)

//桥接管理

/**
 * 当前桥接列表
 */
var mBridgeList = map[*UDPBridge]bool{}
var mBridgeListLock sync.Mutex

// 开始UDP桥接
func Start(isEncodeData bool, targetAddr string, clientSocket *net.UDPConn) {
	bridge := &UDPBridge{
		isEncodeData:  isEncodeData,
		targetAddr:    targetAddr,
		mClientSocket: clientSocket,
	}
	mBridgeListLock.Lock()
	mBridgeList[bridge] = true
	mBridgeListLock.Unlock()
	go bridge.start()
}

// 从桥接列表移除
func removeBridgeList(bridge *UDPBridge) {
	mBridgeListLock.Lock()
	delete(mBridgeList, bridge)
	mBridgeListLock.Unlock()
}

// 清空连接
func closeAll() {
	for bridge, _ := range mBridgeList { //@TODO: 这里调用close之后会不会执行removeBridgeList,待测试
		bridge.close()
	}
}

/**
 * 服务器向客户端同步当前处于激活状态的UDP连接端口
 */
//func syncServerActivePort(ports: String) = GlobalScope.launch(CLCDispatchers.IO) {
//    val serverActivePortList = HashSet(ports.split(",").map { it.toInt() })
//    var closeList: List<UDPBridge>? = null
//    this@UDPBriageManager.mBridgeListLock.withLock {
//
//        //筛选出需要关闭的连接池
//        closeList = if (ports == "0") {
//
//            //全部关闭
//            this@UDPBriageManager.mBridgeList.filter { true }
//        } else {
//
//            //筛选出非活性的连接
//            this@UDPBriageManager.mBridgeList.filter {
//                !serverActivePortList.contains(it.mClientSocket.localPort)
//            }
//        }
//    }
//    closeList!!.forEach {
//        it.close()
//    }
//}
