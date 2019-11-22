package k8s

import (
	"fmt"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	corev1 "k8s.io/api/core/v1"

	stringutils "github.com/nickschuch/rig/internal/utils/string"
	composeconfig "github.com/nickschuch/rig/internal/compose/config"
	"github.com/nickschuch/rig/internal/config"
)

const (
	AnnotationTTL = "rig.io/ttl"
	LabelName = "app.kubernetes.io/name"
	LabelInstance = "app.kubernetes.io/instance"
	LabelVersion = "app.kubernetes.io/version"
	LabelManagedBy = "app.kubernetes.io/managed-by"
)

type Params struct {
	Project string

	Namespace string
	Name string

	Repository string
	Tag string

	Domains []string

	// Compose is a loaded Docker Compose configuration which is used to build a Pod.
	Compose *composeconfig.Config
	Config config.Config
}

func Apply(clientset kubernetes.Interface, params Params) error {
	grace := int64(0)

	objectmeta := metav1.ObjectMeta{
		Name: params.Name,
		Namespace: params.Namespace,
		Annotations: map[string]string{
			// @todo Configurable.
			AnnotationTTL: "24h",
		},
		// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels
		Labels: map[string]string{
			LabelName: params.Name,
			LabelInstance: fmt.Sprintf("%s-%s", params.Project, params.Name),
			LabelVersion: params.Tag,
			LabelManagedBy: "rig",
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: objectmeta,
	}

	for name := range params.Compose.Volumes {
		// @todo, Filter the volumes.

		pod.Spec.InitContainers = append(pod.Spec.InitContainers, corev1.Container{
			Name:                     fmt.Sprintf("volume-%s-cp", name),
			Image:                    fmt.Sprintf("%s:%s-volume-%s", params.Repository, params.Tag, name),
			Command:                  []string{
				"cp", "-rp", ".", "/mnt/volume/",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name: name,
					MountPath: "/mnt/volume",
				},
			},
		})

		pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
			Name:         name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	}

	for name, service := range params.Compose.Services {
		if !stringutils.Contains(params.Config.Services, name) {
			continue
		}

		container :=  corev1.Container{
			Name:                     name,
			Image:                    fmt.Sprintf("%s:%s-service-%s", params.Repository, params.Tag, name),
			// @todo, Requires resource constraints.
		}

		for _, environment := range service.Environment {
			envName, envValue, err := stringutils.SplitBySeparator(environment, "=")
			if err != nil {
				return err
			}

			container.Env = append(container.Env, corev1.EnvVar{
				Name:      envName,
				Value:     envValue,
				ValueFrom: nil,
			})
		}

		for _, volume := range service.Volumes {
			volName, volValue, err := stringutils.SplitBySeparator(volume, ":")
			if err != nil {
				return err
			}

			container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
				Name:             volName,
				MountPath:        volValue,
			})
		}

		pod.Spec.Containers = append(pod.Spec.Containers, container)
	}

	_, err := clientset.CoreV1().Pods(objectmeta.Namespace).Create(pod)
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			err = clientset.CoreV1().Pods(objectmeta.Namespace).Delete(objectmeta.Name, &metav1.DeleteOptions{
				GracePeriodSeconds: &grace,
			})
			if err != nil {
				return err
			}

			_, err = clientset.CoreV1().Pods(objectmeta.Namespace).Create(pod)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	service := &corev1.Service{
		ObjectMeta: objectmeta,
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					// @todo, Needs to be configurable.
					Port: 8080,
				},
			},
			Selector: objectmeta.Labels,
		},
	}

	_, err = clientset.CoreV1().Services(objectmeta.Namespace).Create(service)
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			err = clientset.CoreV1().Services(objectmeta.Namespace).Delete(objectmeta.Name, &metav1.DeleteOptions{
				GracePeriodSeconds: &grace,
			})
			if err != nil {
				return err
			}

			_, err = clientset.CoreV1().Services(objectmeta.Namespace).Create(service)
			if err != nil {
				return err
			}

			return nil
		} else {
			return err
		}
	}

	ingress := &networkingv1beta1.Ingress{
		ObjectMeta: objectmeta,
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
								ServiceName: objectmeta.Name,
								// @todo, Needs configuration.
								ServicePort: intstr.FromInt(8080),
							},
						},
					},
				},
			},
		})
	}

	_, err = clientset.NetworkingV1beta1().Ingresses(objectmeta.Namespace).Create(ingress)
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			_, err = clientset.NetworkingV1beta1().Ingresses(objectmeta.Namespace).Update(ingress)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
