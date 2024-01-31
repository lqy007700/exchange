GO := /usr/local/Cellar/go@1.20/1.20.7/bin/go
GOPATH:=$(shell $(GO) env GOPATH)

.PHONY: init
init:
	go get -u google.golang.org/protobuf/proto
	go install github.com/golang/protobuf/protoc-gen-go@latest
	go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest


.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/asset.proto


