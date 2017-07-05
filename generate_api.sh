#!/usr/bin/env bash
export PATH=$PATH:$GOPATH/bin
cd $GOPATH/src/github.com/foofilers/cfhd
echo "generate users rpc"
protoc -I rpc rpc/auth/auth.proto --go_out=plugins=grpc:rpc
protoc -I rpc rpc/users/users.proto --go_out=plugins=grpc:rpc
protoc -I rpc rpc/applications/applications.proto --go_out=plugins=grpc:rpc
