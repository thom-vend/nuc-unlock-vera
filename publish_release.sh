#!/usr/bin/env bash
set -euxo pipefail
git tag -l
new_tag=${1:-}
if [[ -z "$new_tag" ]] || [[ $# -ne 1 ]]; then
    echo "Error: usage :$0 <tag>"
    exit 1
fi
go fmt
./go_tests.sh
./integration_tests.sh
./go_build_release.sh
git tag -a "$new_tag" -m 'Release tag'
git push origin "$new_tag"
gh release create --generate-notes --latest "$new_tag" ./release/*
