#!/usr/bin/env bash

set -euo pipefail

function update_istio() {
  version="$1"
  target_dir="$2"

  export ISTIO_VERSION=${version}
  curl -L https://istio.io/downloadIstio | sh -
  rm -rf "${target_dir}"
  mv "istio-${ISTIO_VERSION}" "${target_dir}"
}

update_istio "1.16.1" "third_party/istio"
