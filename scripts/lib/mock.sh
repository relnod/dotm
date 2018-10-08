#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

mocks=(
    "pkg/profile Action"
    "pkg/fileutil Visitor"
)

mock_package="mock"

function dotm::mock::generate() {
    cd "$1/$mock_package"

    for mock in "${mocks[@]}"
    do
        package=$(echo "$mock" | cut -f1 -d " ")
        interface=$(echo "$mock" | cut -f2 -d " ")

        pegomock generate "$BASE_PACKAGE/$package" "$interface" \
            --output "$interface.go" \
            --package "mock"
    done
}

function dotm::mock::diff() {
    for mock in "${mocks[@]}";
    do
        interface=$(echo "$mock" | cut -f2 -d " ")

        local ret=0
        diff "$1/$mock_package/$interface.go" "$2/$mock_package/$interface.go" || ret=$?
        if [[ $ret -ne 0 ]]; then
            echo "$interface is outdated."
            echo "Run 'make update' to update generated files."
            exit 1
        fi
    done
}
