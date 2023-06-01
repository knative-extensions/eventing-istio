#!/usr/bin/env bash

export REPO_ROOT_DIR=${REPO_ROOT_DIR:-$(git rev-parse --show-toplevel)}

if [[ -n "${TAG:-}" ]]; then
  LABEL_YAML_CMD=(sed -e "s|app.kubernetes.io/version: devel|app.kubernetes.io/version: \"${TAG:1}\"|")
else
  LABEL_YAML_CMD=(cat)
fi

ko resolve ${KO_FLAGS} -Rf "${REPO_ROOT_DIR}"/config/eventing-istio | "${LABEL_YAML_CMD[@]}" >"${REPO_ROOT_DIR}"/eventing-istio.yaml
