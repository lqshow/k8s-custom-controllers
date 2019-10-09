package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	corev1 "k8s.io/api/core/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	foobarv1 "github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/apis/foobar/v1"
	clientset "github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/generated/clientset/versioned"
	informers "github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/generated/informers/externalversions/foobar/v1"
	listers "github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/generated/listers/foobar/v1"
)

const controllerAgentName = "foobar-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Network is synced
	SuccessSynced = "Synced"

	// MessageResourceSynced is the message used for an Event fired when a Network
	// is synced successfully
	MessageResourceSynced = "Network synced successfully"
)

// Controller is the controller implementation for Foobar resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeClientset kubernetes.Interface

	// foobarClientset is a clientset for our own API group
	foobarClientset clientset.Interface

	foobarsLister listers.FooBarLister
	foobarsSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new foobar controller
func NewController(
	kubeClientSet kubernetes.Interface,
	foobarClientSet clientset.Interface,
	foobarInformer informers.FooBarInformer) *Controller {

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientSet.CoreV1().Events("")})

	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeClientset:   kubeClientSet,
		foobarClientset: foobarClientSet,
		foobarsLister:   foobarInformer.Lister(),
		foobarsSynced:   foobarInformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "FooBars"),
		recorder:        recorder,
	}

	// Set up an event handler for when FooBar resources change
	// 为 foobarInformer 注册了三个 Handler（AddFunc、UpdateFunc 和 DeleteFunc），分别对应 API 对象的“添加”“更新”和“删除”事件
	// 具体的处理操作，都是将该事件对应的 API 对象加入到工作队列中。
	// 实际入队的并不是 API 对象本身，而是它们的 Key，即：该 API 对象的 <namespace>/<name>
	foobarInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueFooBar,
		UpdateFunc: func(old, new interface{}) {
			oldFooBar := old.(*foobarv1.FooBar)
			newFooBar := new.(*foobarv1.FooBar)
			if oldFooBar.ResourceVersion == newFooBar.ResourceVersion {
				// Periodic resync will send update events for all known FooBars.
				// Two different versions of the same FooBar will always have different RVs.
				return
			}
			controller.enqueueFooBar(new)
		},
		DeleteFunc: controller.enqueueFooBarForDelete,
	})
	return controller
}

// enqueueFooBar takes a FooBar resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than FooBar.
func (c *Controller) enqueueFooBar(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// enqueueFooBarForDelete takes a deleted FooBar resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than FooBar.
func (c *Controller) enqueueFooBarForDelete(obj interface{}) {
	var key string
	var err error
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Network resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Network resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Network resource with this namespace/name
	network, err := c.foobarsLister.FooBars(namespace).Get(name)
	if err != nil {
		// The Network resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			glog.Warningf("Network: %s/%s does not exist in local cache, will delete it from Neutron ...",
				namespace, name)

			glog.Infof("[Neutron] Deleting network: %s/%s ...", namespace, name)

			// FIX ME: call Neutron API to delete this network by name.
			//
			// neutron.Delete(namespace, name)

			return nil
		}

		runtime.HandleError(fmt.Errorf("failed to list network by: %s/%s", namespace, name))

		return err
	}

	glog.Infof("[Neutron] Try to process network: %#v ...", network)

	// FIX ME: Do diff().
	//
	// actualNetwork, exists := neutron.Get(namespace, name)
	//
	// if !exists {
	// 	neutron.Create(namespace, name)
	// } else if !reflect.DeepEqual(actualNetwork, network) {
	// 	neutron.Update(namespace, name)
	// }

	c.recorder.Event(network, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Network control loop")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.foobarsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process Foobar resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}
