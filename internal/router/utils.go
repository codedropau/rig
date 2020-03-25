package router

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// AnnotationDomain used to determine a domain.
	AnnotationDomain = "rig.io/domain"
)

// Get a list of Pods.
func getPodList(clientset kubernetes.Interface) (PodList, error) {
	var list PodList

	pods, err := clientset.CoreV1().Pods(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return list, err
	}

	for _, item := range pods.Items {
		if _, ok := item.ObjectMeta.Annotations[AnnotationDomain]; !ok {
			continue
		}

		pod := Pod{
			Name: item.ObjectMeta.Name,
			Domain: item.ObjectMeta.Annotations[AnnotationDomain],
			IP: item.Status.PodIP,
			Status: item.Status.Phase,
		}

		list.Items = append(list.Items, pod)
	}

	return list, nil
}

// Helper function to get a Pod based on the domain.
func getPod(domain string, list *PodList) (Pod, bool) {
	for _, item := range  list.Items {
		if item.Domain == domain {
			return item, true
		}
	}

	return Pod{}, false
}