#!/usr/bin/env bash

set -euo pipefail

source $(dirname $0)/e2e-common.sh

if ! ${SKIP_INITIALIZE}; then
  initialize $@ --skip-istio-addon --num-nodes 3
fi

run_eventing_core_tests || fail_test
run_eventing_kafka_broker_tests || fail_test
