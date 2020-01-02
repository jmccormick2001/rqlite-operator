package rqcluster

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"text/template"

	rqclusterv1alpha1 "github.com/jmccormick2001/rq/pkg/apis/rqcluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// PVC template fields
type PVCFields struct {
	Namespace        string
	ClaimName        string
	AccessMode       string
	StorageCapacity  string
	StorageClassName string
}

// rqlite pod template fields
type PodFields struct {
	Namespace      string
	PodName        string
	ServiceAccount string
	ClusterName    string
	JoinAddress    string
}

// rqlite service template fields
type ServiceFields struct {
	Namespace    string
	ServiceName  string
	ClusterName  string
	LeaderStatus string
	PodNameMatch string
}

type ConfigMapTemplates struct {
	ServiceTemplate *template.Template
	PodTemplate     *template.Template
	PVCTemplate     *template.Template
}

// the rqlite pod template is found in the rqoperator ConfigMap
const containerTemplatePath = "/usr/local/bin/"
const ConfigMapName = "rq-config"
const PodTemplateFile = "pod-template.json"
const ServiceTemplateFile = "service-template.json"
const PVCTemplateFile = "pvc-template.json"

const ServiceAccountName = "default"
const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())

}

// newPodForCR returns a rqlite pod with the same name/namespace as the cr
func newPodForCRFromTemplate(joinAddress string, cr *rqclusterv1alpha1.Rqcluster, client client.Client) (*corev1.Pod, error) {

	pod := corev1.Pod{}

	podName := fmt.Sprintf("%s-%s", cr.Name, generateSuffix())

	myPodInfo := PodFields{
		PodName:        podName,
		Namespace:      cr.Namespace,
		ServiceAccount: ServiceAccountName,
		ClusterName:    cr.Name,
		JoinAddress:    joinAddress,
	}

	templates, err := getTemplates(cr.Namespace, client)
	if err != nil {
		return &pod, err
	}

	var podBuffer bytes.Buffer
	templates.PodTemplate.Execute(&podBuffer, myPodInfo)

	err = json.Unmarshal(podBuffer.Bytes(), &pod)
	pod.ObjectMeta.Namespace = cr.Namespace

	err = setResources(cr, client, &pod)
	if err != nil {
		return &pod, err
	}

	return &pod, nil
}

// newServiceForPod
func newServiceForPod(podName string, cr *rqclusterv1alpha1.Rqcluster, client client.Client) (*corev1.Service, error) {

	mySvcInfo := ServiceFields{
		ServiceName:  podName,
		Namespace:    cr.Namespace,
		ClusterName:  cr.Name,
		PodNameMatch: fmt.Sprintf(`"pod":"%s",`, podName),
		LeaderStatus: "",
	}

	return createService(mySvcInfo, cr, client)
}

// newPVCForPod
func newPVCForPod(podName string, cr *rqclusterv1alpha1.Rqcluster, client client.Client) (*corev1.PersistentVolumeClaim, error) {

	myPVCInfo := PVCFields{
		ClaimName:        podName,
		AccessMode:       "ReadWriteOnce",
		Namespace:        cr.Namespace,
		StorageCapacity:  "10Mi",
		StorageClassName: "my-local-storage",
	}

	return createPVC(myPVCInfo, cr, client)
}

// newServiceForCR returns a rqlite service with the same name/namespace as the cr
func newServiceForCRFromTemplate(leader bool, cr *rqclusterv1alpha1.Rqcluster, client client.Client) (*corev1.Service, error) {

	leaderStatus := ""
	serviceName := cr.Name
	if leader {
		leaderStatus = `"leader":"true",`
		serviceName = serviceName + "-leader"
	}

	mySvcInfo := ServiceFields{
		ServiceName:  serviceName,
		Namespace:    cr.Namespace,
		ClusterName:  cr.Name,
		PodNameMatch: "",
		LeaderStatus: leaderStatus,
	}

	return createService(mySvcInfo, cr, client)
}

func createService(mySvcInfo ServiceFields, cr *rqclusterv1alpha1.Rqcluster, client client.Client) (*corev1.Service, error) {

	svc := corev1.Service{}

	templates, err := getTemplates(cr.Namespace, client)
	if err != nil {
		return &svc, err
	}

	var svcBuffer bytes.Buffer
	templates.ServiceTemplate.Execute(&svcBuffer, mySvcInfo)

	err = json.Unmarshal(svcBuffer.Bytes(), &svc)
	svc.ObjectMeta.Namespace = cr.Namespace
	return &svc, nil
}

func createPVC(myPVCInfo PVCFields, cr *rqclusterv1alpha1.Rqcluster, client client.Client) (*corev1.PersistentVolumeClaim, error) {

	pvc := corev1.PersistentVolumeClaim{}

	templates, err := getTemplates(cr.Namespace, client)
	if err != nil {
		return &pvc, err
	}

	var pvcBuffer bytes.Buffer
	templates.PVCTemplate.Execute(&pvcBuffer, myPVCInfo)

	fmt.Printf("%s\n", pvcBuffer.String())
	err = json.Unmarshal(pvcBuffer.Bytes(), &pvc)
	pvc.ObjectMeta.Namespace = cr.Namespace
	return &pvc, nil
}

func getTemplates(namespace string, client client.Client) (ConfigMapTemplates, error) {
	templates := ConfigMapTemplates{}

	// lookup the rq configmap
	configMapFound := true
	cMap := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: ConfigMapName, Namespace: namespace}, cMap)
	if apierrors.IsNotFound(err) {
		configMapFound = false
	} else if err != nil {
		return templates, err
	}

	var value string
	if configMapFound {
		value = cMap.Data[PodTemplateFile]
	} else {
		templateData, err := ioutil.ReadFile(containerTemplatePath + PodTemplateFile)
		if err != nil {
			return templates, err
		}
		value = string(templateData)
	}

	if value == "" {
		return templates, err
	}
	templates.PodTemplate = template.Must(template.New("podtemplate").Parse(value))
	if templates.PodTemplate == nil {
		return templates, errors.New("pod template didnt parse")
	}

	if configMapFound {
		value = cMap.Data[ServiceTemplateFile]
	} else {
		templateData, err := ioutil.ReadFile(containerTemplatePath + ServiceTemplateFile)
		if err != nil {
			return templates, err
		}
		value = string(templateData)
	}
	if value == "" {
		return templates, err
	}
	templates.ServiceTemplate = template.Must(template.New("servicetemplate").Parse(value))
	if templates.ServiceTemplate == nil {
		return templates, errors.New("service template didnt parse")
	}

	if configMapFound {
		value = cMap.Data[PVCTemplateFile]
	} else {
		templateData, err := ioutil.ReadFile(containerTemplatePath + PVCTemplateFile)
		if err != nil {
			return templates, err
		}
		value = string(templateData)
	}
	templates.PVCTemplate = template.Must(template.New("pvctemplate").Parse(value))
	if templates.PVCTemplate == nil {
		return templates, errors.New("PVC template didnt parse")
	}

	return templates, nil
}

// generate a 4 char random string
func generateSuffix() string {
	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
