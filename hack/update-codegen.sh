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

set -o errexit
set -o nounset
set -o pipefail

# shellcheck disable=SC1091
source "$(dirname "$0")"/../vendor/knative.dev/hack/codegen-library.sh
source "${CODEGEN_PKG}/kube_codegen.sh"

# If we run with -mod=vendor here, then generate-groups.sh looks for vendor files in the wrong place.
export GOFLAGS=-mod=

echo "=== Update Codegen for $MODULE_NAME"

group "Knative Codegen"

# Knative Injection (for istio)
OUTPUT_PKG="knative.dev/eventing-istio/pkg/client/istio/injection" \
${KNATIVE_CODEGEN_PKG}/hack/generate-knative.sh "injection" \
  istio.io/client-go/pkg istio.io/client-go/pkg/apis \
  "networking:v1beta1" \
  --go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

group "Update deps post-codegen"

# Make sure our dependencies are up-to-date
"${REPO_ROOT_DIR}"/hack/update-deps.sh
