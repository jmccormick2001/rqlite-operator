package rqcluster

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	rqclusterv1alpha1 "github.com/jmccormick2001/rq/pkg/apis/rqcluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// rqReconcile implements the Reconcile for the rq-operator
func rqReconcile(r *ReconcileRqcluster, request reconcile.Request, instance *rqclusterv1alpha1.Rqcluster) error {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	reqLogger.Info("rqReconcile called")

	podList, err := getPods(r, request.Namespace, instance.Name)
	if err != nil {
		return err
	}

	err = verifyServices(r, instance)
	if err != nil {
		return err
	}

	requestedPodCount := int(instance.Spec.Size)
	podCount := len(podList.Items)

	if podCount != requestedPodCount {

		//handle the case of a new cluster, we need a leader
		//pod to be created first before creating the followers
		if podCount == 0 {
			err := createClusterPod(true, r, instance)
			if err != nil {
				return err
			}

			// Delay a bit to let the leader start before
			// the followers
			time.Sleep(time.Duration(2) * time.Second)

			podCount += 1
		}

		podsToCreate := requestedPodCount - podCount
		for i := 0; i < podsToCreate; i++ {
			err := createClusterPod(false, r, instance)
			if err != nil {
				return err
			}
		}

		// check for the case where a leader pod has been removed
		leaderPod, err := getLeaderPod(r, request.Namespace, instance.Name)
		if err != nil {
			return err
		}
		if leaderPod == nil {
			reqLogger.Info("would need to see who the new leader is here")
			err := labelNewLeader(r, instance)
			if err != nil {
				return err
			}
		}

	}

	// at this point, the cluster's pods should exist
	return updateStatus(podList.Items, r, instance)
}

// getPods returns the list of pods for a given namespace and instance
func getPods(r *ReconcileRqcluster, requestNamespace, instanceName string) (*corev1.PodList, error) {
	reqLogger := log.WithValues("Request.Namespace", requestNamespace, "Request.Name", instanceName)

	podList := &corev1.PodList{}
	err := r.client.List(context.TODO(), podList, client.InNamespace(requestNamespace), client.MatchingLabels{"cluster": instanceName})
	if err != nil {
		reqLogger.Info("unable to find any pods that match this request: " + err.Error())
		return podList, err
	}

	return podList, nil
}

func createClusterPod(leader bool, r *ReconcileRqcluster, instance *rqclusterv1alpha1.Rqcluster) error {
	reqLogger := log.WithValues("Request.Namespace", instance.Namespace, "Request.Name", instance.Name)

	reqLogger.Info("createClusterPod called")
	var joinAddress string
	if !leader {
		joinAddress = fmt.Sprintf("--join http://%s-leader:4001", instance.Name)
	}
	// define a new cluster Pod
	mypod, err := newPodForCRFromTemplate(joinAddress, instance, r.client)
	if err != nil {
		return err
	}

	// Create a service for the pod if it doesn't exist
	svcfound := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: mypod.Name, Namespace: mypod.Namespace}, svcfound)
	if err != nil && errors.IsNotFound(err) {
		podSvc, err := newServiceForPod(mypod.Name, instance, r.client)
		if err != nil {
			return err
		}
		// Set Rqcluster instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, podSvc, r.scheme); err != nil {
			return err
		}
		err = r.client.Create(context.TODO(), podSvc)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Set Rqcluster instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, mypod, r.scheme); err != nil {
		return err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: mypod.Name, Namespace: mypod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", mypod.Namespace, "Pod.Name", mypod.Name, "Namespace", mypod.ObjectMeta.Namespace)
		if leader {
			mypod.ObjectMeta.Labels["leader"] = "true"
		}
		err = r.client.Create(context.TODO(), mypod)
		if err != nil {
			return err
		}

		// Pod created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

// verifyServices checks to see if there is ...
// a service for the cluster leader
// a service that will select on all pods in the cluster
func verifyServices(r *ReconcileRqcluster, instance *rqclusterv1alpha1.Rqcluster) error {
	reqLogger := log.WithValues("Request.Namespace", instance.Namespace, "Request.Name", instance.Name)

	// Check if the leader service already exists
	var leaderService *corev1.Service
	leaderService, err := newServiceForCRFromTemplate(true, instance, r.client)
	if err != nil {
		return err
	}

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: leaderService.Name, Namespace: leaderService.Namespace}, leaderService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new leader service", "Pod.Namespace", leaderService.Namespace, "Pod.Name", leaderService.Name, "Namespace", leaderService.ObjectMeta.Namespace)

		// Set Rqcluster instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, leaderService, r.scheme); err != nil {
			return err
		}
		err = r.client.Create(context.TODO(), leaderService)
		if err != nil {
			return err
		}

		// leader Service created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	} else {
		reqLogger.Info("leader service already exists")
	}

	// Check if the rqcluster service already exists
	var clusterService *corev1.Service
	clusterService, err = newServiceForCRFromTemplate(false, instance, r.client)
	if err != nil {
		return err
	}

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: clusterService.Name, Namespace: clusterService.Namespace}, clusterService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new cluster service", "Pod.Namespace", clusterService.Namespace, "Pod.Name", clusterService.Name, "Namespace", clusterService.ObjectMeta.Namespace)

		// Set Rqcluster instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, clusterService, r.scheme); err != nil {
			return err
		}
		err = r.client.Create(context.TODO(), clusterService)
		if err != nil {
			return err
		}

		// cluster Service created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	} else {
		reqLogger.Info("cluster service already exists")
	}

	return nil
}

func updateStatus(pods []corev1.Pod, r *ReconcileRqcluster, instance *rqclusterv1alpha1.Rqcluster) error {
	reqLogger := log.WithValues("Request.Namespace", instance.Namespace, "Request.Name", instance.Name)
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("Failed to update rqcluster status: " + err.Error())
			//return err
			// I'm returning nil here per https://github.com/kubernetes-sigs/controller-runtime/issues/403
			if errors.IsConflict(err) {
				reqLogger.Info("conflict error raised in status update: " + err.Error())
			}

			return nil
		}
	}
	return nil

}
