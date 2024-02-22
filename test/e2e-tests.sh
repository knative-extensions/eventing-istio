#!/usr/bin/env bash

set -euo pipefail

source $(dirname $0)/e2e-common.sh

if ! ${SKIP_INITIALIZE}; then
  initialize "$@" --cluster-version="${K8S_CLUSTER_VERSION}" --num-nodes 3
fi

run_eventing_core_tests || fail_test
run_eventing_kafka_broker_tests || fail_test
