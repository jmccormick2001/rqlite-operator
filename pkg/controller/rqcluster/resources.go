package rqcluster

import (
	rqclusterv1alpha1 "github.com/jmccormick2001/rq/pkg/apis/rqcluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// setResources sets the cpu and memory requests and limits from the cr
func setResources(cr *rqclusterv1alpha1.Rqcluster, client client.Client, pod *corev1.Pod) error {

	if cr.Spec.MemoryLimit == "" && cr.Spec.CpuLimit == "" && cr.Spec.MemoryRequest == "" && cr.Spec.CpuRequest == "" {
		return nil
	}

	myLimits := corev1.ResourceList{}
	if cr.Spec.MemoryLimit != "" {
		memLimit, err := resource.ParseQuantity(cr.Spec.MemoryLimit)
		if err != nil {
			return err
		}
		myLimits[corev1.ResourceMemory] = memLimit
	}

	if cr.Spec.CpuLimit != "" {
		cpuLimit, err := resource.ParseQuantity(cr.Spec.CpuLimit)
		if err != nil {
			return err
		}
		myLimits[corev1.ResourceCPU] = cpuLimit
	}

	myRequests := corev1.ResourceList{}
	if cr.Spec.MemoryRequest != "" {
		memRequest, err := resource.ParseQuantity(cr.Spec.MemoryRequest)
		if err != nil {
			return err
		}
		myRequests[corev1.ResourceMemory] = memRequest
	}

	if cr.Spec.CpuRequest != "" {
		cpuRequest, err := resource.ParseQuantity(cr.Spec.CpuRequest)
		if err != nil {
			return err
		}
		myRequests[corev1.ResourceCPU] = cpuRequest
	}

	pod.Spec.Containers[0].Resources = corev1.ResourceRequirements{
		Limits:   myLimits,
		Requests: myRequests,
	}

	return nil
}
