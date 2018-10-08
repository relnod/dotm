#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

get_latest_release() {
  curl --silent "https://api.github.com/repos/relnod/dotm/releases/latest" |
    grep '"tag_name":' |
    sed -E 's/.*"([^"]+)".*/\1/'
}

version="0.1.0"
target=""
goarch=""

case $(uname) in
    "Linux")target=linux;;
    *)
        echo "Target $(uname) is not supported"
        exit 1
        ;;
    esac

case $(uname -m) in
    "x86_64")goarch=amd64;;
    *)
        echo "Arch $(uname -a) is not supported"
        exit 1
        ;;
    esac

version=$(get_latest_release)

name="dotm_${version}_${target}_${goarch}.tar.gz"

echo "Downloading dotm binary at version ${version}"
curl --silent -L "https://github.com/relnod/dotm/releases/download/${version}/${name}" -o "/tmp/${name}"

echo "Installing dotm to /usr/local/bin"
rm -f /usr/local/bin/dotm
tar -C /usr/local/bin -xzf "/tmp/${name}"

echo "Generating bash completions"
dotm --genCompletions > /etc/bash_completion.d/dotm
