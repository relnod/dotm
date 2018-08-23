#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

${ROOT}/hack/update-generated-mocks.sh
