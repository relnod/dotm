#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source "${ROOT}/hack/lib/gotools.sh"

echo "Verifying gofmt"

dotm::fmt::run "${TMP_ROOT}"
dotm::fmt::diff "${ROOT}" "${TMP_ROOT}"

echo "Done verifying gofmt"
