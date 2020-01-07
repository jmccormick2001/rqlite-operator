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

	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

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
	PodTemplate *template.Template
}

// the rqlite pod template is found in the rqoperator ConfigMap
const containerTemplatePath = "/usr/local/bin/"
const ConfigMapName = "rq-config"
const PodTemplateFile = "pod-template.json"

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

	if cr.Spec.StorageClass != "" {
		fmt.Println("jeff here in sc check")
		if pod.Spec.Volumes[0].Name == "rqlite-storage" {
			fmt.Println("jeff here in sc check 2")
			vs := corev1.VolumeSource{}
			pod.Spec.Volumes[0].VolumeSource = vs
			pvc := corev1.PersistentVolumeClaimVolumeSource{}
			pvc.ClaimName = podName
			vs.PersistentVolumeClaim = &pvc
			pod.Spec.Volumes[0].VolumeSource = vs
			fmt.Printf("jeff here in sc check 3 %v\n", pod.Spec.Volumes[0])
		}
	}

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

	// set some sane defaults, we only work with storage classes
	myPVCInfo := PVCFields{
		ClaimName:        podName,
		AccessMode:       "ReadWriteOnce",
		Namespace:        cr.Namespace,
		StorageCapacity:  cr.Spec.StorageLimit,
		StorageClassName: cr.Spec.StorageClass,
	}
	return createPVC(myPVCInfo, cr, client)
}

// newServiceForCR returns a rqlite service with the same name/namespace as the cr
func newServiceForCRFromTemplate(leader bool, cr *rqclusterv1alpha1.Rqcluster, client client.Client) (*corev1.Service, error) {

	leaderStatus := ""
	serviceName := cr.Name
	if leader {
		leaderStatus = "true"
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

	svc.ObjectMeta.Namespace = cr.Namespace
	svc.ObjectMeta.Name = mySvcInfo.ServiceName

	svc.Spec.Ports = make([]corev1.ServicePort, 1)
	svc.Spec.Ports[0].Name = "rq"
	svc.Spec.Ports[0].Port = 4001
	svc.Spec.Ports[0].Protocol = corev1.ProtocolTCP
	svc.Spec.Ports[0].TargetPort = intstr.FromInt(4001)

	svc.Spec.Selector = make(map[string]string)
	svc.Spec.Selector["cluster"] = mySvcInfo.ClusterName
	svc.Spec.Selector["leader"] = mySvcInfo.LeaderStatus

	svc.Spec.Type = corev1.ServiceTypeClusterIP
	//svc.Spec.Type = corev1.ServiceTypeNodePort
	//svc.Spec.Type = corev1.ServiceTypeLoadBalancer
	return &svc, nil
}

func createPVC(myPVCInfo PVCFields, cr *rqclusterv1alpha1.Rqcluster, client client.Client) (*corev1.PersistentVolumeClaim, error) {
	pvc := corev1.PersistentVolumeClaim{}
	/**
	templates, err := getTemplates(cr.Namespace, client)
	if err != nil {
		return &pvc, err
	}

	var pvcBuffer bytes.Buffer
	templates.PVCTemplate.Execute(&pvcBuffer, myPVCInfo)

	fmt.Printf("%s\n", pvcBuffer.String())
	err = json.Unmarshal(pvcBuffer.Bytes(), &pvc)
	pvc.ObjectMeta.Namespace = cr.Namespace
	*/

	pvc.ObjectMeta.Name = myPVCInfo.ClaimName
	pvc.Spec.AccessModes = make([]corev1.PersistentVolumeAccessMode, 1)
	switch myPVCInfo.AccessMode {
	case string(corev1.ReadWriteOnce):
		pvc.Spec.AccessModes[0] = corev1.ReadWriteOnce
	case string(corev1.ReadOnlyMany):
		pvc.Spec.AccessModes[0] = corev1.ReadOnlyMany
	case string(corev1.ReadWriteMany):
		pvc.Spec.AccessModes[0] = corev1.ReadWriteMany
	default:
		return nil, fmt.Errorf("invalid AccessMode specified in CR")
	}

	pvc.Spec.Resources = corev1.ResourceRequirements{}
	pvc.Spec.Resources.Requests = make(corev1.ResourceList)
	q, err := resource.ParseQuantity(myPVCInfo.StorageCapacity)
	//q, err := resource.ParseQuantity("10Mi")
	if err != nil {
		return nil, err
	}
	pvc.Spec.Resources.Requests[corev1.ResourceStorage] = q

	pvc.Spec.StorageClassName = &myPVCInfo.StorageClassName
	pvc.ObjectMeta.Namespace = cr.Namespace

	fmt.Printf("jeff PVC to create is %v\n", pvc)
	d, err := yaml.Marshal(pvc)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("dump PVC %s\n", string(d))

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

	/**
	if configMapFound {
		value = cMap.Data[PVCTemplateFile]
	} else {
		templateData, err := ioutil.ReadFile(containerTemplatePath + PVCTemplateFile)
		if err != nil {
			return templates, err
		}
		value = string(templateData)
	}
	if value == "" {
		return templates, err
	}
	templates.PVCTemplate = template.Must(template.New("pvctemplate").Parse(value))
	if templates.PVCTemplate == nil {
		return templates, errors.New("PVCservice template didnt parse")
	}
	*/

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
