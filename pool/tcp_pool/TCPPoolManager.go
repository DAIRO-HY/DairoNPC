package tcp_pool

import (
	"DairoNPC/constant"
	"net"
	"sync"
)

// 当前激活的TCP连接池
var mTcpPoolList map[*TCPPool]bool
var lock sync.Mutex

// Create 创建TCP连接池
func Create(count int) {
	for i := 0; i < count; i++ {

		//与目标端口建立连接
		tcp, err := net.Dial("tcp", constant.Host+":"+constant.TcpPort)
		if err != nil {
			return
		}
		pool := &TCPPool{
			npcTCP: tcp,
		}
		lock.Lock()
		mTcpPoolList[pool] = true
		lock.Unlock()
		go pool.start()
	}
}

// 连接池列表移除
func removePool(pool *TCPPool) {
	lock.Lock()
	delete(mTcpPoolList, pool)
	lock.Unlock()
}

// ShutdownAll 关闭连接池所有链接
func ShutdownAll() {
	var closeList []*TCPPool
	lock.Lock()
	for pool := range mTcpPoolList {
		closeList = append(closeList, pool)
	}
	mTcpPoolList = map[*TCPPool]bool{}
	lock.Unlock()

	//关闭操作需要通知服务端，可能会造成阻塞，所以最好新开一个协程单独处理
	go func() {
		for _, pool := range closeList {
			pool.npcTCP.Close()
		}
	}()
}
