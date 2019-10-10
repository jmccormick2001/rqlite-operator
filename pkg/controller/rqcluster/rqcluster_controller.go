package rqcluster

import (
	"context"
	"fmt"

	rqclusterv1alpha1 "github.com/jmccormick2001/rq/pkg/apis/rqcluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_rqcluster")

// Add creates a new Rqcluster Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRqcluster{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("rqcluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Rqcluster
	err = c.Watch(&source.Kind{Type: &rqclusterv1alpha1.Rqcluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Rqcluster
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &rqclusterv1alpha1.Rqcluster{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileRqcluster implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileRqcluster{}

// ReconcileRqcluster reconciles a Rqcluster object
type ReconcileRqcluster struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Rqcluster object and makes changes based on the state read
// and what is in the Rqcluster.Spec
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileRqcluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Rqcluster")

	// Fetch the Rqcluster instance
	instance := &rqclusterv1alpha1.Rqcluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	err = rqReconcile(request, instance)

	// see if pods already exist for this rqcluster CR
	podList := &corev1.PodList{}
	err = r.client.List(context.TODO(), podList, client.InNamespace(request.Namespace), client.MatchingLabels{"cluster": instance.Name})
	if err != nil {
		reqLogger.Error(err, "unable to find any pods that match this request")
	} else {
		fmt.Printf("jeff list got back %d\n", len(podList.Items))
	}

	if len(podList.Items) == 0 {
		// create the cluster pods

		// Define a new Pod object
		// get the Pod using the configmap, template, and CR
		mypod, err := newPodForCRFromTemplate(instance, r.client)
		if mypod == nil {
			fmt.Println("mypod is nil")
		}
		if err != nil {
			fmt.Println(err.Error())
		}

		// Set Rqcluster instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, mypod, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		// Check if this Pod already exists
		found := &corev1.Pod{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: mypod.Name, Namespace: mypod.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Pod", "Pod.Namespace", mypod.Namespace, "Pod.Name", mypod.Name, "Namespace", mypod.ObjectMeta.Namespace)
			err = r.client.Create(context.TODO(), mypod)
			if err != nil {
				return reconcile.Result{}, err
			}

			// Pod created successfully - don't requeue
			return reconcile.Result{}, nil
		} else if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		// cluster Pods already exists
		reqLogger.Info("jeff reconcile: here is where we handle checkingt the set of cluster pods")
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}
