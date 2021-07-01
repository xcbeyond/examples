package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type clientManager struct {
	// Path to kubeconfig file. If both kubeConfigPath and apiserverHost are empty
	// inClusterConfig will be used
	kubeConfigPath string
	// Address of apiserver host in format 'protocol://address:port'
	apiserverHost string

	// Kubernetes client created without providing auth info.
	insecureClient kubernetes.Interface
}

func NewClientManager(kubeConfigPath, apiserverHost string) clientManager {
	client := clientManager{
		kubeConfigPath: kubeConfigPath,
		apiserverHost:  apiserverHost,
	}
	return client
}

func (client *clientManager) Client() {
	config, err := clientcmd.BuildConfigFromFlags(client.apiserverHost, client.kubeConfigPath)
	if err != nil {
		panic(err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	client.insecureClient = k8sClient
}
