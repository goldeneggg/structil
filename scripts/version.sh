#!/bin/bash
set -eu

# shellcheck disable=SC1091
source scripts/_prepare.sh

GREP=$(command -v grep)
SED=$(command -v sed)

${GREP} "const VERSION" "${BASE_DIR}"/version.go | ${SED} -e 's/const VERSION = //g' | ${SED} -e 's/\"//g'
