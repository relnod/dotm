#!/bin/bash

base_package="github.com/relnod/dotm"


mocks=(
    "dotfiles Action"
    "internal/util/file Visitor"
)

cd internal/mock

for mock in "${mocks[@]}"
do
    package=$(echo "$mock" | cut -f1 -d " ")
    interface=$(echo "$mock" | cut -f2 -d " ")

    pegomock generate "$base_package/$package" "$interface" --output "$interface.go" --package "mock"
done

# pegomock generate "$package/dotfiles" Action --output=internal/mocks/
