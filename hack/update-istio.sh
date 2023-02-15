#!/usr/bin/env bash

export ISTIO_VERSION=1.16.1

curl -L https://istio.io/downloadIstio | sh -

mv "istio-${ISTIO_VERSION}" "third_party/istio"