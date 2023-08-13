#!/bin/bash

GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o app cmd/main.go
