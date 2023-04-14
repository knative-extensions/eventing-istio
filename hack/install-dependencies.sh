#!/usr/bin/env bash

set -euo pipefail

$(dirname $0)/update-istio.sh

export REPO_ROOT_DIR=${REPO_ROOT_DIR:-$(git rev-parse --show-toplevel)}

export EVENTING_CONFIG=${EVENTING_CONFIG:-"./third_party/eventing-latest/"}
export EVENTING_KAFKA_CONFIG=${EVENTING_KAFKA_CONFIG:-"./third_party/eventing-kafka-broker-latest/"}
export ISTIO_CONFIG_DIR=${ISTIO_CONFIG_DIR:-"./third_party/istio"}
export SYSTEM_NAMESPACE=${SYSTEM_NAMESPACE:-"knative-eventing"}
export ISTIO_NAMESPACE=${SYSTEM_NAMESPACE:-"istio-system"}

export PATH="third_party/istio/bin:$PATH"

echo "Installing Istio"
istioctl version

istioctl x precheck

istioctl install \
  -y \
  --set profile=default \
  --set meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY

kubectl apply -f "${ISTIO_CONFIG_DIR}/samples/addons"

kubectl create namespace knative-eventing --dry-run=client -oyaml | kubectl apply -f -
kubectl label namespace knative-eventing istio-injection=enabled

echo "Installing Eventing from ${EVENTING_CONFIG}"
kubectl apply -Rf "${EVENTING_CONFIG}"

kubectl patch deployment \
  eventing-webhook \
  --type merge \
  -n knative-eventing \
  --patch-file "${REPO_ROOT_DIR}/hack/eventing-injection-disabled.yaml"

kubectl patch deployment \
  imc-controller \
  --type merge \
  -n knative-eventing \
  --patch-file "${REPO_ROOT_DIR}/hack/eventing-injection-disabled.yaml"

echo "Installing Eventing Kafka from ${EVENTING_KAFKA_CONFIG}"
kubectl apply -Rf "${EVENTING_KAFKA_CONFIG}"

kubectl patch deployment \
  kafka-webhook-eventing \
  --type merge \
  -n knative-eventing \
  --patch-file "${REPO_ROOT_DIR}/hack/eventing-injection-disabled.yaml"

"${REPO_ROOT_DIR}"/third_party/eventing-kafka-broker/test/kafka/kafka_setup.sh

kubectl apply -n "${ISTIO_NAMESPACE}" -Rf "${REPO_ROOT_DIR}"/test/config
kubectl apply -n "${SYSTEM_NAMESPACE}" -Rf "${REPO_ROOT_DIR}"/test/config

source "${REPO_ROOT_DIR}/third_party/eventing-kafka-broker/test/e2e-common.sh"

create_sasl_secrets || exit 1
create_tls_secrets || exit 1
