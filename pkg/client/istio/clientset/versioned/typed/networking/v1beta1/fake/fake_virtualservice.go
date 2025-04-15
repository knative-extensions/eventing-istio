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

// fakeVirtualServices implements VirtualServiceInterface
type fakeVirtualServices struct {
	*gentype.FakeClientWithList[*v1beta1.VirtualService, *v1beta1.VirtualServiceList]
	Fake *FakeNetworkingV1beta1
}

func newFakeVirtualServices(fake *FakeNetworkingV1beta1, namespace string) networkingv1beta1.VirtualServiceInterface {
	return &fakeVirtualServices{
		gentype.NewFakeClientWithList[*v1beta1.VirtualService, *v1beta1.VirtualServiceList](
			fake.Fake,
			namespace,
			v1beta1.SchemeGroupVersion.WithResource("virtualservices"),
			v1beta1.SchemeGroupVersion.WithKind("VirtualService"),
			func() *v1beta1.VirtualService { return &v1beta1.VirtualService{} },
			func() *v1beta1.VirtualServiceList { return &v1beta1.VirtualServiceList{} },
			func(dst, src *v1beta1.VirtualServiceList) { dst.ListMeta = src.ListMeta },
			func(list *v1beta1.VirtualServiceList) []*v1beta1.VirtualService {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1beta1.VirtualServiceList, items []*v1beta1.VirtualService) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
