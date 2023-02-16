#!/usr/bin/env bash

set -euo pipefail

source $(dirname $0)/e2e-common.sh

if ! ${SKIP_INITIALIZE}; then
  initialize $@ --skip-istio-addon --min-nodes 3 --max-nodes 3
fi
