#! /bin/sh

# Simple script to build a local image for use with minikube

set -ex

VERSION=test
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR

eval $(minikube -p minikube docker-env)
docker build -t faction:${VERSION} -f Dockerfile .

cd -
