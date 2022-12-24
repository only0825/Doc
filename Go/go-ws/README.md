# go-ws


Mac编译 Mac执行：
go build -o go-data go-data

Mac编译 linux下去执行
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build -o go-data-linux go-data