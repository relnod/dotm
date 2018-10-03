#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

${ROOT}/scripts/update-generated-mocks.sh
${ROOT}/scripts/update-fmt.sh
