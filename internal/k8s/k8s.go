package k8s

import (
	"fmt"
	"github.com/codedropau/rig/cmd/rig/version"
	"k8s.io/apimachinery/pkg/api/resource"

	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	composeconfig "github.com/codedropau/rig/internal/compose/config"
	"github.com/codedropau/rig/internal/config"
	stringutils "github.com/codedropau/rig/internal/utils/string"
)

// ManagedBy is used to determine if this was a Rig environment.
// @todo, Consider moving this to a separate package eg. Could be helpful if we can to "list" environments.
const ManagedBy = "rig"

const (
	// AnnotationTTL is used to define how long an environment should be available for.
	AnnotationTTL = "rig.io/ttl"
	// AnnotationRigVersion is used to determine what version of Rig provisioned his resource.
	AnnotationRigVersion = "rig.io/version"
	// AnnotationRigCommit is used to determine what commit of Rig provisioned his resource.
	AnnotationRigCommit = "rig.io/commit"
)

const (
	// LabelName is used to discovery an object.
	LabelName = "app.kubernetes.io/name"
	// LabelInstance is used to discovery an object.
	LabelInstance = "app.kubernetes.io/instance"
	// LabelVersion is used to discovery an object.
	LabelVersion = "app.kubernetes.io/version"
	// LabelManagedBy is used to discovery an object.
	LabelManagedBy = "app.kubernetes.io/managed-by"
)

const (
	// ResourceRequestCPU which are applied to each container.
	ResourceRequestCPU = "50m"
	// ResourceRequestMemory which are applied to each container.
	ResourceRequestMemory = "128Mi"
	// ResourceLimitDefaultCPU is applied to each container when not specified.
	ResourceLimitDefaultCPU = "150m"
	// ResourceLimitDefaultMemory is applied to each container when not specified.
	ResourceLimitDefaultMemory = "256Mi"
)

const (
	// MountVolume for copying volume contents into an EmptyDir.
	MountVolume = "/mnt/volume"
)

// Params used to create Kubernetes objects.
type Params struct {
	Project string

	// Metadata applied to Kubernetes objects.
	Namespace string
	Name      string

	// Information used to run the correct images.
	Repository string
	Tag        string

	// Domains which the environment will be accessible from.
	Domains []string

	// Compose is a loaded Docker Compose configuration which is used to build a Pod.
	Compose *composeconfig.Config
	Config  *config.Config
}

// Apply objects to the Kubernetes cluster.
func Apply(clientset kubernetes.Interface, params Params) error {
	grace := int64(0)

	metadata := metav1.ObjectMeta{
		Name:      params.Name,
		Namespace: params.Namespace,
		Annotations: map[string]string{
			AnnotationTTL:        params.Config.Retention.String(),
			AnnotationRigVersion: version.GitVersion,
			AnnotationRigCommit:  version.GitCommit,
		},
		// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels
		Labels: map[string]string{
			LabelName:      params.Name,
			LabelInstance:  fmt.Sprintf("%s-%s", params.Project, params.Name),
			LabelVersion:   params.Tag,
			LabelManagedBy: ManagedBy,
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: metadata,
	}

	for name := range params.Compose.Volumes {
		// @todo, Filter the volumes.

		pod.Spec.InitContainers = append(pod.Spec.InitContainers, corev1.Container{
			Name:  fmt.Sprintf("volume-%s-cp", name),
			Image: fmt.Sprintf("%s:%s-volume-%s", params.Repository, params.Tag, name),
			Command: []string{
				"/bin/sh", "-c",
			},
			Args: []string{
				fmt.Sprintf("cp -rp . %s/", MountVolume),
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      name,
					MountPath: MountVolume,
				},
			},
		})

		pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	}

	for name, service := range params.Compose.Services {
		if _, ok := params.Config.Services[name]; !ok {
			continue
		}

		cpu, err := resourceWithFallback(params.Config.Services[name].CPU, ResourceLimitDefaultCPU)
		if err != nil {
			return fmt.Errorf("failed to parse cpu: %w", err)
		}

		memory, err := resourceWithFallback(params.Config.Services[name].Memory, ResourceLimitDefaultMemory)
		if err != nil {
			return fmt.Errorf("failed to parse memory: %w", err)
		}

		container := corev1.Container{
			Name:  name,
			Image: fmt.Sprintf("%s:%s-service-%s", params.Repository, params.Tag, name),
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse(ResourceRequestCPU),
					corev1.ResourceMemory: resource.MustParse(ResourceRequestMemory),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    cpu,
					corev1.ResourceMemory: memory,
				},
			},
		}

		for _, environment := range service.Environment {
			envName, envValue, err := stringutils.SplitBySeparator(environment, "=")
			if err != nil {
				return err
			}

			container.Env = append(container.Env, corev1.EnvVar{
				Name:  envName,
				Value: envValue,
			})
		}

		for _, volume := range service.Volumes {
			volName, volValue, err := stringutils.SplitBySeparator(volume, ":")
			if err != nil {
				return err
			}

			container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
				Name:      volName,
				MountPath: volValue,
			})
		}

		pod.Spec.Containers = append(pod.Spec.Containers, container)
	}

	_, err := clientset.CoreV1().Pods(metadata.Namespace).Create(pod)
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			err = clientset.CoreV1().Pods(metadata.Namespace).Delete(metadata.Name, &metav1.DeleteOptions{
				GracePeriodSeconds: &grace,
			})
			if err != nil {
				return err
			}

			_, err = clientset.CoreV1().Pods(metadata.Namespace).Create(pod)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	service := &corev1.Service{
		ObjectMeta: metadata,
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Port: int32(params.Config.Ingress.Port),
				},
			},
			Selector: metadata.Labels,
		},
	}

	_, err = clientset.CoreV1().Services(metadata.Namespace).Create(service)
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			err = clientset.CoreV1().Services(metadata.Namespace).Delete(metadata.Name, &metav1.DeleteOptions{
				GracePeriodSeconds: &grace,
			})
			if err != nil {
				return err
			}

			_, err = clientset.CoreV1().Services(metadata.Namespace).Create(service)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	ingress := &networkingv1beta1.Ingress{
		ObjectMeta: metadata,
	}

	for _, domain := range params.Domains {
		ingress.Spec.Rules = append(ingress.Spec.Rules, networkingv1beta1.IngressRule{
			Host: domain,
			IngressRuleValue: networkingv1beta1.IngressRuleValue{
				HTTP: &networkingv1beta1.HTTPIngressRuleValue{
					Paths: []networkingv1beta1.HTTPIngressPath{
						{
							Path: "/",
							Backend: networkingv1beta1.IngressBackend{
								ServiceName: metadata.Name,
								ServicePort: intstr.FromInt(params.Config.Ingress.Port),
							},
						},
					},
				},
			},
		})
	}

	_, err = clientset.NetworkingV1beta1().Ingresses(metadata.Namespace).Create(ingress)
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			_, err = clientset.NetworkingV1beta1().Ingresses(metadata.Namespace).Update(ingress)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
