#!/usr/bin/env bash

source $(dirname $0)/e2e-common.sh

if ! ${SKIP_INITIALIZE}; then
  initialize $@ --skip-istio-addon --min-nodes 3 --max-nodes 3
fi

# TODO(approvers): replace with real tests
true
