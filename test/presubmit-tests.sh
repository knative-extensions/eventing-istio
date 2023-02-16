#!/usr/bin/env bash

source $(dirname $0)/e2e-common.sh
source $(dirname $0)/../vendor/knative.dev/hack/presubmit-tests.sh

main $@
