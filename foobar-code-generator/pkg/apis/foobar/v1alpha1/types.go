package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooBar is a specification for a FooBar resource
type FooBar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FooBarSpec `json:"spec"`
}

// FooBarSpec is the spec for a FooBar resource
type FooBarSpec struct {
	Foo     string `json:"foo"`
	Bar     bool   `json:"bar"`
	Command string `json:"command"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooBarList is a list of FooBar resources
type FooBarList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FooBar `json:"items"`
}
