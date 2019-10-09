package main

import (
	"github.com/golang/glog"
	controller "github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/controller"
	"github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/kube"
	"os"
	"os/signal"
	"syscall"
	"time"


	informers "github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/generated/informers/externalversions"
)

func registerSIGINTHandler() (stopCh <-chan struct{}){
	// Register for SIGINT.
	stop := make(chan struct{})
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		close(stop)
		<-signalChan
		os.Exit(1)
	}()

	return stop
}

func main() {
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := registerSIGINTHandler()

	// 创建 Kubernetes client 和 FooBar client
	kubeClient, foobarClient := kube.GetKubernetesClient()

	// 为 FooBar 对象创建一个叫作 InformerFactory 的工厂，并使用它生成一个 FooBar 对象的 Informer，传递给控制器。
	foobarInformerFactory := informers.NewSharedInformerFactory(foobarClient, time.Second*30)

	// 启动上述的 Informer，然后执行 controller.Run，启动自定义控制器
	foobarController := controller.NewController(kubeClient, foobarClient,
		foobarInformerFactory.K8s().V1().FooBars())

	if err := foobarController.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}