#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

mkdir -p "${TMP_ROOT}"
rsync -avr "${ROOT}" "${TMP_ROOT}" --exclude .git

# TMP_ROOT is used to generate files in a seperate directory
export TMP_ROOT=${TMP_ROOT}/dotm

# currently not working
# $ROOT/hack/verify-generated-mocks.sh
$ROOT/hack/verify-fmt.sh
$ROOT/hack/verify-lint.sh
