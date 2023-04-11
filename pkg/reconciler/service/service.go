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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"

	istionetworking "istio.io/client-go/pkg/apis/networking/v1beta1"

	"knative.dev/eventing-istio/pkg/apis/config"
	servicereconciler "knative.dev/eventing-istio/pkg/client/injection/kube/reconciler/core/v1/service"
	"knative.dev/eventing-istio/pkg/client/istio/clientset/versioned"
	networkingv1beta1 "knative.dev/eventing-istio/pkg/client/istio/listers/networking/v1beta1"
)

type Reconciler struct {
	Tracker               tracker.Interface
	GetConfig             func(ctx context.Context, svc *corev1.Service) *config.Config
	IstioClient           versioned.Interface
	DestinationRuleLister networkingv1beta1.DestinationRuleLister
}

var (
	_ servicereconciler.Interface = &Reconciler{}
)

func (r *Reconciler) ReconcileKind(ctx context.Context, svc *corev1.Service) reconciler.Event {
	cfg := r.GetConfig(ctx, svc)
	if !cfg.IsEnabled() {
		// If the flag was disabled after being enabled finalize resources
		return r.finalizeDestinationRule(ctx, svc)
	}

	if err := r.reconcileDestinationRule(ctx, svc); err != nil {
		return fmt.Errorf("failed to reconcile DestinationRule: %w", err)
	}

	return nil
}

func (r *Reconciler) reconcileDestinationRule(ctx context.Context, svc *corev1.Service) error {

	expected := DestinationRule(DestinationRuleConfig{
		Service: svc,
	})

	got, err := r.DestinationRuleLister.DestinationRules(svc.Namespace).Get(svc.Name)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to get DestinationRule %s/%s: %w", svc.Namespace, svc.Name, err)
	}
	if apierrors.IsNotFound(err) {
		return r.createDestinationRule(ctx, expected)
	}

	_ = r.Tracker.TrackReference(tracker.Reference{
		APIVersion: got.APIVersion,
		Kind:       got.Kind,
		Namespace:  got.Namespace,
		Name:       got.Name,
	}, svc)

	if equality.Semantic.DeepDerivative(expected, got) {
		return nil
	}

	updated := &istionetworking.DestinationRule{
		TypeMeta:   expected.TypeMeta,
		ObjectMeta: *got.ObjectMeta.DeepCopy(),
		Spec:       *expected.Spec.DeepCopy(),
		Status:     *got.Status.DeepCopy(),
	}
	updated.Labels = expected.Labels
	updated.Annotations = expected.Annotations

	return r.updateDestinationRule(ctx, svc, updated)
}

func (r *Reconciler) createDestinationRule(ctx context.Context, expected *istionetworking.DestinationRule) error {
	_, err := r.IstioClient.NetworkingV1beta1().
		DestinationRules(expected.GetNamespace()).
		Create(ctx, expected, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create DestinationRule %s/%s: %w", expected.Namespace, expected.Name, err)
	}

	controller.GetEventRecorder(ctx).
		Event(expected, corev1.EventTypeNormal, "Created", fmt.Sprintf("Created DestinationRule %s/%s", expected.Namespace, expected.Name))

	return nil
}

func (r *Reconciler) updateDestinationRule(ctx context.Context, svc *corev1.Service, expected *istionetworking.DestinationRule) error {
	if !metav1.IsControlledBy(expected, svc) {
		return fmt.Errorf("owner: %s with Type %T does not own DestinationRule: %s/%s", svc.Name, svc, expected.Namespace, expected.Name)
	}

	_, err := r.IstioClient.NetworkingV1beta1().
		DestinationRules(expected.GetNamespace()).
		Update(ctx, expected, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update DestinationRule %s/%s: %w", expected.Namespace, expected.Name, err)
	}

	controller.GetEventRecorder(ctx).
		Event(expected, corev1.EventTypeNormal, "Updated", fmt.Sprintf("Updated DestinationRule %s/%s", expected.Namespace, expected.Name))

	return nil
}

func (r *Reconciler) finalizeDestinationRule(ctx context.Context, svc *corev1.Service) error {
	dr, err := r.DestinationRuleLister.DestinationRules(svc.Namespace).Get(svc.Name)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to get ")
	}
	if apierrors.IsNotFound(err) {
		return nil
	}

	if !metav1.IsControlledBy(dr, svc) {
		return nil
	}

	err = r.IstioClient.
		NetworkingV1beta1().
		DestinationRules(svc.Namespace).
		Delete(ctx, svc.Name, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete DestinationRule %s/%s: %w", svc.Namespace, svc.Name, err)
	}

	return nil
}