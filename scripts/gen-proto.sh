#!/bin/bash
CURRENT_DIR=$(pwd)

for module in $(find $CURRENT_DIR/book-store-protos/* -type d); do
    protoc -I /usr/local/include \
           -I $GOPATH/src/github.com/gogo/protobuf/gogoproto \
           -I $CURRENT_DIR/book-store-protos/ \
            --gofast_out=plugins=grpc:$CURRENT_DIR/genproto/ \
            $module/*.proto;
done;

for module in $(find $CURRENT_DIR/genproto/* -type d); do
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i "" -e "s/,omitempty//g" $module/*.go
  else
    sed -i -e "s/,omitempty//g" $module/*.go
  fi
done;