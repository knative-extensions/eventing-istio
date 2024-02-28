/*
Copyright 2024 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apisconfig "knative.dev/eventing-istio/pkg/apis/config"
	"knative.dev/eventing/pkg/apis/feature"
)

func VerifyOIDCAndIstioNotEnabledSameTime(config *corev1.ConfigMap) (feature.Flags, error) {
	flags, err := feature.NewFlagsConfigFromConfigMap(config)
	if err != nil {
		return nil, fmt.Errorf("could not parse features configmap: %w", err)
	}

	if flags.IsOIDCAuthentication() && flags.IsEnabled(apisconfig.IstioConfigKey) {
		return nil, fmt.Errorf("%q feature can't be enabled while Istio is enabled too", feature.OIDCAuthentication)
	}

	return flags, nil
}
