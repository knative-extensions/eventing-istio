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
	"knative.dev/pkg/logging"
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
			},
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := filterServices(ctx)(tc.obj)
			if tc.expected != got {
				t.Fatal("expected", tc.expected, "got", got)
			}
		})
	}
}
