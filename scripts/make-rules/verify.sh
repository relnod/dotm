#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

mkdir -p "${TMP_ROOT}"
rsync -vr "${ROOT}" "${TMP_ROOT}" --exclude .git --exclude vendor

# currently not working
$ROOT/scripts/verify-generated-mocks.sh
$ROOT/scripts/verify-fmt.sh
$ROOT/scripts/verify-lint.sh
