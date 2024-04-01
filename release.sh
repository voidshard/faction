#! /bin/sh

set -ex

VERSION=0.0.1

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR

docker build -t faction:${VERSION} -f Dockerfile .
docker tag faction:${VERSION} uristmcdwarf/faction:${VERSION}
docker push uristmcdwarf/faction:${VERSION}

cd -
