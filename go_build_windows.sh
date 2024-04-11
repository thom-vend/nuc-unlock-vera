#!/usr/bin/env bash
set -euxo pipefail
rm -rf release
rm -f -- nucunlocker.exe nucunlocker.exe.asc nucunlocker.exe.sha256
GOOS=windows GOARCH=amd64 go build -o nucunlocker.exe main.go
sha256sum nucunlocker.exe |tee nucunlocker.exe.sha256
gpg --armor --output nucunlocker.exe.asc --detach-sign nucunlocker.exe
gpg --verify nucunlocker.exe.asc nucunlocker.exe
mkdir -p release
mv nucunlocker.exe nucunlocker.exe.asc nucunlocker.exe.sha256 release/

echo done