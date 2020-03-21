#!/bin/bash
set -eu

source scripts/_prepare.sh

cat ${BASE_DIR}/tools/tools.go | grep _ | awk -F'"' '{print $2}'
