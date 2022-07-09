SHELL := /bin/bash

make_ser_cert:
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-server/CN=server.grpc.io" \
    -key server.key -out server.crt

make_client_cert:
	openssl genrsa -out client.key 2048
	openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-client/CN=client.grpc.io" \
    -key client.key -out client.crt

gen_service:
	protoc --go-grpc_out=require_unimplemented_servers=false:. *.proto

test:
	# go test -v -coverprofile=c.out -coverpkg ./... ./tests/...
	# go tool cover -html=c.out -o coverage.html
	go test -v -run=Test ./prototype

.PHONY: test