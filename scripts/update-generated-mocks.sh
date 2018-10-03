#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source "${ROOT}/scripts/lib/mock.sh"

echo "Generating mocks"

dotm::mock::generate "${ROOT}"

echo "Done generating mocks"
