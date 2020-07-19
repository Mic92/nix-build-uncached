#!/usr/bin/env bash

set -eu -o pipefail

escape() { perl -p -e 's/%/%25/;s/\r/%0D/;s/\n/%0A/'; }
ref="$1"
tag_name=${ref#refs/tags/}
release_name=$(git tag -l --format='%(subject)' "$tag_name"  | escape)
body=$(git tag -l --format='%(contents:body)' "$tag_name" | escape)

echo "Release title: ${release_name}"
echo "Release description: ${body}"

echo "::set-output name=release_name::$release_name"
echo "::set-output name=body::$body"
