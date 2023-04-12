/*
Copyright 2023 The Knative Authors

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

package config

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/eventing/pkg/apis/feature"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/logging"
)

const (
	// IstioConfig is the name of config map for the istio config features.
	IstioConfig = "config-features"

	IstioConfigKey = "istio"
)

type Config struct {
	istio                  feature.Flag
	destinationRuleTLSMode string
}

func (c Config) IsEnabled() bool {
	return strings.EqualFold(string(c.istio), string(feature.Enabled))
}

// Store is a typed wrapper around configmap.Untyped store to handle our configmaps.
// +k8s:deepcopy-gen=false
type Store struct {
	*configmap.UntypedStore
}

func NewStore(ctx context.Context, onAfterStore ...func(name string, value interface{})) *Store {
	return &Store{
		UntypedStore: configmap.NewUntypedStore(
			IstioConfig,
			logging.FromContext(ctx),
			configmap.Constructors{
				IstioConfig: newIstioConfig,
			},
			onAfterStore...,
		),
	}
}

func Load(store *Store) *Config {
	return store.UntypedStore.UntypedLoad(IstioConfig).(*Config)
}

func newIstioConfig(config *corev1.ConfigMap) (*Config, error) {
	return newIstioConfigFromMap(config.Data)
}

func newIstioConfigFromMap(data map[string]string) (*Config, error) {
	c := &Config{}
	if err := configmap.Parse(data, asFlag(IstioConfigKey, &c.istio)); err != nil {
		return c, fmt.Errorf("failed to parse flag %s: %w", IstioConfigKey, err)
	}
	return c, nil
}

func asFlag(key string, f *feature.Flag) configmap.ParseFunc {
	return func(m map[string]string) error {
		v, ok := m[key]
		if !ok {
			return nil
		}
		*f = feature.Flag(v)
		return nil
	}
}
