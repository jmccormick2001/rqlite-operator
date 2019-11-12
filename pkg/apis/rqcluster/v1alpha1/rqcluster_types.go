package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RqPodAffinity defines the pod affinity for an rq cluster
// with the operator spreading rqcluster nodes across this set of affinity
// values
// +k8s:openapi-gen=true
type RqPodAffinity struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// RqclusterSpec defines the desired state of Rqcluster
// +k8s:openapi-gen=true
type RqclusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size          int32         `json:"size"`
	CpuLimit      string        `json:"cpulimit"`
	CpuRequest    string        `json:"cpurequest"`
	MemoryLimit   string        `json:"memorylimit"`
	MemoryRequest string        `json:"memoryrequest"`
	StorageClass  string        `json:"storageclass"`
	StorageLimit  string        `json:"storagelimit"`
	PodAffinity   RqPodAffinity `json:"rqpodaffinity"`
}

// RqclusterStatus defines the observed state of Rqcluster
// +k8s:openapi-gen=true
type RqclusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Rqcluster is the Schema for the rqclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=rqclusters,scope=Namespaced
type Rqcluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RqclusterSpec   `json:"spec,omitempty"`
	Status RqclusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RqclusterList contains a list of Rqcluster
type RqclusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rqcluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Rqcluster{}, &RqclusterList{})
}
