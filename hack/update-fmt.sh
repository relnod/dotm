#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source "${ROOT}/hack/lib/gotools.sh"

echo "Updating gofmt"

dotm::fmt::run "${ROOT}"

echo "Done updating gofmt"
