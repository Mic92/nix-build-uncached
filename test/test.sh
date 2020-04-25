#!/usr/bin/env bash

set -eux -o pipefail

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
store=$(mktemp -p "$dir" -d)
outputs=$(mktemp -p "$dir" -d)
trap "{ chmod -R +w $store; rm -rf $store $outputs; }" EXIT
go run $dir/.. -args "--store '$store'" -build-args "-o '$outputs/result'" $dir/test.nix
count=$(find "$outputs" -type l | wc -l)

if [[ $count != "1" ]]; then
    echo "expect one package to be built, got $count"
fi
