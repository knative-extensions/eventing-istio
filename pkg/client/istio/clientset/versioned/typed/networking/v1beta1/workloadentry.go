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

package v1beta1

import (
	"context"

	v1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
	scheme "knative.dev/eventing-istio/pkg/client/istio/clientset/versioned/scheme"
)

// WorkloadEntriesGetter has a method to return a WorkloadEntryInterface.
// A group's client should implement this interface.
type WorkloadEntriesGetter interface {
	WorkloadEntries(namespace string) WorkloadEntryInterface
}

// WorkloadEntryInterface has methods to work with WorkloadEntry resources.
type WorkloadEntryInterface interface {
	Create(ctx context.Context, workloadEntry *v1beta1.WorkloadEntry, opts v1.CreateOptions) (*v1beta1.WorkloadEntry, error)
	Update(ctx context.Context, workloadEntry *v1beta1.WorkloadEntry, opts v1.UpdateOptions) (*v1beta1.WorkloadEntry, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, workloadEntry *v1beta1.WorkloadEntry, opts v1.UpdateOptions) (*v1beta1.WorkloadEntry, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.WorkloadEntry, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.WorkloadEntryList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.WorkloadEntry, err error)
	WorkloadEntryExpansion
}

// workloadEntries implements WorkloadEntryInterface
type workloadEntries struct {
	*gentype.ClientWithList[*v1beta1.WorkloadEntry, *v1beta1.WorkloadEntryList]
}

// newWorkloadEntries returns a WorkloadEntries
func newWorkloadEntries(c *NetworkingV1beta1Client, namespace string) *workloadEntries {
	return &workloadEntries{
		gentype.NewClientWithList[*v1beta1.WorkloadEntry, *v1beta1.WorkloadEntryList](
			"workloadentries",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *v1beta1.WorkloadEntry { return &v1beta1.WorkloadEntry{} },
			func() *v1beta1.WorkloadEntryList { return &v1beta1.WorkloadEntryList{} }),
	}
}
