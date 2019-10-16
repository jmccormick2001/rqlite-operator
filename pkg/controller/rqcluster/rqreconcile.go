package rqcluster

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
	//logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	rqclusterv1alpha1 "github.com/jmccormick2001/rq/pkg/apis/rqcluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//var log = logf.Log.WithName("rqReconcile")

// rqReconcile implements the Reconcile for the rq-operator
func rqReconcile(r *ReconcileRqcluster, request reconcile.Request, instance *rqclusterv1alpha1.Rqcluster) error {

	log.Info("jeff Reconciling Rqcluster")

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
			//a not so great way to let the leader get started
			//before creating the followers
			time.Sleep(time.Duration(4) * time.Second)

			podCount += 1
		}

		podsToCreate := requestedPodCount - podCount
		for i := 0; i < podsToCreate; i++ {
			err := createClusterPod(false, r, instance)
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
	podList := &corev1.PodList{}
	err := r.client.List(context.TODO(), podList, client.InNamespace(requestNamespace), client.MatchingLabels{"cluster": instanceName})
	if err != nil {
		log.Error(err, "unable to find any pods that match this request")
		return podList, err
	}

	return podList, nil
}

func createClusterPod(leader bool, r *ReconcileRqcluster, instance *rqclusterv1alpha1.Rqcluster) error {

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
		log.Info("Creating a new Pod", "Pod.Namespace", mypod.Namespace, "Pod.Name", mypod.Name, "Namespace", mypod.ObjectMeta.Namespace)
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

	// Check if the leader service already exists
	leaderStatus := []bool{true, false}
	for _, v := range leaderStatus {
		leaderService, err := newServiceForCRFromTemplate(v, instance, r.client)
		if err != nil {
			return err
		}
		found := &corev1.Service{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: leaderService.Name, Namespace: leaderService.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new leader service", "Pod.Namespace", leaderService.Namespace, "Pod.Name", leaderService.Name, "Namespace", leaderService.ObjectMeta.Namespace)

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
		}
	}

	return nil
}

func updateStatus(pods []corev1.Pod, r *ReconcileRqcluster, instance *rqclusterv1alpha1.Rqcluster) error {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			log.Error(err, "Failed to update rqcluster status")
			return err
		}
	}
	return nil

}
