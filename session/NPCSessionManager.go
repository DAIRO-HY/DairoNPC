package session

import (
	"DairoNPC/bridge/tcp_bridge"
	"DairoNPC/bridge/udp_bridge"
	"DairoNPC/constant"
	"DairoNPC/pool/tcp_pool"
	"DairoNPC/pool/udp_pool"
	"fmt"
	"net"
	"sync"
	"time"
)

// 每隔一段时间检测心跳存活状态
const CHECK_HEART_TIME = constant.HEART_TIME * 3

// 最后一次收到服务器端心跳时间
var lastHeartTime int64

// 当前客户端会话
var npcSession *NPCSession

// 是否已经关闭
var isClose = false

// 是否正在运行中
var IsRuning = false

// 同步锁
var lock sync.Mutex

// 开启客户端
func Open() {
	if IsRuning { //如果正在运行中
		return
	}
	lock.Lock()
	isClose = false
	IsRuning = true
	fmt.Println("NPC服务开启成功")
	checkHeart()
	IsRuning = false
	lock.Unlock()
}

// 关闭客户端
// 公开给外部程序调用
func Close() {
	isClose = true
	if npcSession != nil {
		npcSession.shutdown()
	}
}

// 检测心跳
func checkHeart() {
	for {
		if isClose {
			if npcSession != nil {
				npcSession.shutdown()
			}
			return
		}
		if (time.Now().UnixMilli() - lastHeartTime) > CHECK_HEART_TIME { //长时间没有收到心跳，视为掉线

			//关闭上次会话
			if npcSession != nil {
				npcSession.shutdown()
			}

			//长时间没有收到心跳回复,重启客户端
			//println("-->${Date()}长时间没有收到心跳回复,重启客户端")
			//println("-->当前连接数:${TCPBriageManager.mBridgeList.count() + UDPBriageManager.mBridgeList.count()}/${(UDPPoolManager.mPoolList.count() + TCPPoolManager.mTcpPoolList.count())}")

			//关闭所有链接
			tcp_bridge.ShutdownAll()
			udp_bridge.ShutdownAll()

			tcp_pool.ShutdownAll()
			udp_pool.ShutdownAll()
			createConnection()
		}
		time.Sleep(CHECK_HEART_TIME * time.Millisecond)
	}
}

// 创建连接
func createConnection() {
	//fmt.Printf("createConnection:当前UDP连接池:%d UDP桥接数:%d \n", udp_pool.Count(), udp_bridge.Count())

	// 与服务端建立连接
	tcp, err := net.Dial("tcp", constant.Host+":"+constant.TcpPort)
	if err != nil {
		fmt.Println("-->与主机连接失败:" + constant.Host + ":" + constant.TcpPort)
		return
	}
	npcSession = &NPCSession{
		npsTCP: tcp,
	}
	go npcSession.start()
}
