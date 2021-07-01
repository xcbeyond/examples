package deployment

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// AppDeploymentSpec is a specification for an app deployment.
type AppDeploymentSpec struct {
	// Name of the application.
	Name string `json:"name"`

	// Docker image path for the application.
	ContainerImage string `json:"containerImage"`

	// The name of an image pull secret in case of a private docker repository.
	ImagePullSecret *string `json:"imagePullSecret"`

	// Command that is executed instead of container entrypoint, if specified.
	ContainerCommand *string `json:"containerCommand"`

	// Arguments for the specified container command or container entrypoint (if command is not
	// specified here).
	ContainerCommandArgs *string `json:"containerCommandArgs"`

	// Number of replicas of the image to maintain.
	Replicas int32 `json:"replicas"`

	// Port mappings for the service that is created. The service is created if there is at least
	// one port mapping.
	PortMappings []PortMapping `json:"portMappings"`

	// List of user-defined environment variables.
	Variables []EnvironmentVariable `json:"variables"`

	// Whether the created service is external.
	IsExternal bool `json:"isExternal"`

	// Description of the deployment.
	Description *string `json:"description"`

	// Target namespace of the application.
	Namespace string `json:"namespace"`

	// Optional memory requirement for the container.
	MemoryRequirement *resource.Quantity `json:"memoryRequirement"`

	// Optional CPU requirement for the container.
	CpuRequirement *resource.Quantity `json:"cpuRequirement"`

	// Labels that will be defined on Pods/RCs/Services
	Labels []Label `json:"labels"`

	// Whether to run the container as privileged user (essentially equivalent to root on the host).
	RunAsPrivileged bool `json:"runAsPrivileged"`
}

// PortMapping is a specification of port mapping for an application deployment.
type PortMapping struct {
	// Port that will be exposed on the service.
	Port int32 `json:"port"`

	// Docker image path for the application.
	TargetPort int32 `json:"targetPort"`

	// IP protocol for the mapping, e.g., "TCP" or "UDP".
	Protocol apiv1.Protocol `json:"protocol"`
}

// EnvironmentVariable represents a named variable accessible for containers.
type EnvironmentVariable struct {
	// Name of the variable. Must be a C_IDENTIFIER.
	Name string `json:"name"`

	// Value of the variable, as defined in Kubernetes core API.
	Value string `json:"value"`
}

// Label is a structure representing label assignable to Pod/RC/Service
type Label struct {
	// Label key
	Key string `json:"key"`

	// Label value
	Value string `json:"value"`
}

// Protocols is a structure representing supported protocol types for a service
type Protocols struct {
	// Array containing supported protocol types e.g., ["TCP", "UDP"]
	Protocols []apiv1.Protocol `json:"protocols"`
}

func CreateDeployment(client *kubernetes.Clientset, spec AppDeploymentSpec) error {
	deploymentClient := client.AppsV1().Deployments(apiv1.NamespaceDefault)

	annotations := map[string]string{}
	if spec.Description != nil {
		annotations["description"] = *spec.Description
	}
	labels := getLabelsMap(spec.Labels)
	objectMeta := metav1.ObjectMeta{
		Annotations: annotations,
		Name:        spec.Name,
		Labels:      labels,
	}

	containerSpec := apiv1.Container{
		Name:  spec.Name,
		Image: spec.ContainerImage,
		SecurityContext: &apiv1.SecurityContext{
			Privileged: &spec.RunAsPrivileged,
		},
		Resources: apiv1.ResourceRequirements{
			Requests: make(map[apiv1.ResourceName]resource.Quantity),
		},
		Env: convertEnvVarsSpec(spec.Variables),
	}

	if spec.ContainerCommand != nil {
		containerSpec.Command = []string{*spec.ContainerCommand}
	}
	if spec.ContainerCommandArgs != nil {
		containerSpec.Args = []string{*spec.ContainerCommandArgs}
	}

	if spec.CpuRequirement != nil {
		containerSpec.Resources.Requests[apiv1.ResourceCPU] = *spec.CpuRequirement
	}
	if spec.MemoryRequirement != nil {
		containerSpec.Resources.Requests[apiv1.ResourceMemory] = *spec.MemoryRequirement
	}
	podSpec := apiv1.PodSpec{
		Containers: []apiv1.Container{containerSpec},
	}
	if spec.ImagePullSecret != nil {
		podSpec.ImagePullSecrets = []apiv1.LocalObjectReference{{Name: *spec.ImagePullSecret}}
	}

	podTemplate := apiv1.PodTemplateSpec{
		ObjectMeta: objectMeta,
		Spec:       podSpec,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: objectMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: &spec.Replicas,
			Template: podTemplate,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}

	_, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

func convertEnvVarsSpec(variables []EnvironmentVariable) []apiv1.EnvVar {
	var result []apiv1.EnvVar
	for _, variable := range variables {
		result = append(result, apiv1.EnvVar{Name: variable.Name, Value: variable.Value})
	}
	return result
}

// Converts array of labels to map[string]string
func getLabelsMap(labels []Label) map[string]string {
	result := make(map[string]string)

	for _, label := range labels {
		result[label.Key] = label.Value
	}

	return result
}
