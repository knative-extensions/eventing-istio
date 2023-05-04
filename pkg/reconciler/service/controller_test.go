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

package service

import (
	"context"
	"testing"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"

	_ "knative.dev/pkg/client/injection/kube/client/fake"

	_ "knative.dev/eventing-istio/pkg/client/injection/kube/informers/core/v1/service/fake"
	_ "knative.dev/eventing-istio/pkg/client/istio/injection/client/fake"
	istiofilteredfactory "knative.dev/eventing-istio/pkg/client/istio/injection/informers/factory/filtered"
	_ "knative.dev/eventing-istio/pkg/client/istio/injection/informers/factory/filtered/fake"
	_ "knative.dev/eventing-istio/pkg/client/istio/injection/informers/networking/v1beta1/destinationrule/filtered/fake"
)

func TestFilterServices(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := logging.WithLogger(context.Background(), logger.Sugar())

	tt := []struct {
		name     string
		obj      interface{}
		expected bool
	}{
		{
			name: "IMC pass",
			obj: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"messaging.knative.dev/role":  "in-memory-channel",
						"messaging.knative.dev/role1": "dispatcher",
					},
				},
				Spec: corev1.ServiceSpec{
					ExternalName: "example.com",
				},
			},
			expected: true,
		},
		{
			name: "IMC not pass",
			obj: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"messaging.knative.dev/channel": "in-memory-channel2",
					},
				},
				Spec: corev1.ServiceSpec{
					ExternalName: "example.com",
				},
			},
			expected: false,
		},
		{
			name: "KafkaChannel pass",
			obj: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"messaging.knative.dev/role":  "kafka-channel",
						"messaging.knative.dev/role2": "kafka-channel",
					},
				},
				Spec: corev1.ServiceSpec{
					ExternalName: "example.com",
				},
			},
			expected: true,
		},
		{
			name: "KafkaChannel not pass",
			obj: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"messaging.knative.dev/role": "kafka-channel1",
					},
				},
				Spec: corev1.ServiceSpec{
					ExternalName: "example.com",
				},
			},
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			if _, ok := tc.obj.(runtime.Object); !ok {
				panic(tc.obj)
			}

			got := filterServices(ctx)(tc.obj)
			if tc.expected != got {
				t.Fatal("expected", tc.expected, "got", got)
			}
		})
	}
}

func TestNewController(t *testing.T) {
	ctx := istiofilteredfactory.WithSelectors(context.Background(), IstioResourceSelector)
	ctx, _ = injection.Fake.SetupInformers(ctx, &rest.Config{})

	impl := NewController(ctx, &configmap.ManualWatcher{})
	if impl == nil {
		t.Fatal("impl must not be nil")
	}
}
