package deployment

import (
	"k8s-client-go-examples/util/k8sclient"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestCreateDeployment(t *testing.T) {
	type args struct {
		client *kubernetes.Clientset
		spec   AppDeploymentSpec
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{
			client: k8sclient.NewK8sClient(k8sclient.DefaultKubeConfigPath, ""),
			spec:   AppDeploymentSpec{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateDeployment(tt.args.client, tt.args.spec); (err != nil) != tt.wantErr {
				t.Errorf("CreateDeployment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
