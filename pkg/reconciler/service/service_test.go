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
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientgotesting "k8s.io/client-go/testing"
	"k8s.io/utils/pointer"
	"knative.dev/pkg/network"

	istionetworkingapi "istio.io/api/networking/v1beta1"
	istionetworking "istio.io/client-go/pkg/apis/networking/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	logtesting "knative.dev/pkg/logging/testing"
	. "knative.dev/pkg/reconciler/testing"
	"knative.dev/pkg/tracker"

	kubeclient "knative.dev/eventing-istio/pkg/client/injection/kube/client/fake"

	"knative.dev/eventing-istio/pkg/apis/config"
	"knative.dev/eventing-istio/pkg/client/injection/kube/reconciler/core/v1/service"
	istioclientset "knative.dev/eventing-istio/pkg/client/istio/injection/client/fake"
	. "knative.dev/eventing-istio/pkg/reconciler/testing"
)

const (
	configKey = "istio-testing-config"

	serviceName         = "service-name"
	serviceUUID         = "abc"
	serviceNamespace    = "service-namespace"
	serviceExternalName = "channel.namespace.svc.cluster.local"

	destinationRuleUUID = "xyz"
)

var (
	defaultCmpOpts = []cmp.Option{protocmp.Transform()}
)

func TestReconcileKind(t *testing.T) {

	logger := logtesting.TestLogger(t)

	key := fmt.Sprintf("%s/%s", serviceNamespace, serviceName)

	r := TableTest{
		{
			Name: "istio disabled by default",
			Objects: []runtime.Object{
				makeService(),
			},
			Key:     key,
			WantErr: false,
		},
		{
			Name: "istio enabled, create DestinationRule",
			Objects: []runtime.Object{
				makeService(),
			},
			Key:     key,
			WantErr: false,
			OtherTestData: map[string]interface{}{
				configKey: panicOnError(config.NewDefaultConfig(config.WithEnabled())),
			},
			WantCreates: []runtime.Object{
				makeDestinationRuleForService(),
			},
			WantEvents: []string{
				makeDestinationRuleForServiceCreatedEvent(),
			},
			CmpOpts: defaultCmpOpts,
		},
		{
			Name: "istio enabled, no DestinationRule update",
			Objects: []runtime.Object{
				makeService(),
				makeDestinationRuleForService(),
			},
			Key:     key,
			WantErr: false,
			OtherTestData: map[string]interface{}{
				configKey: panicOnError(config.NewDefaultConfig(config.WithEnabled())),
			},
			CmpOpts: defaultCmpOpts,
		},
		{
			Name: "istio enabled, DestinationRule updated",
			Objects: []runtime.Object{
				makeService(),
				makeDestinationRuleForService(func(rule *istionetworking.DestinationRule) {
					rule.Spec.Host = serviceExternalName
				}),
			},
			Key:     key,
			WantErr: false,
			OtherTestData: map[string]interface{}{
				configKey: panicOnError(config.NewDefaultConfig(config.WithEnabled())),
			},
			WantUpdates: []clientgotesting.UpdateActionImpl{
				{Object: makeDestinationRuleForService()},
			},
			WantEvents: []string{
				makeDestinationRuleForServiceUpdatedEvent(),
			},
			CmpOpts: defaultCmpOpts,
		},
		{
			Name: "istio disabled, DestinationRule deleted",
			Objects: []runtime.Object{
				makeService(),
				makeDestinationRuleForService(func(rule *istionetworking.DestinationRule) {
					rule.Spec.Host = network.GetServiceHostname(serviceName, serviceNamespace)
				}),
			},
			Key:     key,
			WantErr: false,
			WantDeletes: []clientgotesting.DeleteActionImpl{
				deleteDestinationRuleForService(),
			},
			WantEvents: []string{
				makeDestinationRuleForServiceDeletedEvent(),
			},
			CmpOpts: defaultCmpOpts,
		},
		{
			Name: "istio disabled, DestinationRule not controlled by us",
			Objects: []runtime.Object{
				makeService(),
				makeDestinationRuleForService(func(rule *istionetworking.DestinationRule) {
					rule.ObjectMeta.OwnerReferences = nil
				}),
			},
			Key:     key,
			WantErr: false,
			CmpOpts: defaultCmpOpts,
		},
	}

	r.Test(t, MakeFactory(func(ctx context.Context, listers *Listers, row *TableRow, watcher configmap.Watcher) controller.Reconciler {

		r := &Reconciler{
			GetConfig: func(ctx context.Context, svc *corev1.Service) *config.Config {
				v, ok := row.OtherTestData[configKey]
				if ok {
					return v.(*config.Config)
				}
				c, err := config.NewDefaultConfig()
				if err != nil {
					panic(err)
				}
				return c
			},
			IstioClient:           istioclientset.Get(ctx),
			DestinationRuleLister: listers.GetDestinationRuleLister(),
			Tracker:               tracker.New(func(types.NamespacedName) {}, time.Second),
		}

		return service.NewReconciler(ctx,
			logging.FromContext(ctx),
			kubeclient.Get(ctx),
			listers.GetServiceLister(),
			controller.GetEventRecorder(ctx),
			r,
		)

	}, logger))
}

