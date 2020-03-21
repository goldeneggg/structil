#!/bin/bash
set -eu

echo "$(go list ./... | \grep -v 'vendor' | \grep -v '/examples')"
