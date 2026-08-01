#!/bin/bash

set -e

mkdir -p gen

protoc \
    -I proto \
    --go_out=gen \
    --go_opt=paths=source_relative \
    --go-grpc_out=gen \
    --go-grpc_opt=paths=source_relative \
    proto/shard.proto

echo "done"