func deleteDestinationRuleForService() clientgotesting.DeleteActionImpl {
	uid := types.UID(destinationRuleUUID)
	return clientgotesting.DeleteActionImpl{
		ActionImpl: clientgotesting.ActionImpl{
			Namespace: serviceNamespace,
			Verb:      "",
			Resource: schema.GroupVersionResource{
				Group:    istionetworking.SchemeGroupVersion.Group,
				Version:  istionetworking.SchemeGroupVersion.Version,
				Resource: "destinationrules",
			},
		},
		Name: serviceName,
		DeleteOptions: metav1.DeleteOptions{
			Preconditions: &metav1.Preconditions{
				UID: &uid,
			},
		},
	}
}

func makeDestinationRuleForService(opts ...func(rule *istionetworking.DestinationRule)) runtime.Object {
	s := makeService()
	s.Labels["app.kubernetes.io/managed-by"] = "knative-eventing-istio-controller"

	dr := &istionetworking.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: serviceNamespace,
			Name:      serviceName,
			OwnerReferences: []metav1.OwnerReference{
				makeServiceOwnerReference(),
			},
			Labels: s.Labels,
		},
		Spec: istionetworkingapi.DestinationRule{
			Host: network.GetServiceHostname(serviceName, serviceNamespace),
			TrafficPolicy: &istionetworkingapi.TrafficPolicy{
				Tls: &istionetworkingapi.ClientTLSSettings{
					Mode: istionetworkingapi.ClientTLSSettings_ISTIO_MUTUAL,
				},
			},
		},
	}

	for _, opt := range opts {
		opt(dr)
	}

	return dr
}

func makeService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: serviceNamespace,
			Name:      serviceName,
			UID:       serviceUUID,
			Labels: map[string]string{
				"messaging.knative.dev/role": "in-memory-channel",
			},
		},
		Spec: corev1.ServiceSpec{
			ExternalName: serviceExternalName,
		},
	}
}

func makeServiceOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion:         corev1.SchemeGroupVersion.String(),
		Kind:               "Service",
		Name:               serviceName,
		UID:                serviceUUID,
		Controller:         pointer.Bool(true),
		BlockOwnerDeletion: pointer.Bool(true),
	}
}

func makeDestinationRuleForServiceCreatedEvent() string {
	return Eventf(corev1.EventTypeNormal, "Created", fmt.Sprintf("Created DestinationRule %s/%s", serviceNamespace, serviceName))
}

func makeDestinationRuleForServiceDeletedEvent() string {
	return Eventf(corev1.EventTypeNormal, "Deleted", fmt.Sprintf("Deleted DestinationRule %s/%s", serviceNamespace, serviceName))
}

func makeDestinationRuleForServiceUpdatedEvent() string {
	return Eventf(corev1.EventTypeNormal, "Updated", fmt.Sprintf("Updated DestinationRule %s/%s", serviceNamespace, serviceName))
}

func panicOnError(defaultConfig interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return defaultConfig
}
