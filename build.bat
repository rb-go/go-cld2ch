@echo off
TITLE CollectD2Clickhouse Building
cd %~dp0
set /p Build=<version
for /f %%x in ('wmic path win32_utctime get /format:list ^| findstr "="') do set %%x
set today=%Year%-%Month%-%Day%

echo Getting required packages
go get ./...
Rem goimports -d .
Rem golint ./...

set GOARCH=amd64

echo Building for Windows
set GOOS=windows
go build --ldflags "-s -w -X 'main.version=%Build%' -X 'main.build=release-%Build%' -X 'main.buildDate=%today%'" -o %~dp0/out/cld2ch_64.exe

echo Building for Linux
set GOOS=linux
go build --ldflags "-s -w -X 'main.version=%Build%' -X 'main.build=release-%Build%' -X 'main.buildDate=%today%'" -o %~dp0/out/cld2ch_l_64

echo Building for Mac OS
set GOOS=darwin
go build --ldflags "-s -w -X 'main.version=%Build%' -X 'main.build=release-%Build%' -X 'main.buildDate=%today%'" -o %~dp0/out/cld2ch_darwin