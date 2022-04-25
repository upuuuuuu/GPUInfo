package main

import (
	"context"
	clientset "crdsample/pkg/generated/clientset/versioned"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func sche() {

	type sche struct {
		client   clientset.Interface
		queue    workqueue.RateLimitingInterface
		nodeName string
	}

	var c *sche

	item, shutdown := c.queue.Get()
	if shutdown {
		return
	}

	key := item.(string)

	nameSpaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return
	}
	gpu, err := c.client.TestV1().Gpus("default").Get(context.TODO(), c.nodeName, v1.GetOptions{})

	fmt.Println(nameSpaceKey, name, gpu)
	defer c.queue.Done(item)

}
