#!/bin/bash

#github登录票据
github_token=$GITHUB_TOKEN

#docker用户名
docker_user=$DOCKER_USER
docker_pwd=$DOCKER_PWD

#项目名
projectName="DairoNPC"

repo="DAIRO-HY/$projectName"
branch="release"

#--------------------------------------获取代码-----------------------------------------
if [ -d $projectName ]; then
    cd $projectName

    #删除所有新添加的文件
    git clean -f

    #取消所有更改
    git reset --hard
    git pull
else
    CLONE_URL="https://${github_token}@github.com/${repo}.git"
    git clone --branch $branch $CLONE_URL
    cd $projectName
fi

# 获取版本号
version=$(grep -oP '(?<=VERSION = ")[^"]+' ./constant/NPCConstant.go)
#------------------------------------------编译二进制文件-----------------------------------------
echo "正在编译二进制文件..."
if [ -d "./build" ]; then
    rm -rf "./build"
fi
mkdir "./build"

#编译linux-amd64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ./build/dairo-npc-linux-amd64

#编译linux-arm64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ./build/dairo-npc-linux-arm64

#编译linux-arm
GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -ldflags="-s -w" -o ./build/dairo-npc-linux-arm

#编译win-amd64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ./build/dairo-npc-win-amd64.exe

#编译mac-amd64
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ./build/dairo-npc-mac-amd64

#编译mac-arm64
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ./build/dairo-npc-mac-arm64

if [ ! -e "./build/dairo-npc-mac-arm64" ]; then
    echo "编译失败"
    exit 1
fi

#------------------------------------------编译APK-----------------------------------------
echo "正在编译APK..."
mv main.go main.go.bk
mv DairoNPC.go.bk DairoNPC.go

#编译aar
gomobile bind "-target=android/amd64,android/arm64,android/arm" -androidapi 21 -o ./android/app/libs/dairo-npc.aar DairoNPC
cd ./android/app
old_version_string="versionName = \"1.0\""
new_version_string="versionName = \"$version\""
sed -i "s/${old_version_string}/${new_version_string}/g" build.gradle.kts

#打包apk
gradle clean
gradle assembleRelease

if [ ! -e "./build/outputs/apk/release/app-release.apk" ]; then
    echo "编译失败"
    exit 1
fi

#回到项目根目录
cd ../../


#---------------------------------------创建标签----------------------------------------
echo "正在创建标签..."

#删除本地已经存在的标签
git tag -d $version

#删除远程标签
git push origin --delete tag $version

git tag $version
git push origin $version

release_message="本次发布版本:$version"
create_release_api_response=$(curl -L -X POST "https://api.github.com/repos/$repo/releases" \
                        -H "Accept: application/vnd.github.v3+json" \
                        -H "Authorization: Bearer $github_token" \
                        -H "X-GitHub-Api-Version: 2022-11-28" \
                        -d "{\"tag_name\":\"$version\",\"name\":\"$version\",\"body\":\"$release_message\"}")
echo "创建标签结果:${create_release_api_response}"

#通过正则匹配ReleaseId, [head -n 1]功能是从匹配到的多个字符串中去第一个字符串
release_id=$(echo "$create_release_api_response" | grep -oP '(?<="id": )[^,]+' | head -n 1)
echo "本地发布ID:${release_id}"


#---------------------------------------上传编译好的二进制文件----------------------------------
echo "正在上传编译好的二进制文件..."
upload_file_api_response=$(curl -s -X POST \
                            -H "Accept: application/vnd.github+json" \
                            -H "Authorization: Bearer ${github_token}" \
                            -H "X-GitHub-Api-Version: 2022-11-28" \
                            -H "Content-Type: application/octet-stream" \
                            --data-binary "@./build/dairo-npc-linux-amd64" \
                            "https://uploads.github.com/repos/${repo}/releases/${release_id}/assets?name=dairo-npc-linux-amd64")
