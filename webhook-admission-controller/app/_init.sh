#!/bin/bash
if [ $# -eq 0 ]
  then
    echo "No module name passed..."
    exit 0
fi
export GO111MODULE="on"
go mod init $1
go mod tidy

go get -u github.com/gorilla/mux