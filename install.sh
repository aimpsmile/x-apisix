#!/bin/bash
# 命令名称
BUILD_NAME=x-apisix
go build -i -o ${GOPATH}/bin/${BUILD_NAME} -ldflags "-w -s"  ./monitor/main.go


AIM_CONFIG=/storage/code/x-apisix/config AIM_ENV=local MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=etcd.service:2379 x-apisix monitor

