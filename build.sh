mkdir -p bin

go mod tidy

GOOS=windows GOARCH=amd64 go build -o bin/homeserver-amd64.exe .
GOOS=windows GOARCH=386 go build -o bin/homeserver-386.exe .

GOOS=darwin GOARCH=amd64 go build -o bin/homeserver-amd64-darwin .

GOOS=linux GOARCH=amd64 go build -o bin/homeserver-amd64-linux .
GOOS=linux GOARCH=386 go build -o bin/homeserver-386-linux .

cp -r web bin/
touch bin/.env

read -p "Enter the path of directory you want to serve, press interto skip " path
echo "HOME_PATH=$path" >> bin/.env

read -p "Enter the your discord bot token, press inter to skip " discord
echo "DISCORD_BOT_TOKEN=$discord" >> bin/.env

read -p "Enter the your telegram bot token, press inter to skip " telegram
echo "TELEGRAM_BOT_TOKEN=$telegram" >> bin/.env


echo=done
