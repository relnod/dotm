#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source "${ROOT}/hack/lib/gotools.sh"

echo "Verifying golint"

dotm::lint::run "${ROOT}"

echo "Done verifying golint"
