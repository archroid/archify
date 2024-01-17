@echo off
mkdir bin

go mod tidy

set GOOS=windows
set GOARCH=amd64
go build -o bin/archify-amd64.exe .
set GOARCH=386
go build -o bin/archify-386.exe .

set GOOS=darwin
set GOARCH=amd64
go build -o bin/archify-amd64-darwin .

set GOOS=linux
set GOARCH=amd64
go build -o bin/archify-amd64-linux .
set GOARCH=386
go build -o bin/archify-386-linux .

xcopy /E /I templates bin\templates
xcopy /E /I static bin\static
type nul > bin\.env

set /p path="Enter the path of directory you want to serve, press enter to skip: "
echo HOME_PATH=%path% >> bin\.env

set /p discord="Enter your discord bot token, press enter to skip: "
echo DISCORD_BOT_TOKEN=%discord% >> bin\.env

set /p telegram="Enter your telegram bot token, press enter to skip: "
echo TELEGRAM_BOT_TOKEN=%telegram% >> bin\.env

echo done