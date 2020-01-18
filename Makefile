.PHONY: install build image clean run

CWD = $(shell pwd)
GOPATH = ${CWD}/build

install:
	[ -d ${GOPATH} ] || mkdir -p ${GOPATH}
	GOPATH=${GOPATH} go get -d -v ./...

build: install
	GOPATH=${GOPATH} go build

clean:
	rm -fr ${GOPATH} | true
	rm hazrd | true
	docker ps -a | grep "hazrd-image:latest" | awk '{print $$1 }' | xargs -I {} docker stop {} | true
	docker ps -a | grep "hazrd-image:latest" | awk '{print $$1 }' | xargs -I {} docker rm {} --force | true
	docker rmi hazrd-image:latest --force | true

image: build
	docker build -t hazrd-image .

run:
	docker run -d -p 8080:8080 hazrd-image:latest hazrd

