package constant

// 客户端版本号
const VERSION = "1.1.1"

// 心跳间隔时间(毫秒)
const HEART_TIME = 3000

// 关闭UDP连接池标记
const UDP_POOL_CLOSE_FLAG = ":[NPS-POOL-CLOSE]"

// 关闭UDP桥接标记
const UDP_BRIDIGE_CLOSE_FLAG = ":[NPS-BRIDGE-CLOSE]"

// 服务器
var Host string

// 服务器端TCP端口
var TcpPort string

// 服务器端UDP端口
var UdpPort string

// 认证秘钥
var Key string

// 客户端id,该值有服务器端返回
var ClientId int
