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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	testv1 "crdsample/pkg/apis/test/v1"
	versioned "crdsample/pkg/generated/clientset/versioned"
	internalinterfaces "crdsample/pkg/generated/informers/externalversions/internalinterfaces"
	v1 "crdsample/pkg/generated/listers/test/v1"
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// GpuInformer provides access to a shared informer and lister for
// Gpus.
type GpuInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.GpuLister
}

type gpuInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewGpuInformer constructs a new informer for Gpu type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGpuInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGpuInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredGpuInformer constructs a new informer for Gpu type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGpuInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TestV1().Gpus(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TestV1().Gpus(namespace).Watch(context.TODO(), options)
			},
		},
		&testv1.Gpu{},
		resyncPeriod,
		indexers,
	)
}

func (f *gpuInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGpuInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *gpuInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&testv1.Gpu{}, f.defaultInformer)
}

func (f *gpuInformer) Lister() v1.GpuLister {
	return v1.NewGpuLister(f.Informer().GetIndexer())
}
