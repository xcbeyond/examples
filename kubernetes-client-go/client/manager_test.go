package client

import (
	"flag"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

// default kubeConfigPath,from homeDir
var defaultKubeConfigPath string

func init()  {
	if home := homedir.HomeDir(); home != "" {
		defaultKubeConfigPath = *flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		defaultKubeConfigPath = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
}


