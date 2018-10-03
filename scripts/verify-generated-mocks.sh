#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source "${ROOT}/scripts/lib/mock.sh"
source "${ROOT}/scripts/lib/gotools.sh"

echo "Verifying mocks"

dotm::mock::generate "${TMP_ROOT}"
dotm::fmt::run "${TMP_ROOT}"
dotm::mock::diff "${ROOT}" "${TMP_ROOT}"

echo "Done verifying mocks"
