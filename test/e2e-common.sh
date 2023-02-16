#!/usr/bin/env bash

REPO_ROOT_DIR=$(dirname $0)/..

export SKIP_INITIALIZE=${SKIP_INITIALIZE:-false}

source $(pwd)/vendor/knative.dev/hack/e2e-tests.sh

git submodule update --init --recursive

function knative_setup() {
  "${REPO_ROOT_DIR}"/hack/update-istio.sh || return $?
  "${REPO_ROOT_DIR}"/hack/install-dependencies.sh || return $?
}
