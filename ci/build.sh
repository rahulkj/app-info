#!/bin/bash -e

if [[ -d releases ]]; then
  rm -rf releases
fi

mkdir releases

GOOS=linux GOARCH=amd64 go build -o releases/app-info-linux-amd64 github.com/rahulkj/app-info

GOOS=darwin GOARCH=amd64 go build -o releases/app-info-darwin-amd64 github.com/rahulkj/app-info

GOOS=windows GOARCH=386 go build -o releases/app-info-windows-amd64.exe github.com/rahulkj/app-info
