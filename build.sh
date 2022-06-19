rm -rf liansu-server.exe
$GOPATH/bin/rsrc.exe -manifest main.manifest -o rsrc.syso
# go build -o liansu-server.exe -ldflags="-H windowsgui"
go build -o liansu-server.exe