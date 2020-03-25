package router

import corev1 "k8s.io/api/core/v1"

// PodList is returned containing a list of Pods.
type PodList struct {
	Items []Pod
}

// Pod which is routable by Rig.
type Pod struct {
	Name string
	Domain string
	IP string
	Status corev1.PodPhase
}