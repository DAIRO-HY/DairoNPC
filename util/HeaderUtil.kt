package cn.dairo.clc.util

import java.net.Socket
import java.net.SocketTimeoutException

/**
 * 该类直接从服务端复制
 */
object HeaderUtil {

    /**
     * 客户端与服务器端通信连接标记
     */
    const val CLIENT_TO_SERVER_MAIN_CONNECTION = 0

    /**
     * 与客户端通信心跳标记
     */
    const val MAIN_HEART_BEAT = 1

    /**
     * 向客户端发送clientId
     */
    const val SERVER_TO_CLIENT_ID = 2

    /**
     * 向客户端申请TCP连接池请求
     */
    const val SERVER_TCP_POOL_REQUEST = 3

    /**
     * 向客户端申请UDP连接池请求
     */
    const val SERVER_UDP_POOL_REQUEST = 4

    /**
     * 服务器向客户端同步当前处于激活状态的UDP连接池端口
     */
    const val SYNC_ACTIVE_POOL_UDP_PORT = 5

    /**
     * 向客户端同步当前处于激活状态的UDP连接端口
     */
    const val SYNC_ACTIVE_BRIDGE_UDP_PORT = 6


    /**
     * 获取客户端Socket头部信息
     */
    fun getHeader(clientSocket: Socket): String? {

        //得到输入流
        val iStream = clientSocket.inputStream

        //设置读取数据超时
        clientSocket.soTimeout = 3000

        try {//读取一个字节,该字节代表key长度
            val headerLen = iStream.read()
            val headerData = ByteArray(headerLen)

            //一直读到指定长度的数据
            val realLen = iStream.readNBytes(headerData, 0, headerData.size)
            if (realLen != headerLen) {

                //如果读取到的长度与实际key长度不匹配,说明数据不完成,则该连接作废
                clientSocket.close()
            }

            //设置读取数据超时
            clientSocket.soTimeout = 0

            //得到头部信息
            val header = String(headerData)
            return header
        } catch (e: SocketTimeoutException) {
            println("-->服务器端未在指定时间内收到头部响应数据")
            e.printStackTrace()
        } catch (e: Exception) {
            e.printStackTrace()
        }finally {

            //设置永不超时
            clientSocket.soTimeout = 0
        }
        return null
    }
}