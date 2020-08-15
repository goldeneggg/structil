#!/bin/bash
MYDIR=$(cd $(dirname $0) && pwd)
BASE_DIR=${MYDIR}/..
FULL_PACKAGE=github.com/goldeneggg/structil
LOCAL_GO=go${LOCAL_GO_VERSION:-}
