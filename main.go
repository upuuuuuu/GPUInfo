package main

import (
	"log"
	"time"

	controller "github.com/upuuuuuu/gpuinfo/pkg/controller"
	clientset "github.com/upuuuuuu/gpuinfo/pkg/generated/clientset/versioned"
	informers "github.com/upuuuuuu/gpuinfo/pkg/generated/informers/externalversions"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. config
	configPath := "admin.conf"
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Fatal(err)
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
		config = inClusterConfig
	}

	// 2. clientset
	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// gpuList, err := clientset.TestV1().Gpus("default").List(context.TODO(), v1.ListOptions{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, gpu := range gpuList.Items {
	// 	fmt.Println(gpu.Spec.Memory)
	// }

	// 3. factory informer
	informerFactory := informers.NewSharedInformerFactory(clientset, 1*time.Second)

	GpuInformer := informerFactory.Test().V1().Gpus()

	// 4. eventHandler --> controller
	controller := controller.NewController(clientset, GpuInformer)

	// 5. start
	stopCh := make(chan struct{})
	informerFactory.Start(stopCh)
	informerFactory.WaitForCacheSync(stopCh)

	controller.Run(stopCh)

}

// informers.AddEventHandler(cache.ResourceEventHandlerFuncs{
// 	AddFunc:    onAdd,
// 	UpdateFunc: onUpadate,
// 	DeleteFunc: onDelete,
// })

// func onAdd(obj interface{}) {
// 	gpu := obj.(*gpuv1.Gpu)
// 	fmt.Println("add", gpu.Name)
// }
// func onUpadate(old, new interface{}) {
// 	oldgpu := old.(*gpuv1.Gpu)
// 	newgpu := new.(*gpuv1.Gpu)

// 	fmt.Println("update", newgpu.Spec.Memory, oldgpu.Spec.Memory)

// }
// func onDelete(obj interface{}) {
// 	gpu := obj.(*gpuv1.Gpu)
// 	fmt.Println("delete", gpu.Name)

// }
