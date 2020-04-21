
BUILD_NAME:=x-apisix
EXEC_ENV:=AIM_CONFIG=/storage/code/x-apisix/config AIM_ENV=local MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=etcd.service:2379

.PHONY: install monitor checkall lint

install:
	go build -i -o ${GOPATH}/bin/${BUILD_NAME} -ldflags "-w -s"  ./monitor/main.go ./monitor/plugins.go

monitor:
	${EXEC_ENV} x-apisix monitor

checkall:
	${EXEC_ENV} x-apisix checkall

lint:
	 golangci-lint run  -c .golangci.yml --deadline=1000s


