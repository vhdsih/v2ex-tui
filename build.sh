#!/bin/bash

# Windows
#GOOS=windows GOARCH=amd64 go build -o bin/windows_v2ex_tui.exe ./cmd/v2ex/main.go

# Mac Intel
GOOS=darwin GOARCH=amd64 go build -o bin/mac_v2ex_tui ./cmd/v2ex/main.go

# Linux
#GOOS=linux GOARCH=amd64 go build -o bin/linux_v2ex_tui ./cmd/v2ex/main.go
