#!/usr/bin/env bash

export REPO_ROOT_DIR=${REPO_ROOT_DIR:-$(git rev-parse --show-toplevel)}

function eventing_istio_setup() {
  "${REPO_ROOT_DIR}"/hack/build-from-source.sh || return $?

  kubectl apply -f "${REPO_ROOT_DIR}"/eventing-istio.yaml || return $?
}

eventing_istio_setup || exit 1
