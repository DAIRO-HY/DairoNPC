FROM ubuntu:24.10

RUN apt update
USER root

#安装xz解压工具，否则可能导致解压xz失败
RUN yes|apt-get install xz-utils
RUN yes|apt-get install unzip

#安装curl
RUN yes|apt install curl

#安装git
RUN yes|apt install git

#安装sqlite编译环境
RUN yes|apt-get install libsqlite3-dev

#安装gcc,编译sqlite需要
RUN yes|apt-get install gcc


#####################################################################安装docker开始#######################################
RUN mkdir /etc/docker
#COPY daemon.json /etc/docker
RUN curl -fsSL get.docker.com -o get-docker.sh

##若执行报错,请加参数 --mirror Aliyun
#RUN sh get-docker.sh --mirror Aliyun
RUN sh get-docker.sh

#安装docker-compose
RUN curl -SL https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
RUN chmod +x /usr/local/bin/docker-compose
####################################################################安装docker结束#######################################

WORKDIR /usr/local
#####################################################################安装JDK-START#######################################

## 将jdk压缩包添加到容器的 /root 目录,ADD指令会自动解压,解压之后的文件名:jdk1.8.0_241
## openjdk官方下载地址https://jdk.java.net/archive/
ADD https://download.java.net/java/GA/jdk17.0.1/2a2082e5a09d4267845be086888add4f/12/GPL/openjdk-17.0.1_linux-x64_bin.tar.gz /usr/local/openjdk.tar.gz

RUN tar -zvxf openjdk.tar.gz
RUN rm -rf openjdk.tar.gz


# 配置JAVA_HOME环境变量
ENV JAVA_HOME /usr/local/jdk-17.0.1/

# 将JAVA_HOME/bin 添加至PATH环境变量
ENV PATH $JAVA_HOME/bin:$PATH
####################################################################安装JDK-END##########################################





#######################################################################安装gradle开始#####################################
ENV GRADLE_VERSION 8.5
ADD https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-all.zip /usr/local/gradle-all.zip

RUN unzip gradle-all.zip

#重命名文件夹
RUN mv gradle-${GRADLE_VERSION} gradle

#删除安装包
RUN rm -rf /usr/local/gradle-all.zip

# 配置GRADLE_HOME环境变量
ENV GRADLE_HOME /usr/local/gradle

# 将GRADLE_HOME/bin 添加至PATH环境变量
ENV PATH $PATH:$GRADLE_HOME/bin

# 配置GRADLE_USER_HOME环境变量,默认jar包下载位置,这个目录最好是做数据卷映射，否则每次进入容器都会重新下载依赖
ENV GRADLE_USER_HOME /usr/.gradle
#####################################################################安装gradle结束#######################################



##########################################################安装Andrid SDK相关#######################################START
WORKDIR /home/SDK/cmdline-tools
RUN curl -o commandlinetools.zip https://dl.google.com/android/repository/commandlinetools-linux-11076708_latest.zip
RUN unzip commandlinetools.zip
RUN rm commandlinetools.zip
RUN mv cmdline-tools latest
ENV ANDROID_HOME="/home/SDK"
ENV PATH="$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools:$ANDROID_HOME/tools/bin:$ANDROID_HOME/cmdline-tools/latest/bin:${PATH}"

#安装NDK
RUN yes|sdkmanager "ndk;28.0.12433566"

#安装build-tools
RUN yes|sdkmanager "build-tools;34.0.0"

#安装platform-tools
RUN yes|sdkmanager "platform-tools"

#安装platforms
RUN yes|sdkmanager "platforms;android-34"

#添加授权
RUN yes|sdkmanager --licenses
##########################################################安装Andrid SDK相关#######################################END




##########################################################安装GO相关#######################################START
WORKDIR /opt

#下载go
ADD https://go.dev/dl/go1.23.1.linux-amd64.tar.gz /opt

#解压
RUN tar -xzvf go1.23.1.linux-amd64.tar.gz

#删除安装包
RUN rm go1.23.1.linux-amd64.tar.gz

# 追加到PATH环境变量
ENV PATH="/opt/go/bin:${PATH}"

#安装gomobile
RUN go install golang.org/x/mobile/cmd/gomobile@latest

# 追加到PATH环境变量
ENV PATH="/root/go/bin:${PATH}"

RUN gomobile init
##########################################################安装GO相关#######################################END


WORKDIR /home/nps-build
COPY build-linux.sh build-linux.sh
RUN chmod +x build-linux.sh

#运行程序
ENTRYPOINT "./build-linux.sh"
