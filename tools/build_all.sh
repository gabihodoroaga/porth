#!/bin/bash
set -e
# clean 

rm -rf ./bin/*

# build linux
echo "build linux"
env GOOS=linux GOARCH=amd64 go build -o bin/porth-server-linux64 ./server
env GOOS=linux GOARCH=amd64 go build -o bin/porth-operator-linux64 ./operator
env GOOS=linux GOARCH=amd64 go build -o bin/porth-client-linux64 ./client

# build windows
echo "build windows"
env GOOS=windows GOARCH=amd64 go build -o bin/porth-server-win64.exe ./server
env GOOS=windows GOARCH=amd64 go build -o bin/porth-operator-win64.exe ./operator
env GOOS=windows GOARCH=amd64 go build -o bin/porth-client-win64.exe ./client

# build darwin
echo "build darwin"
env GOOS=darwin GOARCH=amd64 go build -o bin/porth-server-darwin64 ./server
env GOOS=darwin GOARCH=amd64 go build -o bin/porth-operator-darwin64 ./operator
env GOOS=darwin GOARCH=amd64 go build -o bin/porth-client-darwin64 ./client

