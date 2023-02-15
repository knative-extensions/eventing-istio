#!/usr/bin/env bash

set -euo pipefail

export REPO_ROOT_DIR=${REPO_ROOT_DIR:-$(git rev-parse --show-toplevel)}

export EVENTING_CONFIG=${EVENTING_CONFIG:-"./third_party/eventing-latest/"}
export EVENTING_KAFKA_CONFIG=${EVENTING_KAFKA_CONFIG:-"./third_party/eventing-kafka-broker-latest/"}
export ISTIO_CONFIG_DIR=${ISTIO_CONFIG_DIR:-"./third_party/istio"}
export EVENTING_ISTIO_RESOURCES_CONFIG=${EVENTING_ISTIO_RESOURCES_CONFIG:-"./third_party/istio-resources/"}

export PATH="third_party/istio/bin:$PATH"

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

echo "Installing Istio"
istioctl version

istioctl install \
  -y \
  --set profile=default \
  --set meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY

kubectl apply -f "${ISTIO_CONFIG_DIR}/samples/addons"

kubectl apply -Rf "${EVENTING_ISTIO_RESOURCES_CONFIG}"

echo "Installing Eventing Kafka from ${EVENTING_KAFKA_CONFIG}"
kubectl apply -Rf "${EVENTING_KAFKA_CONFIG}"

kubectl patch deployment \
  kafka-webhook-eventing \
  --type merge \
  -n knative-eventing \
  --patch-file "${REPO_ROOT_DIR}/hack/eventing-injection-disabled.yaml"

"${REPO_ROOT_DIR}"/third_party/kafka/kafka_setup.sh

# curl -X POST -v -H "content-type: application/json" -H "ce-specversion: 1.0" -H "ce-source: my/curl/command" -H "ce-type: my.demo.event" -H "ce-id: 0815" -d '{"value":"Hello Knative"}' http://broker-ingress.knative-eventing.svc.cluster.local/test/default

