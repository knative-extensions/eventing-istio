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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	gentype "k8s.io/client-go/gentype"
	networkingv1beta1 "knative.dev/eventing-istio/pkg/client/istio/clientset/versioned/typed/networking/v1beta1"
)

// fakeSidecars implements SidecarInterface
type fakeSidecars struct {
	*gentype.FakeClientWithList[*v1beta1.Sidecar, *v1beta1.SidecarList]
	Fake *FakeNetworkingV1beta1
}

func newFakeSidecars(fake *FakeNetworkingV1beta1, namespace string) networkingv1beta1.SidecarInterface {
	return &fakeSidecars{
		gentype.NewFakeClientWithList[*v1beta1.Sidecar, *v1beta1.SidecarList](
			fake.Fake,
			namespace,
			v1beta1.SchemeGroupVersion.WithResource("sidecars"),
			v1beta1.SchemeGroupVersion.WithKind("Sidecar"),
			func() *v1beta1.Sidecar { return &v1beta1.Sidecar{} },
			func() *v1beta1.SidecarList { return &v1beta1.SidecarList{} },
			func(dst, src *v1beta1.SidecarList) { dst.ListMeta = src.ListMeta },
			func(list *v1beta1.SidecarList) []*v1beta1.Sidecar { return gentype.ToPointerSlice(list.Items) },
			func(list *v1beta1.SidecarList, items []*v1beta1.Sidecar) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
