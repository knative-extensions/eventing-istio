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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// DestinationRuleLister helps list DestinationRules.
// All objects returned here must be treated as read-only.
type DestinationRuleLister interface {
	// List lists all DestinationRules in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.DestinationRule, err error)
	// DestinationRules returns an object that can list and get DestinationRules.
	DestinationRules(namespace string) DestinationRuleNamespaceLister
	DestinationRuleListerExpansion
}

// destinationRuleLister implements the DestinationRuleLister interface.
type destinationRuleLister struct {
	listers.ResourceIndexer[*v1beta1.DestinationRule]
}

// NewDestinationRuleLister returns a new DestinationRuleLister.
func NewDestinationRuleLister(indexer cache.Indexer) DestinationRuleLister {
	return &destinationRuleLister{listers.New[*v1beta1.DestinationRule](indexer, v1beta1.Resource("destinationrule"))}
}

// DestinationRules returns an object that can list and get DestinationRules.
func (s *destinationRuleLister) DestinationRules(namespace string) DestinationRuleNamespaceLister {
	return destinationRuleNamespaceLister{listers.NewNamespaced[*v1beta1.DestinationRule](s.ResourceIndexer, namespace)}
}

// DestinationRuleNamespaceLister helps list and get DestinationRules.
// All objects returned here must be treated as read-only.
type DestinationRuleNamespaceLister interface {
	// List lists all DestinationRules in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.DestinationRule, err error)
	// Get retrieves the DestinationRule from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.DestinationRule, error)
	DestinationRuleNamespaceListerExpansion
}

// destinationRuleNamespaceLister implements the DestinationRuleNamespaceLister
// interface.
type destinationRuleNamespaceLister struct {
	listers.ResourceIndexer[*v1beta1.DestinationRule]
}
