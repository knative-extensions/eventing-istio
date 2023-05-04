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
	istionetworkingapi "istio.io/api/networking/v1beta1"
	istionetworking "istio.io/client-go/pkg/apis/networking/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/network"
)

const (
	istioResourceLabelKey = "app.kubernetes.io/managed-by"
	istioResource         = "knative-eventing-istio-controller"

	IstioResourceSelector = istioResourceLabelKey + "=" + istioResource
)

type DestinationRuleConfig struct {
	Service *corev1.Service
}

func DestinationRule(cfg DestinationRuleConfig) *istionetworking.DestinationRule {
	dr := &istionetworking.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   cfg.Service.Namespace,
			Name:        cfg.Service.Name,
			Labels:      cfg.Service.Labels,
			Annotations: cfg.Service.Annotations,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cfg.Service, schema.GroupVersionKind{
					Group:   corev1.SchemeGroupVersion.Group,
					Version: corev1.SchemeGroupVersion.Version,
					Kind:    "Service",
				}),
			},
		},
		Spec: istionetworkingapi.DestinationRule{
			Host: network.GetServiceHostname(cfg.Service.Name, cfg.Service.Namespace),
			TrafficPolicy: &istionetworkingapi.TrafficPolicy{
				Tls: &istionetworkingapi.ClientTLSSettings{
					Mode: istionetworkingapi.ClientTLSSettings_ISTIO_MUTUAL,
				},
			},
		},
	}

	if dr.Labels == nil {
		dr.Labels = make(map[string]string, 1)
	}
	dr.Labels[istioResourceLabelKey] = istioResource

	return dr
}
