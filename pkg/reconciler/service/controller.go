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

	"go.uber.org/zap"
	istionetworking "istio.io/client-go/pkg/apis/networking/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	"knative.dev/eventing-istio/pkg/apis/config"
	serviceinformer "knative.dev/eventing-istio/pkg/client/injection/kube/informers/core/v1/service"
	"knative.dev/eventing-istio/pkg/client/injection/kube/reconciler/core/v1/service"
	istioclientset "knative.dev/eventing-istio/pkg/client/istio/injection/client"
	istionetworkinginformer "knative.dev/eventing-istio/pkg/client/istio/injection/informers/networking/v1beta1/destinationrule/filtered"
)

func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {

	logger := logging.FromContext(ctx).Desugar()

	var globalResync func()

	store := config.NewStore(ctx, func(name string, value interface{}) {
		if globalResync == nil {
			return
		}

		logger.Info("Config store changed",
			zap.String("name", name),
			zap.Any("value", value.(*config.Config)))

		globalResync()
	})
	store.WatchConfigs(cmw)

	// Get a filtered informer for destination rules
	drInformer := istionetworkinginformer.Get(ctx, IstioResourceSelector)

	ic := istioclientset.Get(ctx)
	serviceInformer := serviceinformer.Get(ctx).Informer()

	r := &Reconciler{
		IstioClient:           ic,
		DestinationRuleLister: drInformer.Lister(),
		GetConfig: func(ctx context.Context, svc *corev1.Service) *config.Config {
			return config.Load(store)
		},
	}

	impl := service.NewImpl(ctx, r, func(impl *controller.Impl) controller.Options {
		return controller.Options{
			SkipStatusUpdates: true,
			PromoteFilterFunc: filterServices(ctx),
		}
	})

	globalResync = func() {
		impl.FilteredGlobalResync(filterServices(ctx), serviceInformer)
	}

	r.Tracker = impl.Tracker

	serviceInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: filterServices(ctx),
		Handler: controller.HandleAll(controller.EnsureTypeMeta(
			impl.Enqueue,
			schema.GroupVersionKind{
				Group:   corev1.SchemeGroupVersion.Group,
				Version: corev1.SchemeGroupVersion.Version,
				Kind:    "Service",
			},
		)),
	})

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

func filterServices(ctx context.Context) func(obj interface{}) bool {
	logger := logging.FromContext(ctx).Desugar()

	return func(obj interface{}) bool {
		imcLabels := map[string]string{
			"messaging.knative.dev/role": "in-memory-channel",
		}
		kcLabels := map[string]string{
			"messaging.knative.dev/role": "kafka-channel",
		}

		svc, ok := obj.(*corev1.Service)
		if !ok {
			return false
		}

		logger.Debug("Filtering Service",
			zap.String("namespace", svc.GetNamespace()),
			zap.String("name", svc.GetName()),
			zap.Any("labels", svc.GetLabels()),
		)

		if svc.Spec.ExternalName == "" {
			return false
		}

		l := labels.Set(svc.GetLabels())

		imcSelector := labels.SelectorFromSet(imcLabels)
		kcSelector := labels.SelectorFromSet(kcLabels)

		return imcSelector.Matches(l) || kcSelector.Matches(l)
	}
}
