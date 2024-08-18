rm -rf deploy/liansu-server
$GOPATH/bin/rsrc.exe -manifest main.manifest -o rsrc.syso
# go build -o liansu-server.exe -ldflags="-H windowsgui"
GOARCH=amd64 GOOS=linux go build -o deploy/liansu-server