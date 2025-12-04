@ECHO off

SET BIN_PATH="./bin"
IF NOT EXIST %BIN_PATH% (
    mkdir %BIN_PATH%
)

SET CGO_ENABLED=0

SET GOOS=windows
ECHO compiling to windows...

SET GOARCH=amd64
go build -o %BIN_PATH%/srunClient_windows_x64.exe .
ECHO [+] srunClient_windows_x64.exe

SET GOARCH=386
go build -o %BIN_PATH%/srunClient_windows_x32.exe .
ECHO [+] srunClient_windows_x32.exe

SET GOOS=linux
ECHO compiling to linux...

SET GOARCH=amd64
go build -o %BIN_PATH%/srunClient_linux_x64 .
ECHO [+] srunClient_linux_x64

SET GOARCH=386
go build -o %BIN_PATH%/srunClient_linux_x32 .
ECHO [+] srunClient_linux_x32

SET GOARCH=arm64
go build -o %BIN_PATH%/srunClient_linux_arm64 .
ECHO [+] srunClient_linux_arm64

SET GOARCH=arm
go build -o %BIN_PATH%/srunClient_linux_arm32 .
ECHO [+] srunClient_linux_arm32

SET GOOS=darwin
ECHO compiling to darwin...

SET GOARCH=amd64
go build -o %BIN_PATH%/srunClient_darwin_x64 .
ECHO [+] srunClient_darwin_x64

SET GOARCH=arm64
go build -o %BIN_PATH%/srunClient_darwin_arm64 .
ECHO [+] srunClient_darwin_arm64

ECHO done