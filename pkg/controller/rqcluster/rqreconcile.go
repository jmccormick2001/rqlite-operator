package rqcluster

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

	if len(podList.Items) == 0 {
		for i := 0; i < 3; i++ {
			err := createClusterPod(r, instance)
			if err != nil {
				return err
			}
		}
	}

	// cluster Pods already exists
	log.Info("jeff reconcile: here is where we handle checkingt the set of cluster pods")
	return nil
}

// getPods returns the list of pods for a given namespace and instance
func getPods(r *ReconcileRqcluster, requestNamespace, instanceName string) (*corev1.PodList, error) {
	podList := &corev1.PodList{}
	err := r.client.List(context.TODO(), podList, client.InNamespace(requestNamespace), client.MatchingLabels{"cluster": instanceName})
	if err != nil {
		log.Error(err, "unable to find any pods that match this request")
		return podList, err
	}

	fmt.Printf("jeff list got back %d\n", len(podList.Items))
	return podList, nil
}

func createClusterPod(r *ReconcileRqcluster, instance *rqclusterv1alpha1.Rqcluster) error {
	// create the cluster pods
	// Define a new Pod object
	// get the Pod using the configmap, template, and CR
	mypod, err := newPodForCRFromTemplate(instance, r.client)
	if err != nil {
		fmt.Println(err.Error())
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

// verifyServices checks to see if there are 2 services for the
// rqcluster, one is a leader service and the other a cluster service,
// they are created if not found
func verifyServices(r *ReconcileRqcluster, instance *rqclusterv1alpha1.Rqcluster) error {

	// Check if the leader service already exists
	leaderService, err := newServiceForCRFromTemplate(instance, r.client)
	if err != nil {
		fmt.Println(err.Error())
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

	return nil
}
