#!/usr/bin/env bash

set -eux -o pipefail

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
store=$(mktemp -p "$dir" -d)
outputs=$(mktemp -p "$dir" -d)
trap "{ chmod -R +w $store; rm -rf $store $outputs; }" EXIT

go run $dir/.. -flags "--store '$store'" -build-flags "-o '$outputs/result1'" $dir/test.nix
go run $dir/.. -flags "--store '$store'" -build-flags "-o '$outputs/result2'" $dir/test2.nix
count=$(find "$outputs" -type l | wc -l)

if [[ $count != "2" ]]; then
    echo "expect two package to be built, got $count"
fi
