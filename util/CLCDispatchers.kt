package cn.dairo.clc.util

import kotlinx.coroutines.Dispatchers

object CLCDispatchers {

    /**
     * IO调度器,可动态调整分配最大线程数
     */
    val IO = Dispatchers.IO.limitedParallelism(1024)
}