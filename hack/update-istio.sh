#!/usr/bin/env bash

set -euo pipefail

function update_istio() {
  version="$1"
  target_dir="$2"

  echo "Updating Istio: Version ${version}, target directory ${target_dir}"

  export ISTIO_VERSION=${version}
  curl -L https://istio.io/downloadIstio | sh -
  rm -rf "${target_dir}"
  mv "istio-${ISTIO_VERSION}" "${target_dir}"
}

update_istio "$(cat $(dirname $0)/../third_party/istio_version.txt)" "third_party/istio"
