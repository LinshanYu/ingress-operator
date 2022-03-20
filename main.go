package main

import (
	"github.com/linshanyu/ingress-operator/controller"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if nil != err {
		clusterConfig, err := rest.InClusterConfig()
		if nil != err {
			panic(err)
		}
		config = clusterConfig
	}

	clientset, err := kubernetes.NewForConfig(config)
	if nil != err {
		panic(err)
	}
	infoFactory := informers.NewSharedInformerFactory(clientset, 0)
	servicesInfo := infoFactory.Core().V1().Services()
	ingressInfo := infoFactory.Networking().V1().Ingresses()
	//ForCrd
	//geInfor,err := infoFactory.ForResource()
	//if nil != err {
	//	panic(err)
	//}

	con := controller.NewController(clientset, servicesInfo, ingressInfo)
	stopChan := make(chan struct{})
	infoFactory.Start(stopChan)
	con.Run(stopChan)
}
