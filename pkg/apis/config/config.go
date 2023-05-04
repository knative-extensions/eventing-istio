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
	// IstioConfig is the name of config map for the Istio config features.
	IstioConfig = "config-features"

	// IstioConfigKey is the key in IstioConfig ConfigMap for the feature flag enabled or disable.
	IstioConfigKey = "istio"
)

type Config struct {
	Istio feature.Flag `json:"istio"`
}

type Option func(config *Config) error

func NewDefaultConfig(options ...Option) (*Config, error) {
	c := &Config{
		Istio: feature.Disabled,
	}

	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithEnabled() Option {
	return func(config *Config) error {
		config.Istio = feature.Enabled
		return nil
	}
}

func (c Config) IsEnabled() bool {
	return strings.EqualFold(string(c.Istio), string(feature.Enabled))
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
	c, err := NewDefaultConfig()
	if err != nil {
		return nil, err
	}
	if err := configmap.Parse(data, asFlag(IstioConfigKey, &c.Istio)); err != nil {
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
