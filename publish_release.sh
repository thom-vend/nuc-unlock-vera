#!/usr/bin/env bash
set -euxo pipefail
new_tag=${1}
comment=${2}
./go_build_windows.sh
git tag -a "$new_tag" -m "$comment"
git push origin "$new_tag"
gh release create --generate-notes --latest "$new_tag" ./release/*
