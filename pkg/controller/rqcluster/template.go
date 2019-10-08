package rqcluster

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"

	rqclusterv1alpha1 "github.com/jmccormick2001/rq/pkg/apis/rqcluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"text/template"
)

// rqlite pod template fields
type PodFields struct {
	Namespace      string
	PodName        string
	ServiceAccount string
}

// the rqlite pod template is found in the rqoperator ConfigMap
const ConfigMapName = "rq-configs"
const TemplateRoot = "/rq-configs"
const PodTemplateFile = "pod-template.yaml"

// newPodForCR returns a rqlite pod with the same name/namespace as the cr
func newPodForCRFromTemplate(cr *rqclusterv1alpha1.Rqcluster) (*corev1.Pod, error) {

	var pod *corev1.Pod

	myPodInfo := PodFields{
		PodName:        "rqpod1",
		Namespace:      "default",
		ServiceAccount: "default",
	}

	podBuffer, err := getPodTemplate(myPodInfo)
	if err != nil {
		return pod, err
	}

	err = yaml.Unmarshal(podBuffer.Bytes(), pod)
	/**
	labels := map[string]string{
		"app": cr.Name,
	}
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
	*/
	return pod, nil
}

func getPodTemplate(myPodInfo PodFields) (bytes.Buffer, error) {
	var podBuffer bytes.Buffer
	var client *kubernetes.Clientset
	cMap, err := getConfigMap(client, myPodInfo.Namespace, ConfigMapName)
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

func getConfigMap(client *kubernetes.Clientset, namespace, name string) (*corev1.ConfigMap, error) {
	cfg, err := client.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	return cfg, err

}
