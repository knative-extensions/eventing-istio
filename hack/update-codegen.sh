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

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.

group "Kubernetes Codegen"

# Based on: https://github.com/kubernetes/kubernetes/blob/8ddabd0da5cc54761f3216c08e99fa1a9f7ee2c5/hack/lib/init.sh#L116
# The '-path' is a hack to workaround the lack of portable `-depth 2`.
K8S_TYPES=$(find ./vendor/k8s.io/api -type d -path '*/*/*/*/*/*' | cut -d'/' -f 5-6 | sort | sed 's@/@:@g' |
  grep -v "admission:" | grep -v "imagepolicy:" | grep -v "abac:" | grep -v "componentconfig:")

# Generate our own client for istio (otherwise injection won't work)
kube::codegen::gen_client \
  --boilerplate "${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt" \
  --output-dir "${REPO_ROOT_DIR}/pkg/client/istio" \
  --output-pkg "knative.dev/eventing-istio/pkg/client/istio" \
  --with-watch \
  "${REPO_ROOT_DIR}/vendor/istio.io/client-go/pkg/apis"

group "Knative Codegen"

# Knative Injection (for istio)
"${KNATIVE_CODEGEN_PKG}"/hack/generate-knative.sh "injection" \
  knative.dev/eventing-istio/pkg/client/istio istio.io/client-go/pkg/apis \
  "networking:v1beta1" \
  --lister-has-pointer-elem=true \
  --go-header-file "${REPO_ROOT_DIR}"/hack/boilerplate/boilerplate.go.txt

group "Update deps post-codegen"

# Make sure our dependencies are up-to-date
"${REPO_ROOT_DIR}"/hack/update-deps.sh
