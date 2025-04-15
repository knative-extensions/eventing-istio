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

// fakeWorkloadEntries implements WorkloadEntryInterface
type fakeWorkloadEntries struct {
	*gentype.FakeClientWithList[*v1beta1.WorkloadEntry, *v1beta1.WorkloadEntryList]
	Fake *FakeNetworkingV1beta1
}

func newFakeWorkloadEntries(fake *FakeNetworkingV1beta1, namespace string) networkingv1beta1.WorkloadEntryInterface {
	return &fakeWorkloadEntries{
		gentype.NewFakeClientWithList[*v1beta1.WorkloadEntry, *v1beta1.WorkloadEntryList](
			fake.Fake,
			namespace,
			v1beta1.SchemeGroupVersion.WithResource("workloadentries"),
			v1beta1.SchemeGroupVersion.WithKind("WorkloadEntry"),
			func() *v1beta1.WorkloadEntry { return &v1beta1.WorkloadEntry{} },
			func() *v1beta1.WorkloadEntryList { return &v1beta1.WorkloadEntryList{} },
			func(dst, src *v1beta1.WorkloadEntryList) { dst.ListMeta = src.ListMeta },
			func(list *v1beta1.WorkloadEntryList) []*v1beta1.WorkloadEntry {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1beta1.WorkloadEntryList, items []*v1beta1.WorkloadEntry) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
