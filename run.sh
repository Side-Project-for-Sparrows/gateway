#!/bin/bash

set -e //immediate exit when error 

brew install go

DEBUG_MODE=false
if [ "$1" == "-d" ]; then
  DEBUG_MODE=true
fi

if $DEBUG_MODE; then
  if ! command -v dlv &> /dev/null; then
    echo "select dbug mode.."
    go install github.com/go-delve/delve/cmd/dlv@latest
    export PATH=$PATH:$(go env GOPATH)/bin
  fi
fi

go install golang.org/x/tools/gopls@latest

echo "dependency install.."
go mod tidy

sh ./etc/tls/generate-local-tls.sh

if $DEBUG_MODE; then
  echo "debug mode init..."
  dlv debug main.go
else
  CERT_DIR=./etc/tls ENV=dev go run main.go
fi

