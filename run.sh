#!/usr/bin/env bash

#代理设置
go env -w GOPROXY="https://goproxy.cn,direct"

go run main.go -c ./configs/local/config.toml -e local