package controller

import (
	v1beta12 "k8s.io/api/networking/v1beta1"
	v16 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v13 "k8s.io/client-go/informers/core/v1"
	v14 "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	v15 "k8s.io/client-go/listers/core/v1"
	v1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	client        kubernetes.Interface
	ingressLister v1.IngressLister
	serviceLister v15.ServiceLister
	queue         workqueue.RateLimitingInterface
}

func (c *controller) addService(obj interface{}) {
	c.addqueue(obj)
}

func (c *controller) updateService(oldobj interface{}, newObje interface{}) {
	//对比 TODO
	//不相同
	c.addqueue(newObje)
}

func (c controller) deleteIngress(obj interface{}) {
	ingre := obj.(v1beta12.Ingress)
	service := v16.GetControllerOf(ingre)
	if nil == service || service.Kind != "service" {
		return
	}

	c.addqueue(service)
}

func (c controller) Run(stopCh <-chan struct{}) {
	select {
	case <-stopCh:
		return
	default:
		item, _ := c.queue.Get()
		ns, name, _ := cache.SplitMetaNamespaceKey(item.(string))
		//TODO 	操作svc和ingress
	}

}

func (c *controller) deleteService(obj interface{}) {
	//service := obj.(v12.Service)
	// TODO
	//get annotation and delete ingress
}

func (c *controller) addqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if nil != err {
		return
	}

	c.queue.Add(key)

}

func NewController(client kubernetes.Interface, serviceInfor v13.ServiceInformer, ingressInformer v14.IngressInformer) controller {
	c := controller{
		client:        client,
		ingressLister: ingressInformer.Lister(),
		serviceLister: serviceInfor.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}

	serviceInfor.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
		DeleteFunc: c.deleteService,
	})
	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})
	return c
}
