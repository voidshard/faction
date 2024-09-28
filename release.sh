#! /bin/sh

set -ex

APP=faction
VERSION=test

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR

# build, tag the image
docker build -t ${APP}:${VERSION} -f Dockerfile .
docker tag ${APP}:${VERSION} uristmcdwarf/${APP}:${VERSION}

# set latest tag
docker tag ${APP}:${VERSION} uristmcdwarf/${APP}:latest

# push the image
docker push uristmcdwarf/${APP}:${VERSION}
docker push uristmcdwarf/${APP}:latest

cd -
