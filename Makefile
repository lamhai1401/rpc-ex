SHELL := /bin/bash

gen_service:
	protoc --go-grpc_out=require_unimplemented_servers=false:. *.proto

test:
	# go test -v -coverprofile=c.out -coverpkg ./... ./tests/...
	# go tool cover -html=c.out -o coverage.html
	go test -v -run=Test ./prototype

.PHONY: test