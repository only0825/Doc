#! /bin/bash

git checkout main
git pull origin main
git submodule init
git submodule update

DIR=$(pwd)
PROJECT="merchant"
GitReversion=`git rev-parse HEAD`
BuildTime=`date +'%Y.%m.%d.%H%M%S'`
BuildGoVersion=`go version`

go build -ldflags "-X main.gitReversion=${GitReversion}  -X 'main.buildTime=${BuildTime}' -X 'main.buildGoVersion=${BuildGoVersion}'" -o $PROJECT
# shellcheck disable=SC2164
cd /opt/deploy/cg/$PROJECT
git checkout uat
git pull origin uat
mv $DIR/$PROJECT /opt/deploy/cg/$PROJECT
git commit -am "${GitReversion}"
git push

