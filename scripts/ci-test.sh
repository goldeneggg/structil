#!/bin/bash
set -eux

source scripts/_prepare.sh

echo "" > ${BASE_DIR}/coverage.txt
for d in $(${MYDIR}/_packages.sh); do
  go test -race -parallel 2 -coverprofile=profile.out -covermode=atomic $d
  if [ -f profile.out ]; then
    cat profile.out >> ${BASE_DIR}/coverage.txt
    rm profile.out
  fi
done
