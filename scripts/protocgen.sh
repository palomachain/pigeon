#!/usr/bin/env bash

set -e

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  proto_files=$(find "${dir}" -maxdepth 1 -name '*.proto')
  for file in $proto_files; do
    protoc \
    -I "proto" \
    -I $(dirname $file) \
    -I $(dirname $(dirname $file)) \
    -I "third_party_proto/proto" \
    --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
    $file
  done

done

cp -r ./github.com/palomachain/* types

rm -rf ./github.com
