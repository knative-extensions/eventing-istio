#!/usr/bin/env bash

export REPO_ROOT_DIR=${REPO_ROOT_DIR:-$(git rev-parse --show-toplevel)}

ko resolve ${KO_FLAGS} -Rf "${REPO_ROOT_DIR}"/config/eventing-istio >"${REPO_ROOT_DIR}"/eventing-istio.yaml
