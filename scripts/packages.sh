#!/bin/bash
set -eu

source scripts/_prepare.sh

echo "$(${LOCAL_GO} list ./... | \grep -v 'vendor' | \grep -v '/examples')"
