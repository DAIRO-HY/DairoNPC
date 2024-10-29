package cn.dairo.npc.activity

import android.app.Activity
import android.os.Bundle
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import android.widget.Toast
import cn.dairo.npc.R
import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlin.concurrent.thread

class MainActivity : Activity() {

    /**
     * 服务器
     */
    private lateinit var etHost: EditText

    /**
     * 客户端KEY
     */
    private lateinit var etKey: EditText

    /**
     * TCP端口
     */
    private lateinit var etTCP: EditText

    /**
     * UDP端口
     */
    private lateinit var etUDP: EditText

    /**
     * 运行状态
     */
    private lateinit var tvState: TextView

    /**
     * 轮询显示运行状态是否在运行中
     */
    private var isShowStateRuning = false

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.layout_main)
        initView()
    }

    private fun initView() {
        this.etHost = findViewById(R.id.etHost)
        this.etKey = findViewById(R.id.etKey)
        this.etTCP = findViewById(R.id.etTCP)
        this.etUDP = findViewById(R.id.etUDP)
        this.tvState = findViewById(R.id.tvState)
        findViewById<Button>(R.id.btnConnect).setOnClickListener {
            this.onConnectClick()
        }
        findViewById<Button>(R.id.btnClose).setOnClickListener {//关闭按钮点击事件
            DairoNPC.DairoNPC.close()
        }

        //测试用
//        this.etHost.setText("172.16.3.91")
//        this.etKey.setText("njeHds*fs4tfsd")
    }

    /**
     * 连接按钮点击事件
     */
    private fun onConnectClick() {
        val host = etHost.text.toString()
        if (host.isEmpty()) {
            Toast.makeText(this, "服务器必填", Toast.LENGTH_SHORT).show()
            return
        }
        val key = etKey.text.toString()
        if (key.isEmpty()) {
            Toast.makeText(this, "客户端Key必填", Toast.LENGTH_SHORT).show()
            return
        }

        var tcp = etTCP.text.toString()
        if (tcp.isEmpty()) {
            tcp = "1781"
        }

        var udp = etUDP.text.toString()
        if (udp.isEmpty()) {
            udp = "1782"
        }
        val connectStr = "-h:$host\n-k:$key\n-t:$tcp\n-u:$udp"
        thread {
            showState()
            DairoNPC.DairoNPC.open(connectStr)
        }
    }

    /**
     * 显示运行状态
     */
    @OptIn(DelicateCoroutinesApi::class)
    private fun showState(){
        if(this.isShowStateRuning){
            return
        }
        this.isShowStateRuning = true
        GlobalScope.launch {
            while (true){
                if(DairoNPC.DairoNPC.isRuning()){
                    runOnUiThread{
                        this@MainActivity.tvState.text = "正在运行"
                    }
                }else{
                    this@MainActivity.tvState.text = "已停止"
                }
                delay(1000)
            }
        }
    }
}