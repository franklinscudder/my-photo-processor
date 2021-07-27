@echo off
setlocal 
    set GOOS=linux
    set GOARCH=amd64
    echo Building for linux-amd64...
    go build -o bin/photo-processor-linux-amd64.bin ./photo-processor.go
    set GOARCH=arm
    echo Building for linux-arm...
    go build -o bin/photo-processor-linux-arm.bin ./photo-processor.go
endlocal
echo Building for win64...
go build -o bin/photo-processor-win64.exe ./photo-processor.go


