#!/bin/sh

set -e

install_dir="/usr/local/bin"
bash_completion_dir="/etc/bash_completion.d"
if [ "$1" = "--user" ]; then
    install_dir="$HOME/.local/bin"
    bash_completion_dir="$HOME/.bash_completion.d"
fi

get_latest_release() {
  curl --silent "https://api.github.com/repos/relnod/dotm/releases/latest" |
    grep '"tag_name":' |
    sed -E 's/.*"([^"]+)".*/\1/'
}

target=""
goarch=""

case $(uname) in
    "Linux")target=Linux;;
    "Darwin")target=Darwin;;
    *)
        echo "Target $(uname) is not supported"
        exit 1
        ;;
    esac

case $(uname -m) in
    "x86_64")goarch=x86_64;;
    "i386")goarch=i386;;
    *)
        echo "Arch $(uname -a) is not supported"
        exit 1
        ;;
    esac

version=$(get_latest_release)
current_version=$(dotm --version | cut -d ' ' -f3)
if [ -x "$(command -v dotm)" ] && [ "$version" = "$current_version" ]; then
    echo "dotm is already installed at the latest version ($version)"
    exit
fi

# strip the "v" from the version
rawversion="$(echo "$version" | cut -d "v" -f 2)"

name="dotm_${rawversion}_${target}_${goarch}.tar.gz"

echo "Downloading dotm binary at version ${version}"
curl --silent -L "https://github.com/relnod/dotm/releases/download/${version}/${name}" -o "/tmp/${name}"

[ ! -d "$install_dir" ] && mkdir -p "$install_dir"
[ ! -d "$bash_completion_dir" ] && mkdir -p "$bash_completion_dir"

echo "Installing dotm to $install_dir"
if [ -f "$install_dir/dotm" ]; then
    rm -f "$install_dir/dotm"
fi

tar -C "$install_dir" -xzf "/tmp/${name}" dotm

echo "Generating bash completions at $bash_completion_dir"
dotm --genCompletions > "$bash_completion_dir/dotm"

if [ ! -z "$current_version" ]; then
    echo "Running 'dotm fix'"
    dotm fix
fi
