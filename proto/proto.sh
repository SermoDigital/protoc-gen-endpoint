#!/usr/bin/env bash

set -euo pipefail

protoc -I. \
--go_out=:. \
*.proto

# Rename records a proto rename from $1->$2 for file $3
function rename() {
    sed -i "s,\"$1\",\"$2\" // from $1," $3
}

rename "google/protobuf" "github.com/golang/protobuf/protoc-gen-go/descriptor" "annotations.pb.go"
