package rqcluster

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	rqclusterv1alpha1 "github.com/jmccormick2001/rq/pkg/apis/rqcluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"text/template"
)

// rqlite pod template fields
type PodFields struct {
	PodName        string
	ServiceAccount string
}

// the rqlite pod template is found in the rqoperator ConfigMap
const PodTemplateFile = "pod-template.yaml"

// newPodForCR returns a rqlite pod with the same name/namespace as the cr
func newPodForCRFromTemplate(cr *rqclusterv1alpha1.Rqcluster) (*corev1.Pod, error) {

	var pod *corev1.Pod

	myPodInfo := PodFields{
		PodName: "rqpod1",
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
	buf, err := ioutil.ReadFile(PodTemplateFile)
	if err != nil {
		fmt.Println(err.Error())
		return podBuffer, err
	}
	value := string(buf)
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
