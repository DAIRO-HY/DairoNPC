package cn.dairo.clc.extension

import java.text.SimpleDateFormat
import java.util.*

/**
 * 日期扩展
 * 默认格式yyyy-MM-dd
 * @return 格式化之后的日期
 */
fun Date?.format(pattern: String = "yyyy-MM-dd"): String? {
    if (this == null) {
        return null
    }
    return SimpleDateFormat(pattern).format(this)
}