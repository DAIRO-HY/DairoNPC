package constant

/**
 * 客户端版本号
 */
//val version by lazy {
//    String(this.javaClass.classLoader.getResourceAsStream("version.txt").readAllBytes())
//}

/**
 * 心跳间隔时间
 */
const Version = "1.0"

/**
 * 心跳间隔时间
 */
const HEART_TIME = 3000

/**
 * 关闭UDP连接池标记
 */
const CLOSE_UDP_POOL_FLAG = "CLOSE"

/**
 * 服务主机
 */
var Host string

/**
 * 服务器端TCP端口
 */
var TcpPort string

/**
 * 服务器端UDP端口
 */
var UdpPort string

/**
 * 认证秘钥
 */
var Key string

/**
 * 客户端id,该值有服务器端返回
 */
var ClientId int

///**
// * 获取TCP连接
// */
//val tcp: Socket?
//    get() = try {// 与服务端建立连接
//        val clientSocket = Socket(this.host, this.tcpPort)
//        clientSocket.keepAlive = true
//        clientSocket.soTimeout = 0
//        clientSocket
//    } catch (e: Exception) {//创建socket失败
//        null
//    }
//
///**
// * 获取UDP连接
// */
//val udp: DatagramSocket?
//    get() {
//        repeat(100) {
//            try {
//
//                //因为端口是随机占用,所以这里有端口被占用的可能
//                return DatagramSocket()
//            } catch (e: BindException) {
//            }
//        }
//        return null
//    }
