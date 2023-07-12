#!/usr/bin/env bash

# Copyright 2023 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/../vendor/knative.dev/hack/release.sh

readonly EVENTING_ISTIO_ARTIFACT="eventing-istio.yaml"
readonly EVENTING_ISTIO_CONFIG_DIR=config

function fail() {
  echo "$1"
  exit 1
}

function generate_artifacts() {
  "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"/build-from-source.sh || return $?
}

function update_release_labels() {
  TAG=${TAG:-$(git rev-parse HEAD)}
  
  echo "Updating release labels to app.kubernetes.io/version: \"${TAG}\""

  sed -e "s|app.kubernetes.io/version: devel|app.kubernetes.io/version: \"${TAG}\"|" -i ${EVENTING_ISTIO_ARTIFACT}
}

function build_release() {
  [ -f "${EVENTING_ISTIO_ARTIFACT}" ] && rm "${EVENTING_ISTIO_ARTIFACT}"

  generate_artifacts
  if [[ $? -ne 0 ]]; then
    fail "failed to generate artifacts"
  fi

  update_release_labels
  if [[ $? -ne 0 ]]; then
    fail "failed to update release labels on artifacts"
  fi

  export ARTIFACTS_TO_PUBLISH=(
    "${EVENTING_ISTIO_ARTIFACT}"
  )

  # ARTIFACTS_TO_PUBLISH has to be a string, not an array.
  # shellcheck disable=SC2178
  # shellcheck disable=SC2124
  export ARTIFACTS_TO_PUBLISH="${ARTIFACTS_TO_PUBLISH[@]}"
}

main $@
