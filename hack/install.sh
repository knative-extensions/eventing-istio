#!/usr/bin/env bash

export REPO_ROOT_DIR=${REPO_ROOT_DIR:-$(git rev-parse --show-toplevel)}
export SYSTEM_NAMESPACE=${SYSTEM_NAMESPACE:-"knative-eventing"}
export ISTIO_NAMESPACE=${ISTIO_NAMESPACE:-"istio-system"}

function eventing_istio_setup() {
  "${REPO_ROOT_DIR}"/hack/build-from-source.sh || return $?

  kubectl apply -f "${REPO_ROOT_DIR}"/eventing-istio.yaml || return $?

  kubectl apply -f "${REPO_ROOT_DIR}"/test/config -n "${SYSTEM_NAMESPACE}"
}

eventing_istio_setup || exit 1
