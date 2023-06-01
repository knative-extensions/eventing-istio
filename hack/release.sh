#!/usr/bin/env bash

# Copyright 2023 The Knative Authors
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

# Documentation about this script and how to use it can be found
# at https://github.com/knative/test-infra/tree/main/ci

source $(dirname $0)/../vendor/knative.dev/hack/release.sh

export EVENTING_ISTIO_ARTIFACT="eventing-istio.yaml"

function build_release() {
  export TAG
  $(dirname $0)/build-from-source.sh || return $?

  export ARTIFACTS_TO_PUBLISH=(
    "${EVENTING_ISTIO_ARTIFACT}"
  )
}

main $@
