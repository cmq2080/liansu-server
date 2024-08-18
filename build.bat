@echo off
del .\\deploy\\liansu-server.exe
%GOPATH%\\bin\\rsrc.exe -manifest main.manifest -o rsrc.syso
rem go build -o liansu-server.exe -ldflags="-H windowsgui"
%GOROOT%\\bin\\go.exe build -o .\\deploy\\liansu-server.exe
@echo on