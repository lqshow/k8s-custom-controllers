package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/controller"
	"github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/kube"
	"go.uber.org/zap"

	informers "github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/generated/informers/externalversions"
)

var (
	masterURL  string
	kubeconfig string

	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func registerSIGINTHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := registerSIGINTHandler()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// 创建 Kubernetes client 和 FooBar client
	kubeClient, foobarClient := kube.GetKubernetesClient(masterURL, kubeconfig)

	// 为 FooBar 对象创建一个叫作 InformerFactory 的工厂，并使用它生成一个 FooBar 对象的 Informer，传递给控制器。
	foobarInformerFactory := informers.NewSharedInformerFactory(foobarClient, time.Second*30)

	//foo, e := foobarInformerFactory.Samplecrd().V1alpha1().FooBars().Lister().FooBars("dev").Get("foobar-sample")
	//zap.S().Infof("foobar sample: %v", foo, e)

	// 启动上述的 Informer，然后执行 controller.Run，启动自定义控制器
	foobarController := controller.NewController(kubeClient, foobarClient,
		foobarInformerFactory.Samplecrd().V1alpha1().FooBars())

	if err := foobarController.Run(2, stopCh); err != nil {
		zap.S().Fatalf("Error running controller: %s", err.Error())
	}
}
