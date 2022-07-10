# SHELL := /bin/bash

# build proto file
proto:
	protoc \
		--go_out=. \
		--go_opt=paths=import \
		--go-grpc_out=require_unimplemented_servers=false:. \
		--go-grpc_opt=paths=import \
		third_party/google/api/annotations.proto \
		third_party/google/api/http.proto \
		security/rest_service.proto

# build gateway endpoint
gateway:
	protoc \
		--go_out=$(GOPATH)/src \
		--go_opt=paths=import \
		--grpc-gateway_out=logtostderr=true:$(GOPATH)/src/livestreaming/rpc-ex \
		security/rest_service.proto

# must run this cmd
compile: proto gateway

valid:
	# protoc --govalidators_out=. --go-grpc_out=require_unimplemented_servers=false:. *.proto
	protoc \
    --proto_path=. \
	--proto_path=${GOPATH}/src \
    --proto_path=${GOPATH}/pkg/mod/github.com/mwitkow/go-proto-validators@v0.3.2/ \
    --proto_path=${GOPATH}/pkg/mod/google.golang.org/protobuf@v1.28.0/ \
	--go-grpc_out=require_unimplemented_servers=false:. *.proto


	protoc \
	--proto_path=./ \
	--go_out=./user/*.proto \
	--go-grpc_out=require_unimplemented_servers=false:. *.proto

build_folder:
	protoc \
	--go_out=paths=source_relative:.\
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

init:
	go get \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
	google.golang.org/protobuf/cmd/protoc-gen-go \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: test compile