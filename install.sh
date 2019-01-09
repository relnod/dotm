#!/bin/sh

set -e

get_latest_release() {
  curl --silent "https://api.github.com/repos/relnod/dotm/releases/latest" |
    grep '"tag_name":' |
    sed -E 's/.*"([^"]+)".*/\1/'
}

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
if [ -x "$(command -v dotm)" ] && [ "$version" = "$(dotm --version | cut -d ' ' -f3)" ]; then
    echo "dotm is already installed at the latest version ($version)"
    exit
fi

name="dotm_${version}_${target}_${goarch}.tar.gz"

echo "Downloading dotm binary at version ${version}"
curl --silent -L "https://github.com/relnod/dotm/releases/download/${version}/${name}" -o "/tmp/${name}"

install_dir="/usr/local/bin"
bash_completion_dir="/etc/bash_completion.d"
if [ "$1" = "--user" ]; then
    install_dir="$HOME/.local/bin"
    bash_completion_dir="$HOME/.bash_completion.d"
fi
if [ ! -d "$install_dir" ]; then
    mkdir -p "$install_dir"
fi
if [ ! -d "$bash_completion_dir" ]; then
    mkdir -p "$bash_completion_dir"
fi

echo "Installing dotm to $install_dir"
if [ -f "$install_dir/dotm" ]; then
    rm -f "$install_dir/dotm"
fi
tar -C "$install_dir" -xzf "/tmp/${name}"

echo "Generating bash completions at $bash_completion_dir"
dotm --genCompletions > "$bash_completion_dir/dotm"
