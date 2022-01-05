#!/bin/bash

CURRENT_DIR=$(pwd)

protoc -I /usr/local/include \
       -I $GOPATH/src/github.com/gogo/protobuf/gogoproto \
       -I $CURRENT_DIR/book-store-protos/ \
        --gofast_out=plugins=grpc:$CURRENT_DIR/genproto/ \
        $CURRENT_DIR/book-store-protos/*.proto;

if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i "" -e "s/,omitempty//g" $CURRENT_DIR/genproto/*.go
  else
    sed -i -e "s/,omitempty//g" $CURRENT_DIR/genproto/*.go
fi