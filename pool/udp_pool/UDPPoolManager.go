package udp_pool

import (
	"DairoNPC/constant"
	"fmt"
	"net"
	"sync"
)

// Socket连接池管理

/**
 * 当前激活UDP连接池
 */
var mPoolList = make(map[*UDPPool]bool)
var mPoolListLock sync.Mutex

func Count() int {
	return len(mPoolList)
}

/**
 * 创建UDP连接池
 */
func Create(count int) {

	// 创建一个 UDP 地址
	serverAddr, err := net.ResolveUDPAddr("udp", constant.Host+":"+constant.UdpPort)
	if err != nil {
		fmt.Println("创建udp连接失败", err)
		return
	}

	mPoolListLock.Lock()
	for i := 0; i < count; i++ {

		// 创建一个 UDP 连接
		udp, err := net.DialUDP("udp", nil, serverAddr)
		if err != nil {
			fmt.Println("创建udp连接失败", err)
			return
		}
		pool := &UDPPool{
			NpsUDP: udp,
		}
		mPoolList[pool] = true
		go pool.start()
	}
	mPoolListLock.Unlock()
}

// 连接池列表移除
func removePoolList(pool *UDPPool) {
	mPoolListLock.Lock()
	delete(mPoolList, pool)
	mPoolListLock.Unlock()
}

// 关闭连接池所有链接
func ShutdownAll() {
	mPoolListLock.Lock()
	for pool, _ := range mPoolList {

		//关闭后,UDPPool内部会自动调用removePoolList,这里无需手动调用
		pool.close()
	}
	mPoolListLock.Unlock()
}

/**
 * 服务器向客户端同步当前处于激活状态的UDP连接池端口
 */
//func syncServerActivePort(ports: String) = GlobalScope.launch(CLCDispatchers.IO) {
//    val serverActivePortList = HashSet(ports.split(",").map { it.toInt() })
//    var closeList: List<UDPPool>? = null
//    this@UDPPoolManager.mPoolListLock.withLock {
//
//        //筛选出需要关闭的连接池
//        closeList = if (ports == "0") {
//
//            //全部关闭
//            this@UDPPoolManager.mPoolList.filter { true }
//        } else {
//            this@UDPPoolManager.mPoolList.filter {
//                !serverActivePortList.contains(it.socket.localPort)
//            }
//        }
//    }
//    closeList!!.forEach {
//        it.close()
//    }
//}
