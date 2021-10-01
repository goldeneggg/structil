#!/bin/bash
set -eu

# shellcheck disable=SC1091
source scripts/_prepare.sh

${LOCAL_GO} list ./... | \grep -v 'vendor' | \grep -v '/examples'