echo "上传文件dairo-npc-linux-amd64结果:${upload_file_api_response}"

upload_file_api_response=$(curl -s -X POST \
                            -H "Accept: application/vnd.github+json" \
                            -H "Authorization: Bearer ${github_token}" \
                            -H "X-GitHub-Api-Version: 2022-11-28" \
                            -H "Content-Type: application/octet-stream" \
                            --data-binary "@./build/dairo-npc-linux-arm64" \
                            "https://uploads.github.com/repos/${repo}/releases/${release_id}/assets?name=dairo-npc-linux-arm64")
echo "上传文件dairo-npc-linux-arm64结果:${upload_file_api_response}"

upload_file_api_response=$(curl -s -X POST \
                            -H "Accept: application/vnd.github+json" \
                            -H "Authorization: Bearer ${github_token}" \
                            -H "X-GitHub-Api-Version: 2022-11-28" \
                            -H "Content-Type: application/octet-stream" \
                            --data-binary "@./build/dairo-npc-linux-arm" \
                            "https://uploads.github.com/repos/${repo}/releases/${release_id}/assets?name=dairo-npc-linux-arm")
echo "上传文件dairo-npc-linux-arm结果:${upload_file_api_response}"

upload_file_api_response=$(curl -s -X POST \
                            -H "Accept: application/vnd.github+json" \
                            -H "Authorization: Bearer ${github_token}" \
                            -H "X-GitHub-Api-Version: 2022-11-28" \
                            -H "Content-Type: application/octet-stream" \
                            --data-binary "@./build/dairo-npc-win-amd64.exe" \
                            "https://uploads.github.com/repos/${repo}/releases/${release_id}/assets?name=dairo-npc-win-amd64.exe")
echo "上传文件dairo-npc-win-amd64.exe结果:${upload_file_api_response}"

upload_file_api_response=$(curl -s -X POST \
                            -H "Accept: application/vnd.github+json" \
                            -H "Authorization: Bearer ${github_token}" \
                            -H "X-GitHub-Api-Version: 2022-11-28" \
                            -H "Content-Type: application/octet-stream" \
                            --data-binary "@./build/dairo-npc-mac-amd64" \
                            "https://uploads.github.com/repos/${repo}/releases/${release_id}/assets?name=dairo-npc-mac-amd64")
echo "上传文件dairo-npc-mac-amd64结果:${upload_file_api_response}"

upload_file_api_response=$(curl -s -X POST \
                            -H "Accept: application/vnd.github+json" \
                            -H "Authorization: Bearer ${github_token}" \
                            -H "X-GitHub-Api-Version: 2022-11-28" \
                            -H "Content-Type: application/octet-stream" \
                            --data-binary "@./build/dairo-npc-mac-arm64" \
                            "https://uploads.github.com/repos/${repo}/releases/${release_id}/assets?name=dairo-npc-mac-arm64")
echo "上传文件dairo-npc-mac-arm64结果:${upload_file_api_response}"

upload_file_api_response=$(curl -s -X POST \
                            -H "Accept: application/vnd.github+json" \
                            -H "Authorization: Bearer ${github_token}" \
                            -H "X-GitHub-Api-Version: 2022-11-28" \
                            -H "Content-Type: application/octet-stream" \
                            --data-binary "@./android/app/build/outputs/apk/release/app-release.apk" \
                            "https://uploads.github.com/repos/${repo}/releases/${release_id}/assets?name=dairo-npc.apk")
echo "上传文件dairo-npc.apk结果:${upload_file_api_response}"


#---------------------------------------上传Docker镜像-----------------------------------------
echo "正在打包Docker镜像..."
mv ./build/dairo-npc-linux-amd64 ./document/docker/
cd ./document/docker/
docker build -t $docker_user/dairo-npc:$version .
docker login -u $docker_user --password $docker_pwd
docker push $docker_user/dairo-npc:$version
docker logout

echo "---------------------------------------docker镜像推送完成--------------------------------------"
