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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serviceinformer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	istionetworking "istio.io/client-go/pkg/apis/networking/v1beta1"

	"knative.dev/eventing-istio/pkg/apis/config"
	"knative.dev/eventing-istio/pkg/client/injection/kube/reconciler/core/v1/service"
	istioclient "knative.dev/eventing-istio/pkg/client/istio/injection/client"
	istionetworkinginformer "knative.dev/eventing-istio/pkg/client/istio/injection/informers/networking/v1beta1/destinationrule/filtered"
)

func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {

	store := config.NewStore(ctx)
	store.WatchConfigs(cmw)

	// Get a filtered informer
	drInformer := istionetworkinginformer.Get(ctx, IstioResourceSelector)

	ic := istioclient.Get(ctx)

	r := &Reconciler{
		IstioClient:           ic,
		DestinationRuleLister: drInformer.Lister(),
		GetConfig: func(ctx context.Context, svc *corev1.Service) *config.Config {
			return config.Load(store)
		},
	}

	imcLabels := map[string]string{
		"messaging.knative.dev/channel": "in-memory-channel",
		"messaging.knative.dev/role":    "dispatcher",
	}
	kcLabels := map[string]string{
		"messaging.knative.dev/role": "kafka-channel",
	}

	filterServices := func(obj interface{}) bool {
		svc, ok := obj.(metav1.Object)
		if !ok {
			return false
		}
		l := labels.SelectorFromSet(svc.GetLabels())

		imcSelector := labels.Set(imcLabels)
		kcSelector := labels.Set(kcLabels)

		return l.Matches(imcSelector) || l.Matches(kcSelector)
	}

	impl := service.NewImpl(ctx, r, func(impl *controller.Impl) controller.Options {
		return controller.Options{
			SkipStatusUpdates: true,
			PromoteFilterFunc: filterServices,
		}
	})

	handleServices(ctx, metav1.LabelSelector{MatchLabels: imcLabels}, impl)
	handleServices(ctx, metav1.LabelSelector{MatchLabels: kcLabels}, impl)

	// Notify the tracker that a destination rule we're tracking changed.
	drInformer.Informer().AddEventHandler(controller.HandleAll(controller.EnsureTypeMeta(
		impl.Tracker.OnChanged,
		schema.GroupVersionKind{
			Group:   istionetworking.SchemeGroupVersion.Group,
			Version: istionetworking.SchemeGroupVersion.Version,
			Kind:    "DestinationRule",
		},
	)))

	return impl
}

func handleServices(ctx context.Context, labelSelector metav1.LabelSelector, impl *controller.Impl) {
	imcSelector, err := metav1.LabelSelectorAsSelector(&labelSelector)
	if err != nil {
		panic(err)
	}
	informer := serviceinformer.NewFilteredServiceInformer(
		kubeclient.Get(ctx),
		corev1.NamespaceAll,
		controller.GetResyncPeriod(ctx),
		nil,
		func(options *metav1.ListOptions) {
			options.LabelSelector = imcSelector.String()
		},
	)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    impl.Enqueue,
		UpdateFunc: controller.PassNew(impl.Enqueue),
	})
}
