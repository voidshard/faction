#! /bin/bash

set -ex

protoc -I=./proto --go_out=./ ./proto/*.proto
