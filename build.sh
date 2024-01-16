GOOS=windows GOARCH=amd64 go build -o bin/homeserver-amd64.exe .
GOOS=windows GOARCH=386 go build -o bin/homeserver-386.exe .

GOOS=darwin GOARCH=amd64 go build -o bin/homeserver-amd64-darwin .

GOOS=linux GOARCH=amd64 go build -o bin/homeserver-amd64-linux .
GOOS=linux GOARCH=386 go build -o bin/homeserver-386-linux .

echo=done