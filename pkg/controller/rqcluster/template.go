package rqcluster

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"math/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	rqclusterv1alpha1 "github.com/jmccormick2001/rq/pkg/apis/rqcluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"text/template"
)

// rqlite pod template fields
type PodFields struct {
	Namespace      string
	PodName        string
	ServiceAccount string
	ClusterName    string
}

// the rqlite pod template is found in the rqoperator ConfigMap
const ConfigMapName = "rq-config"
const PodTemplateFile = "pod-template.json"

const ServiceAccountName = "default"
const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())

}

// newPodForCR returns a rqlite pod with the same name/namespace as the cr
func newPodForCRFromTemplate(cr *rqclusterv1alpha1.Rqcluster, client client.Client) (*corev1.Pod, error) {

	pod := corev1.Pod{}

	podName := fmt.Sprintf("%s-%s", cr.Name, generateSuffix())

	myPodInfo := PodFields{
		PodName:        podName,
		Namespace:      cr.Namespace,
		ServiceAccount: ServiceAccountName,
		ClusterName:    cr.Name,
	}

	podBuffer, err := getPodTemplate(myPodInfo, cr.Namespace, client)
	if err != nil {
		return &pod, err
	}

	fmt.Println("podBuffer is %s\n", podBuffer.String())
	err = json.Unmarshal(podBuffer.Bytes(), &pod)
	pod.ObjectMeta.Namespace = cr.Namespace
	return &pod, nil
}

func getPodTemplate(myPodInfo PodFields, namespace string, client client.Client) (bytes.Buffer, error) {
	var podBuffer bytes.Buffer

	// lookup the rq configmap
	cMap := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: ConfigMapName, Namespace: namespace}, cMap)
	if err != nil {
		return podBuffer, err
	}

	value := cMap.Data[PodTemplateFile]
	if value == "" {
		fmt.Println("pod template value is empty")
		return podBuffer, err
	}
	fmt.Println(value)
	var tmpl *template.Template
	tmpl = template.Must(template.New("podtemplate").Parse(value))
	if tmpl == nil {
		fmt.Println("error in template")
		return podBuffer, errors.New("template didnt parse")
	}

	tmpl.Execute(os.Stdout, myPodInfo)
	tmpl.Execute(&podBuffer, myPodInfo)

	return podBuffer, nil
}

// generate a 4 char random string
func generateSuffix() string {
	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
