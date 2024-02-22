#!/usr/bin/env bash

repo_root_dir=$(dirname $(realpath $0))/..

echo "$repo_root_dir"

export SKIP_INITIALIZE=${SKIP_INITIALIZE:-false}
export SYSTEM_NAMESPACE=${SYSTEM_NAMESPACE:-"knative-eventing"}
export ISTIO_NAMESPACE=${ISTIO_NAMESPACE:-"istio-system"}
export KAFKA_BROKER_TEMPLATES=${KAFKA_BROKER_TEMPLATES:-"${repo_root_dir}/test/e2e/templates/kafka-broker"}
export KAFKA_NAMESPACED_BROKER_TEMPLATES=${KAFKA_BROKER_TEMPLATES:-"${repo_root_dir}/test/e2e/templates/kafka-namespaced-broker"}

readonly K8S_CLUSTER_VERSION=1.28

source "${repo_root_dir}"/vendor/knative.dev/hack/e2e-tests.sh

function knative_setup() {
  git submodule update --init --recursive
  "${repo_root_dir}"/hack/install-dependencies.sh || return $?
  "${repo_root_dir}"/hack/install.sh || return $?

  wait_until_pods_running "knative-eventing" || return $?
}

function run_eventing_core_tests() {
  pushd "${repo_root_dir}"/third_party/eventing || return $?

  BROKER_TEMPLATES="${KAFKA_BROKER_TEMPLATES}" go_test_e2e \
    -timeout=1h \
    -parallel=12 \
    -run TestPingSource \
    ./test/rekt/ \
    --istio.enabled=true || return $?

  BROKER_TEMPLATES="${KAFKA_BROKER_TEMPLATES}" go_test_e2e \
    -timeout=1h \
    -parallel=12 \
    -run TestBrokerConformance \
    ./test/rekt/ \
    --istio.enabled=true || return $?

  BROKER_TEMPLATES="${KAFKA_NAMESPACED_BROKER_TEMPLATES}" go_test_e2e \
    -timeout=1h \
    -parallel=12 \
    -run TestBrokerConformance \
    ./test/rekt/ \
    --istio.enabled=true || return $?

  BROKER_TEMPLATES="${KAFKA_BROKER_TEMPLATES}" go_test_e2e \
    -timeout=1h \
    -parallel=12 \
    -run TestContainerSource \
    ./test/rekt/ \
    --istio.enabled=true || return $?

  BROKER_TEMPLATES="${KAFKA_BROKER_TEMPLATES}" go_test_e2e \
    -timeout=1h \
    -parallel=12 \
    -run TestSinkBinding \
    ./test/rekt/ \
    --istio.enabled=true || return $?

  CHANNEL_GROUP_KIND="InMemoryChannel.messaging.knative.dev" \
  CHANNEL_VERSION="v1" \
  go_test_e2e \
    -timeout=1h \
    -parallel=18 \
    -run TestChannel \
    ./test/rekt/ \
    --istio.enabled=true || return $?

  CHANNEL_GROUP_KIND="KafkaChannel.messaging.knative.dev" \
  CHANNEL_VERSION="v1beta1" \
  go_test_e2e \
    -timeout=1h \
    -parallel=18 \
    -run TestChannel \
    ./test/rekt/ \
    --istio.enabled=true || return $?

  popd
}

function run_eventing_kafka_broker_tests() {
  pushd "${repo_root_dir}"/third_party/eventing-kafka-broker || return $?

  BROKER_TEMPLATES="${KAFKA_BROKER_TEMPLATES}" BROKER_CLASS="Kafka" go_test_e2e \
    -timeout=1h \
    -parallel=12 \
    -run TestKafkaSource \
    ./test/e2e_new/... \
    --istio.enabled=true || return $?

  BROKER_TEMPLATES="${KAFKA_BROKER_TEMPLATES}" BROKER_CLASS="Kafka" go_test_e2e \
    -timeout=1h \
    -parallel=12 \
    -run TestKafkaSink \
    ./test/e2e_new/... \
    --istio.enabled=true || return $?

  popd
}
