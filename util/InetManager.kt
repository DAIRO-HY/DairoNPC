package cn.dairo.clc.util

import java.net.InetAddress

object InetManager {

    /**
     * 创建InetAddress会阻塞线程,静态存储提升速度
     */
    private val mInetMap = HashMap<String, InetAddress>()

    /**
     * 获取InetAddress
     * @param name 域名或者IP
     */
    fun get(name: String): InetAddress {
        var inet = mInetMap[name]
        if (inet != null) {
            return inet
        }
        inet = InetAddress.getByName(name)
        mInetMap[name] = inet
        return inet
    }

    /**
     * 清空缓存
     */
    fun clear(){
        mInetMap.clear()
    }
}