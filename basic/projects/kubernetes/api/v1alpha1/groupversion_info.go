// Package v1alpha1 contains the API for the GoPro custom resource
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is the group version used in this package
	GroupVersion = schema.GroupVersion{Group: "gopro.example.com", Version: "v1alpha1"}

	// SchemeBuilder is used to add the scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme is used to add the scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// Kind returns the GroupKind for GoPro
func Kind(kind string) schema.GroupKind {
	return GroupVersion.WithKind(kind).GroupKind()
}

// Resource returns the GroupVersionResource for GoPro
func Resource(resource string) schema.GroupVersionResource {
	return GroupVersion.WithResource(resource)
}

// GetScheme returns the scheme
func GetScheme() *runtime.Scheme {
	sc := runtime.NewScheme()
	if err := AddToScheme(sc); err != nil {
		panic(err)
	}
	return sc
}
