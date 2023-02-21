#!/usr/bin/env bash

REPO_ROOT_DIR=$(dirname $(realpath $0))/..

echo "$REPO_ROOT_DIR"

export SKIP_INITIALIZE=${SKIP_INITIALIZE:-false}
export SYSTEM_NAMESPACE=${SYSTEM_NAMESPACE:-"knative-eventing"}
export ISTIO_NAMESPACE=${ISTIO_NAMESPACE:-"istio-system"}

source "${REPO_ROOT_DIR}"/vendor/knative.dev/hack/e2e-tests.sh

git submodule update --init --recursive

function knative_setup() {
  "${REPO_ROOT_DIR}"/hack/update-istio.sh || return $?
  "${REPO_ROOT_DIR}"/hack/install-dependencies.sh || return $?
  kubectl apply -n "${ISTIO_NAMESPACE}" -Rf "${REPO_ROOT_DIR}"/test/config
  kubectl apply -n "${SYSTEM_NAMESPACE}" -Rf "${REPO_ROOT_DIR}"/test/config
}

function run_eventing_core_tests() {
  pushd "${REPO_ROOT_DIR}"/third_party/eventing || return $?

  go_test_e2e \
    -timeout=1h \
    -run TestPingSource \
    ./test/rekt/ \
    --istio.enabled=true || return $?

  pod
}

function run_eventing_kafka_broker_tests() {
  pushd "${REPO_ROOT_DIR}"/third_party/eventing-kafka-broker || return $?

  BROKER_CLASS=Kafka go_test_e2e \
    -timeout=1h \
    -run TestBrokerConformance \
    ./test/e2e_new/... \
    --istio.enabled=true || return $?

  BROKER_CLASS=KafkaNamespaced go_test_e2e \
    -timeout=1h \
    -run TestBrokerConformance \
    ./test/e2e_new/... \
    --istio.enabled=true || return $?

  BROKER_CLASS=Kafka go_test_e2e \
    -timeout=1h \
    -run TestKafkaSource \
    ./test/e2e_new/... \
    --istio.enabled=true || return $?

  pod
}
