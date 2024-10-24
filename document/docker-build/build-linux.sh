#!/bin/bash

#删除上次编译文件
rm DairoNPC.zip
rm -rf DairoNPC-main
rm /app/DairoNPC/dairo-npc-linux-amd64

curl -L -o DairoNPC.zip https://github.com/DAIRO-HY/DairoNPC/archive/refs/heads/main.zip?561
unzip DairoNPC.zip
cd DairoNPC-main

#开始编译
go build -o /app/DairoNPC/dairo-npc-linux-amd64

/app/DairoNPC/dairo-npc-linux-amd64