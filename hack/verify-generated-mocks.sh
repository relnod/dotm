#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source "${ROOT}/hack/lib/mock.sh"

echo "Verifying mocks"

dotm::mock::generate "${TMP_ROOT}"
dotm::mock::diff "${ROOT}" "${TMP_ROOT}"

echo "Done verifying mocks"
