@echo off
SET GOOS=windows
SET GOARCH=amd64
@go build -ldflags "-s -w" -o webserver.exe "%CD%\cmd\main.go"
