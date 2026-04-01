// Package v1alpha1 contains the API for the GoPro custom resource
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// GoProSpec defines the desired state of GoPro
type GoProSpec struct {
	// Image is the container image to deploy
	Image string `json:"image,omitempty"`

	// Replicas is the number of pods to deploy
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Env is a map of environment variables
	// +optional
	Env map[string]string `json:"env,omitempty"`

	// Port is the container port
	// +optional
	Port int32 `json:"port,omitempty"`
}

// GoProStatus defines the observed state of GoPro
type GoProStatus struct {
	// AvailableReplicas is the number of running pods
	// +optional
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`

	// ReadyReplicas is the number of ready pods
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
//+kubebuilder:printcolumn:name="Available",type="integer",JSONPath=".status.availableReplicas"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// GoPro is the Schema for the goproes API
type GoPro struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoProSpec   `json:"spec,omitempty"`
	Status GoProStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GoProList contains a list of GoPro resources
type GoProList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GoPro `json:"items"`
}

// DeepCopyInto copies the receiver into out
func (in *GoPro) DeepCopyInto(out *GoPro) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy returns a deep copy of the GoPro
func (in *GoPro) DeepCopy() *GoPro {
	if in == nil {
		return nil
	}
	out := new(GoPro)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is needed for runtime.Object interface
func (in *GoPro) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto copies the receiver into out
func (in *GoProList) DeepCopyInto(out *GoProList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]GoPro, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

// DeepCopy returns a deep copy of the GoProList
func (in *GoProList) DeepCopy() *GoProList {
	if in == nil {
		return nil
	}
	out := new(GoProList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is needed for runtime.Object interface
func (in *GoProList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
