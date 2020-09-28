#!/usr/bin/env bash

set -o xtrace

BIN_NAME=./terraform-provider-scaleway_v1.14.0_x4

DARWIN_DIR_PATH=$TF_DIR/plugins/darwin_amd64
LINUX_DIR_PATH=$TF_DIR/plugins/linux_amd64

DARWIN_BIN_PATH=$DARWIN_DIR_PATH/$BIN_NAME
LINUX_BIN_PATH=$LINUX_DIR_PATH/$BIN_NAME

# Build
env GOOS=darwin GOARCH=amd64 go build -o $DARWIN_BIN_PATH .
env GOOS=linux GOARCH=amd64 go build -o $LINUX_BIN_PATH .

# Shasum
cat $DARWIN_DIR_PATH/lock.json | jq ".scaleway = \"$(shasum -a 256 $DARWIN_BIN_PATH | awk '{ print $1 }')\"" > $DARWIN_DIR_PATH/lock.json
cat $LINUX_DIR_PATH/lock.json | jq ".scaleway = \"$(shasum -a 256 $LINUX_BIN_PATH | awk '{ print $1 }')\"" > $LINUX_DIR_PATH/lock.json
