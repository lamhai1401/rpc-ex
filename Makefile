SHELL := /bin/bash

make_valid:
	# protoc --govalidators_out=. --go-grpc_out=require_unimplemented_servers=false:. *.proto
	protoc \
    --proto_path=. \
	--proto_path=${GOPATH}/src \
    --proto_path=${GOPATH}/pkg/mod/github.com/mwitkow/go-proto-validators@v0.3.2/ \
    --proto_path=${GOPATH}/pkg/mod/google.golang.org/protobuf@v1.28.0/ \
	--go_out=. \
	--go-grpc_out=require_unimplemented_servers=false:. *.proto

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