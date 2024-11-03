package session

import (
	"DairoNPC/HeaderUtil"
	"DairoNPC/constant"
	"DairoNPC/extension"
	"DairoNPC/pool/tcp_pool"
	"DairoNPC/pool/udp_pool"
	"DairoNPC/util/SecurityUtil"
	"DairoNPC/util/TcpUtil"
	"fmt"
	"net"
	"strconv"
	"time"
)

// 与服务端通信连接
type NPCSession struct {
	npsTCP net.Conn
}

/**
 * 开始
 */
func (mine *NPCSession) start() {
	mine.readServerInfoAndReceive()

	//关闭会话
	mine.npsTCP.Close()
}

// 从服务端读取基本信息
func (mine *NPCSession) readServerInfoAndReceive() {
	header := constant.Key + "|" + constant.VERSION

	//发送标记信息
	if HeaderUtil.SendFlag(mine.npsTCP, HeaderUtil.CLIENT_TO_SERVER_MAIN_CONNECTION, header) != nil {
		return
	}

	//获取客户端ID
	if mine.readClientId() != nil {
		return
	}

	//客户端加解密秘钥
	if mine.readClientSecurityKey() != nil {
		return
	}

	//发送心跳数据
	go mine.heart()

	//接收数据
	mine.receive()
}

/**
 * 获取客户端ID
 */
func (mine *NPCSession) readClientId() error {

	//第一个字节为标记
	flagData, err := TcpUtil.ReadNByte(mine.npsTCP, 1)
	if err != nil {
		return err
	}
	if flagData[0] != HeaderUtil.SERVER_TO_CLIENT_ID {
		mine.npsTCP.Close()
		return &extension.BusinessException{
			Message: "非法标记:$flag",
		}
	}

	//得到头部数据
	header, err := HeaderUtil.GetHeader(mine.npsTCP)
	if err != nil {
		return err
	}

	//得到客户端ID
	clientId, _ := strconv.ParseInt(header, 10, 64)

	//得到客户端ID
	constant.ClientId = int(clientId)
	return nil
}

// 客户端加解密秘钥
func (mine *NPCSession) readClientSecurityKey() error {
	clientSecurityKey, err := TcpUtil.ReadNByte(mine.npsTCP, 256)
	if err != nil {
		return err
	}

	//将数据复制到数组中
	copy(SecurityUtil.ClientSecurityKey[:], clientSecurityKey)
	return nil
}

/**
 * 从服务端收到数据
 */
func (mine *NPCSession) receive() {
	for {
		flagData, err := TcpUtil.ReadNByte(mine.npsTCP, 1)
		if err != nil {
			return
		}
		flag := flagData[0]
		//fmt.Printf("-->收到标记：%d : %c\n", flag, rune(flag))
		switch flag {

		//服务器向客户端申请TCP连接池请求
		case HeaderUtil.REQUEST_TCP_POOL:
			header, err := HeaderUtil.GetHeader(mine.npsTCP)
			if err != nil {
				return
			}

			//创建数量
			count, _ := strconv.ParseInt(header, 10, 64)

			//创建连接池
			tcp_pool.Create(int(count))

		//服务器向客户端申请UDP连接池请求
		case HeaderUtil.REQUEST_UDP_POOL:
			header, err := HeaderUtil.GetHeader(mine.npsTCP)
			if err != nil {
				return
			}

			//创建数量
			count, _ := strconv.ParseInt(header, 10, 64)

			//创建连接池
			udp_pool.Create(int(count))

		//服务器端回复了心跳
		case HeaderUtil.MAIN_HEART_BEAT:
			//println("-->收到服务器心跳数据:${System.currentTimeMillis()}")
			//fmt.Printf("当前UDP连接池:%d UDP桥接数:%d \n", udp_pool.Count(), udp_bridge.Count())
			lastHeartTime = time.Now().UnixMilli()

		////服务器向客户端同步当前处于激活状态的UDP连接池端口
		//case HeaderUtil.SYNC_ACTIVE_POOL_UDP_PORT : {
		//    val ports = HeaderUtil.getHeader(this.npcTCP) ?: continue
		//    UDPPoolManager.syncServerActivePort(ports)
		//}
		//
		//向客户端同步当前保留的UDP连接端口
		case HeaderUtil.SYNC_ACTIVE_BRIDGE_UDP_PORT:
			ports, err := HeaderUtil.GetHeader(mine.npsTCP)
			if err != nil {
				return
			}
			fmt.Println(ports)
			//UDPBriageManager.syncServerActivePort(ports)

		}
	}
}

// 定期心跳
func (mine *NPCSession) heart() {
	for { //每个一段时间发送一次心跳包
		time.Sleep(constant.HEART_TIME * time.Millisecond)
		err := HeaderUtil.SendFlag(mine.npsTCP, HeaderUtil.MAIN_HEART_BEAT, "")
		if err != nil {
			break
		}
	}
}

/**
 * 发送标记数据
 */
//func (mine *NPCSession) send(flag uint8, message string) error {
//	data := []uint8{flag}
//	if len(message) > 0 {
//		data = append(data, uint8(len(message)))
//
//		//向服务器端发送客户端ID信息
//		data = append(data, []uint8(message)...)
//	}
//	err := TcpUtil.WriteAll(mine.npcTCP, data)
//	if err != nil {
//		mine.npcTCP.Close()
//		return err
//	}
//	return nil
//}

/**
 * 关闭服务
 */
func (mine *NPCSession) shutdown() {
	mine.npsTCP.Close()
}
