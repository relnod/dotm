#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

list=(
    "linux amd64"
    "linux arm64"
)

mkdir -p $ROOT/artifacts

for entry in "${list[@]}"
do
    target=$(echo "${entry}" | cut -f1 -d " ")
    goarch=$(echo "${entry}" | cut -f2 -d " ")
    cd $ROOT && TARGET=${target} GOARCH=${goarch} make build
    cd $ROOT/build && tar -czf "${ROOT}/artifacts/dotm_${VERSION}_${target}_${goarch}.tar.gz" dotm
done
