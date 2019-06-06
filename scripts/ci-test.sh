#!/bin/bash
set -eux

source scripts/_prepare.sh

echo "" > coverage.txt
for d in $(${MYDIR}/_packages.sh); do
  go test -race -coverprofile=profile.out -covermode=atomic $d
  if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
  fi
done
