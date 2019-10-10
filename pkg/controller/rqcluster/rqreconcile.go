package rqcluster

import (
	//	"context"
	//	"errors"
	"fmt"
	//	"k8s.io/apimachinery/pkg/types"
	//	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	rqclusterv1alpha1 "github.com/jmccormick2001/rq/pkg/apis/rqcluster/v1alpha1"
	//	corev1 "k8s.io/api/core/v1"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// rqReconcile implements the Reconcile for the rq-operator
func rqReconcile(request reconcile.Request, instance *rqclusterv1alpha1.Rqcluster) error {

	var err error
	fmt.Printf("jeff in rqReconcile\n")
	return err
}
