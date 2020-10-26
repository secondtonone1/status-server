::编译linux版本
cd src
set GOOS=linux
set GOARCH=amd64
set GOHOSTOS=linux
go.exe build  -o ../bin/client main/client.go
go.exe build  -o ../bin/status-server main/main.go
cd ../
::go modules proxy 设置 https://goproxy.io