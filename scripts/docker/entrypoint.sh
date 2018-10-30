#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

useradd $USER
su - $USER
export HOME=/home/$USER

/repo/build/dotm $@
