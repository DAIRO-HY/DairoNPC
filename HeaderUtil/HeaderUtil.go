package HeaderUtil

import (
	"DairoNPC/util/TcpUtil"
	"net"
)

// 客户端与服务器端通信连接标记
const CLIENT_TO_SERVER_MAIN_CONNECTION = 0

// 与客户端通信心跳标记
const MAIN_HEART_BEAT = 1

// 向客户端发送clientId
const SERVER_TO_CLIENT_ID = 2

// 向客户端申请TCP连接池请求
const REQUEST_TCP_POOL = 3

/**
 * 向客户端申请UDP连接池请求
 */
const REQUEST_UDP_POOL = 4

/**
 * 服务器向客户端同步当前处于激活状态的UDP连接池端口
 */
const SYNC_ACTIVE_POOL_UDP_PORT = 5

/**
 * 向客户端同步当前处于激活状态的UDP连接端口
 */
const SYNC_ACTIVE_BRIDGE_UDP_PORT = 6

/**
 * 向客户端发送clientId
 */
const SECURITY_CLIENT_KEY = 7

/**
 * 获取客户端Socket头部信息
 */
func GetHeader(clientSocket net.Conn) (string, error) {

	//读取一个字节,该字节代表key长度
	lenData, err := TcpUtil.ReadNByte(clientSocket, 1)
	if err != nil {
		return "", err
	}

	//得到头部部分数据长度
	headerLen := lenData[0]
	headerData, err := TcpUtil.ReadNByte(clientSocket, int(headerLen))
	if err != nil {
		return "", err
	}
	return string(headerData), nil
}

// SendFlag 发送一个标记
// - flag 标记
// - header 头部信息
func SendFlag(tcp net.Conn, flag uint8, header string) error {
	data := []uint8{flag}
	if len(header) > 0 {
		data = append(data, uint8(len(header)))

		//向服务器端发送客户端ID信息
		data = append(data, []uint8(header)...)
	}
	return TcpUtil.WriteAll(tcp, data)
}
