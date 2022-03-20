package controller

import (
	v13 "k8s.io/client-go/informers/core/v1"
	v14 "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	v15 "k8s.io/client-go/listers/core/v1"
	v1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
)

type controller struct {
	client kubernetes.Interface
	ingressLister v1.IngressLister
	serviceLister v15.ServiceLister 
}

func (c *controller)addService(obj interface{}) {

}

func (c *controller)updateService(oldobj interface{}, newObje interface{}) {

}

func (c controller)deleteIngress(obj interface{}){

}

func (c controller)Run(stopCh <-chan struct{}){
	select {
		case  <- stopCh:
			return


	}

}

func NewController(client kubernetes.Interface, serviceInfor v13.ServiceInformer, ingressInformer v14.IngressInformer)controller{
	c :=  controller{
		client:       client,
		ingressLister: ingressInformer.Lister(),
		serviceLister: serviceInfor.Lister(),
	}
	
	serviceInfor.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
		DeleteFunc: nil,
	})
	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})
	return c
}