.PHONY: install build image clean run

CWD = $(shell pwd)
GOPATH = ${CWD}/build

install:
	[ -d ${GOPATH} ] || mkdir -p ${GOPATH}
	GOPATH=${GOPATH} go get -d -v ./...

build: install
	GOPATH=${GOPATH} go build

clean:
	rm -fr ${GOPATH}
	rm hazrd

image: build
	docker build -t hazrd-image .
	docker image tag hazrd-image:latest hazrd-image:v1

run:
	docker run -p 8080:8080 hazrd-image:latest hazrd

