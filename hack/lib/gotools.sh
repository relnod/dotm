#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOFILES=$(find . -iname '*.go' -type f | grep -v "vendor" | cut -c 3-)

# Checks if goimports is installed.
function dotm::fmt::check() {
    if [[ -z $(which goimports) ]]; then
        echo ""
        echo "'goimports' is not installed."
        echo "Run 'go get golang.org/x/tools/cmd/goimports' to install it."
        echo ""
        exit 1;
    fi
}

# Run gofmt on all $1/GOFILES.
function dotm::fmt::run() {
    dotm::fmt::check

    for file in $GOFILES; do
        goimports -w "$1/$file"
    done
}

# Creates a diff fo all files to check if they have been gofmt'ed.
function dotm::fmt::diff() {
    for file in $GOFILES; do
        local ret=0
        diff "$1/$file" "$2/$file" || ret=$?
        if [[ $ret -ne 0 ]]; then
            echo "$file is outdated."
            echo "Run 'make update' to update generated files."
            exit 1
        fi
    done
}

# Checks if golint is installed.
function dotm::lint::check() {
    if [[ -z $(which golint) ]]; then
        echo ""
        echo "'golint' is not installed."
        echo "Run 'go get golang.org/x/lint/golint' to install it."
        echo ""
        exit 1;
    fi
}

# Run golint on all $1/GOFILES.
function dotm::lint::run() {
    dotm::lint::check

    local failed=false
    for file in $GOFILES; do
        bla=$(golint "$file")
        if [[ ! -z "$bla" ]]; then
            echo "$bla"
            failed=true
        fi
    done

    if [[ $failed = true ]]; then
        echo ""
        echo "golint faild!"
        exit 1
    fi
}
