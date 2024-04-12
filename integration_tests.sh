#!/usr/bin/env bash
# Simple integration test
set -eEuo pipefail
srvpid=""
output=$(mktemp)
function exit_trap {
    rm -f -- "$output"
    kill "$srvpid"
    echo cleaned
}
trap exit_trap EXIT
pushd tests
./start-python-server.sh &
srvpid=$!
popd
sleep 0.1
go run main.go -m unlock 2>&1 |tee "$output"
# Check output
if ! grep -Fq 'word 1 word 2 word 3' "$output"; then
    echo "failed ❌ partial payload content not found, should be outputed by the executed cmd"
    exit 1
fi
if ! grep -Fq 'NUC unlocked' "$output"; then
    echo "failed ❌ expected sucess msg not found"
    exit 1
fi
echo "Test passed ✅"
