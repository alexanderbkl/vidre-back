#!/usr/bin/env bash

if [[ -z $1 ]] || [[ -z $2 ]]; then
  echo "Usage: ${0##*/} [debug|race|static|prod] [filename]" 1>&2
  exit 1
fi

set -e

BUILD_OS=$(uname -s)
BUILD_DATE=$(date -u +%y%m%d)
BUILD_VERSION=$(git describe --always)
BUILD_TAG=${BUILD_DATE}-${BUILD_VERSION}
BUILD_ID=${BUILD_TAG}-${BUILD_OS}
BUILD_BIN=${2:-hello}
GO_BIN=${GO_BIN:-go}
GO_VER=$($GO_BIN version)

echo "Building Hello ${BUILD_ID} ($1)..."

if [[ $1 == "debug" ]]; then
  BUILD_CMD=("$GO_BIN" build -tags=debug -ldflags "-X main.version=${BUILD_ID}-DEBUG" -o "build/${BUILD_BIN}" cmd/main.go)
elif [[ $1 == "race" ]]; then
  BUILD_CMD=("$GO_BIN" build -tags=debug -race -ldflags "-X main.version=${BUILD_ID}-DEBUG" -o "build/${BUILD_BIN}" cmd/main.go)
elif [[ $1 == "static" ]]; then
  BUILD_CMD=("$GO_BIN" build -a -v -ldflags "-linkmode external -extldflags \"-static -L /usr/lib -ltensorflow\" -s -w -X main.version=${BUILD_ID}" -o "build/${BUILD_BIN}" cmd/main.go)
else
  BUILD_CMD=("$GO_BIN" build -ldflags "-extldflags \"-Wl,-rpath -Wl,\$ORIGIN/../lib\" -s -w -X main.version=${BUILD_ID}" -o "build/${BUILD_BIN}" cmd/main.go)
fi

# Build app binary.
echo "=> compiling \"$BUILD_BIN\" with \"${GO_VER}\""
echo "=> ${BUILD_CMD[*]}"
"${BUILD_CMD[@]}"

# Display binary size.
du -h "build/${BUILD_BIN}"

echo "Done."