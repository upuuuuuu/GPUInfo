package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	Gpuv1 "github.com/upuuuuuu/gpuinfo/pkg/apis/test/v1"
	clientset "github.com/upuuuuuu/gpuinfo/pkg/generated/clientset/versioned"
	informerv1 "github.com/upuuuuuu/gpuinfo/pkg/generated/informers/externalversions/test/v1"
	listerv1 "github.com/upuuuuuu/gpuinfo/pkg/generated/listers/test/v1"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
)

const (
	numWork  = 5
	maxRetry = 10
)

// struct usually define the func return type
// the func return Lister
type controller struct {
	client clientset.Interface
	// Lister = indexer
	gpuLister listerv1.GpuLister
	queue     workqueue.RateLimitingInterface
	nodeName  string
}

// get obj turn into key and add key to queue
func (c *controller) enQueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.queue.Add(key)
}

func (c *controller) addFun(obj interface{}) {
	fmt.Println("This is a add gpu func")
	c.enQueue(obj)
}

func (c *controller) updateFun(oldobj, newobj interface{}) {
	// if reflect.DeepEqual(oldobj, newobj) {
	// 	fmt.Println("This ia a update gpu but can not update")
	// 	return
	// }
	fmt.Println("This is a update gpu")
	c.enQueue(newobj)
}

func (c *controller) Run(stopCh chan struct{}) {

	c.nodeName = os.Getenv("NODE_NAME")
	if err := c.createGpu(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < numWork; i++ {
		fmt.Println("The Run func is Start")
		go wait.Until(c.worker, time.Minute, stopCh)
	}
	<-stopCh
}

func (c *controller) createGpu() error {

	var replicas int32 = 1
	gpu := Gpuv1.Gpu{
		ObjectMeta: v1.ObjectMeta{
			Name: c.nodeName,
		},
		Spec: Gpuv1.GpuSpec{
			DeploymentName: "example-gpu",
			Replicas:       &replicas,
			Memory:         1024,
		},
	}

	gpucreated, err := c.client.TestV1().Gpus("default").Create(context.TODO(), &gpu, v1.CreateOptions{})
	if err != nil {
		return err
	} else {
		fmt.Println("gpucreatedname: ", gpucreated.Name)
	}

	return nil
}

func (c *controller) worker() {
	fmt.Println("The Worker func is Start")
	for c.processNextItem() {
	}
}

func (c *controller) processNextItem() bool {
	fmt.Println("The processNextItem func is Start")
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	defer c.queue.Done(item)

	key := item.(string)

	err := c.SynGpu(key)
	if err != nil {
		c.handlerError(key, err)
	}
	return true
}

func (c *controller) handlerError(key string, err error) {
	if c.queue.NumRequeues(key) <= maxRetry {
		c.queue.AddRateLimited(key)
		return
	}
	runtime.HandleError(err)
	c.queue.Forget(key)
}

func (c *controller) SynGpu(key string) error {
	fmt.Println("The SynGpu func is Start")
	// key build by ns and name
	nameSpaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// qishibuyongxuan next zhijeiyongnodename
	if name != c.nodeName {
		fmt.Println("This is not my event")
		return nil

	}

	gpu, err := c.client.TestV1().Gpus(nameSpaceKey).Get(context.TODO(), name, v1.GetOptions{})

	nvml.Init()
	defer nvml.Shutdown()
	count, err := nvml.GetDeviceCount()
	if err != nil {
		return err
	}
	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		status, err := device.Status()
		if err != nil {
			return err
		}
		gpu.Spec.Memory = int32(*status.Memory.Global.Free)
	}

	newgpu, err := c.client.TestV1().Gpus(nameSpaceKey).Update(context.TODO(), gpu, v1.UpdateOptions{})

	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Println("update", newgpu.Spec.Memory)

	return nil
}

// in: sevice infomer (fake)
// out: not infomer.lister is contorler but chabuduo ; why return informer.lister
func NewController(client clientset.Interface, gpuInfomer informerv1.GpuInformer) controller {
	c := controller{
		client:    client,
		gpuLister: gpuInfomer.Lister(),
		queue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "GpuInfo"),
	}

	// add event keneng ye add dao queue queue from c
	gpuInfomer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addFun,
		UpdateFunc: c.updateFun,
	})

	return c
}
