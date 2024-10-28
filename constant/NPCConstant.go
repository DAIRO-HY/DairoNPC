package constant

// 客户端版本号
const VERSION = "1.0.2"

// 心跳间隔时间(毫秒)
const HEART_TIME = 3000

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
