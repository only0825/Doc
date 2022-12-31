#! /bin/bash

PROJECT="merchant"
GitReversion=`git rev-parse HEAD`
BuildTime=`date +'%Y.%m.%d.%H%M%S'`
BuildGoVersion=`go version`

go build -ldflags "-X main.gitReversion=${GitReversion}  -X 'main.buildTime=${BuildTime}' -X 'main.buildGoVersion=${BuildGoVersion}'" -o $PROJECT

scp -i /home/gocloud-yiy-rich -P 10087 $PROJECT p3test@34.92.240.177:/home/centos/workspace/cg/merchant/merchant_cg
ssh -i /home/gocloud-yiy-rich -p 10087 p3test@34.92.240.177 "sh /home/centos/workspace/cg/merchant/p3.sh"