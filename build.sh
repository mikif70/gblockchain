#!/bin/bash

GOOS=windows GOARCH=386 go build -o gocryptochain.exe

CGO_ENABLED=0 GOOS=linux go build -o gocryptochain -a -installsuffix cgo .