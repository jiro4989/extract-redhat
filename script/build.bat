@echo off

set binname=%1

set srcfile="src\%binname%.go"
set bindir="bin"
mkdir %bindir%

set GOOS=windows
set GOARCH=amd64
go build -o %bindir%\%binname%_win.exe %srcfile%

set GOOS=linux
set GOARCH=amd64
go build -o %bindir%\%binname%_linux %srcfile%

set GOOS=darwin
set GOARCH=amd64
go build -o %bindir%\%binname%_mac %srcfile%
