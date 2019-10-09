package kube

import (
	"k8s.io/client-go/rest"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	//apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	//apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	//apierrors "k8s.io/apimachinery/pkg/api/errors"
	//meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clientset "github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/generated/clientset/versioned"
)

const (
	CRDPlural   string = "sslconfigs"
	CRDGroup    string = "blog.velotio.com"
	CRDVersion  string = "v1alpha1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

func GetKubeConfig() (*rest.Config, error) {
	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"

	// create the config from the path
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)

	return cfg, err
}

// retrieve the Kubernetes cluster client from outside of the cluster
func GetKubernetesClient() (kubernetes.Interface, clientset.Interface) {
	cfg, err := GetKubeConfig()
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	// generate the client based off of the config
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	foobarClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Print("Successfully constructed k8s client")
	return kubeClient, foobarClient
}

//func CreateCRD(clientset apiextension.Interface) error {
//	crd := &apiextensionv1beta1.CustomResourceDefinition{
//		ObjectMeta: meta_v1.ObjectMeta{Name: FullCRDName},
//		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
//			Group:   CRDGroup,
//			Version: CRDVersion,
//			Scope:   apiextensionv1beta1.NamespaceScoped,
//			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
//				Plural: CRDPlural,
//				Kind:   reflect.TypeOf(SslConfig{}).Name(),
//			},
//		},
//	}
//
//	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
//	if err != nil && apierrors.IsAlreadyExists(err) {
//		return nil
//	}
//	return err
//}