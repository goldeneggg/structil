#!/bin/bash
MYDIR=$(cd "$(dirname "${0}")" && pwd)

# shellcheck disable=SC2034
BASE_DIR=${MYDIR}/..

# shellcheck disable=SC2034
FULL_PACKAGE=github.com/goldeneggg/structil

# shellcheck disable=SC2034
LOCAL_GO=go${LOCAL_GO_VERSION:-}
