#!/bin/bash
set -eux

source scripts/_prepare.sh

echo "" > ${BASE_DIR}/coverage.txt
for d in $(${MYDIR}/packages.sh); do
  # Note: must use "go" instead of "${LOCAL_GO}" because this script is executed on CI
  go test -race -p 2 -coverprofile=profile.out -covermode=atomic $d

  if [ -f profile.out ]; then
    cat profile.out >> ${BASE_DIR}/coverage.txt
    rm profile.out
  fi
done
