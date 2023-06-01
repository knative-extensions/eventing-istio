#!/usr/bin/env bash

# Copyright 2020 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

source $(dirname "$0")/../vendor/knative.dev/hack/library.sh

version=$(echo $@ | grep -o "\-\-release \S*" | awk '{print $2}' || echo "")
upgrade=$(echo $@ | grep '\-\-upgrade' || echo "")
upgrade_artifacts=${UPGRADE_ARTIFACTS:-""}

function fetch_submodule() {
  branch=${1}
  module=${2}

  echo "Fetching branch ${branch} for submodule ${module}"
  pushd "${module}"
  git fetch origin "${branch}:${branch}" || return $?
  popd

  git submodule set-branch -b "${branch}" "${module}" || return $?
  git submodule sync --recursive || return $?
  git submodule update --init --recursive --remote || return $?
}

function update_submodule() {
  module=${1}
  branch=""

  if [ "${version}" = "" ] || [ "${version}" = "v9000.1" ]; then
    branch="main"
  else
    major_minor=${version:1} # Remove 'v' prefix
    branch="release-${major_minor}"
  fi

  if [ "${upgrade}" != "" ]; then
    fetch_submodule "${branch}" "${module}" || return $?
  else
    git submodule sync --recursive || return $?
    git submodule update --init --recursive --remote || return $?
  fi

}

function fetch_artifacts() {
  url="https://storage.googleapis.com/knative-nightly/${1}"
  echo "Fetch $url to ${2}"
  curl "${url}" > "${2}"
}

function update_submodules() {
  update_submodule "third_party/eventing" || return $?
  update_submodule "third_party/eventing-kafka-broker" || return $?
}

update_submodules || exit $?

if [ "${upgrade_artifacts}" != "" ]; then
  # Eventing
  e_dir="third_party/eventing-latest"
  rm -rf "${e_dir}" && mkdir "${e_dir}"
  e="eventing"
  fetch_artifacts "${e}/latest/eventing-core.yaml" "${e_dir}/eventing-core.yaml"
  fetch_artifacts "${e}/latest/eventing-crds.yaml" "${e_dir}/eventing-crds.yaml"
  fetch_artifacts "${e}/latest/in-memory-channel.yaml" "${e_dir}/in-memory-channel.yaml"
  fetch_artifacts "${e}/latest/mt-channel-broker.yaml" "${e_dir}/mt-channel-broker.yaml"
  # Eventing Kafka Broker
  ekb_dir="third_party/eventing-kafka-broker-latest"
  ekb="eventing-kafka-broker"
  rm -rf "${ekb_dir}" && mkdir "${ekb_dir}"
  fetch_artifacts "${ekb}/latest/eventing-kafka-controller.yaml" "${ekb_dir}/eventing-kafka-controller.yaml"
  fetch_artifacts "${ekb}/latest/eventing-kafka-broker.yaml" "${ekb_dir}/eventing-kafka-broker.yaml"
  fetch_artifacts "${ekb}/latest/eventing-kafka-channel.yaml" "${ekb_dir}/eventing-kafka-channel.yaml"
  fetch_artifacts "${ekb}/latest/eventing-kafka-sink.yaml" "${ekb_dir}/eventing-kafka-sink.yaml"
  fetch_artifacts "${ekb}/latest/eventing-kafka-source.yaml" "${ekb_dir}/eventing-kafka-source.yaml"
fi

$(dirname $0)/update-istio.sh

# Remove Istio binaries to avoid comparing them with verify-codegen.sh
rm -rf $(dirname $0)/../third_party/istio/bin

go_update_deps "$@"
