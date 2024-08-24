#!/bin/bash -e

go get -u
go mod tidy

if [[ -d releases ]]; then
  rm -rf releases
fi

mkdir releases

GOOS=linux GOARCH=amd64 go build -o releases/app-info-linux-amd64 github.com/rahulkj/app-info

GOOS=darwin GOARCH=amd64 go build -o releases/app-info-darwin-amd64 github.com/rahulkj/app-info

GOOS=windows GOARCH=386 go build -o releases/app-info-windows-amd64.exe github.com/rahulkj/app-info

OS=$(uname)
CF_CLI_EXISTS=$(which cf)
if [[ "${OS}" == "Darwin" && ! -z ${CF_CLI_EXISTS} ]]; then
  cf install-plugin releases/app-info-darwin-amd64 -f
elif [[ "${OS}" == "Linux" && ! -z ${CF_CLI_EXISTS} ]]; then
  cf install-plugin releases/app-info-darwin-amd64 -f
fi