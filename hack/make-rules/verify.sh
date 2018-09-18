#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

mkdir -p "${TMP_ROOT}"
rsync -vr "${ROOT}" "${TMP_ROOT}" --exclude .git --exclude vendor

# currently not working
# $ROOT/hack/verify-generated-mocks.sh
$ROOT/hack/verify-fmt.sh
$ROOT/hack/verify-lint.sh
