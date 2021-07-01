package k8sclient

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

// DefaultKubeConfigPath default kubeConfigPath from homeDir
var DefaultKubeConfigPath string

func init()  {
	if home := homedir.HomeDir(); home != "" {
		DefaultKubeConfigPath = *flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		DefaultKubeConfigPath = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
}

func NewK8sClient(kubeConfigPath, apiserverHost string) (*kubernetes.Clientset){
	config, err := clientcmd.BuildConfigFromFlags(apiserverHost, kubeConfigPath)
	if err != nil {
		panic(err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return k8sClient
}

