package session

import (
	"DairoNPC/bridge"
	"DairoNPC/constant"
	"DairoNPC/pool"
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
	checkHeart()
	IsRuning = false
	lock.Unlock()
}

// 关闭客户端
func Close() {
	isClose = true
}

// 检测心跳
func checkHeart() {
	for {
		if isClose {
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
			bridge.ShutdownAll()
			pool.ShutdownAll()
			createConnection()
		}
		time.Sleep(CHECK_HEART_TIME * time.Millisecond)
	}
}

// 创建连接
func createConnection() {

	// 与服务端建立连接
	tcp, err := net.Dial("tcp", constant.Host+":"+constant.TcpPort)
	if err != nil {
		fmt.Println("-->与主机连接失败:" + constant.Host + ":" + constant.TcpPort)
		return
	}
	npcSession = &NPCSession{
		npcTCP: tcp,
	}
	go npcSession.start()
	fmt.Println("-->服务开启成功")
}
