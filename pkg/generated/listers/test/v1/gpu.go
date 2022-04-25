/*
Copyright The Kubernetes Authors.

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

package v1

import (
	v1 "crdsample/pkg/apis/test/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GpuLister helps list Gpus.
// All objects returned here must be treated as read-only.
type GpuLister interface {
	// List lists all Gpus in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Gpu, err error)
	// Gpus returns an object that can list and get Gpus.
	Gpus(namespace string) GpuNamespaceLister
	GpuListerExpansion
}

// gpuLister implements the GpuLister interface.
type gpuLister struct {
	indexer cache.Indexer
}

// NewGpuLister returns a new GpuLister.
func NewGpuLister(indexer cache.Indexer) GpuLister {
	return &gpuLister{indexer: indexer}
}

// List lists all Gpus in the indexer.
func (s *gpuLister) List(selector labels.Selector) (ret []*v1.Gpu, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Gpu))
	})
	return ret, err
}

// Gpus returns an object that can list and get Gpus.
func (s *gpuLister) Gpus(namespace string) GpuNamespaceLister {
	return gpuNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GpuNamespaceLister helps list and get Gpus.
// All objects returned here must be treated as read-only.
type GpuNamespaceLister interface {
	// List lists all Gpus in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Gpu, err error)
	// Get retrieves the Gpu from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Gpu, error)
	GpuNamespaceListerExpansion
}

// gpuNamespaceLister implements the GpuNamespaceLister
// interface.
type gpuNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Gpus in the indexer for a given namespace.
func (s gpuNamespaceLister) List(selector labels.Selector) (ret []*v1.Gpu, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Gpu))
	})
	return ret, err
}

// Get retrieves the Gpu from the indexer for a given namespace and name.
func (s gpuNamespaceLister) Get(name string) (*v1.Gpu, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("gpu"), name)
	}
	return obj.(*v1.Gpu), nil
}
