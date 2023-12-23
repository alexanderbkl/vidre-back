#!/usr/bin/env bash

set -e

# see https://docs.docker.com/develop/develop-images/build_enhancements/#to-enable-buildkit-builds
export DOCKER_BUILDKIT=1

if [[ -z $1 ]]; then
    echo "Usage: ${0##*/} [name]" 1>&2
    exit 1
fi

NUMERIC='^[0-9]+$'
BUILD_DATE=$(/bin/date -u +%y%m%d)

echo "Building image 'vidre-backend'";

docker buildx build \
  --pull \
  --no-cache \
  --rm \
  --build-arg BUILD_TAG=$BUILD_DATE \
  -t vidre-backend:latest \
  -f Dockerfile \
  --push .

echo "Done."
