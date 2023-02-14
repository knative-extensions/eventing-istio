#!/usr/bin/env bash

source $(pwd)/vendor/knative.dev/hack/e2e-tests.sh

git submodule update --init --recursive
