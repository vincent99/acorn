package v1

import (
	v1 "github.com/acorn-io/runtime/pkg/apis/internal.admin.acorn.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ComputeClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Memory           v1.ComputeClassMemory `json:"memory,omitempty"`
	Description      string                `json:"description,omitempty"`
	Default          bool                  `json:"default"`
	SupportedRegions []string              `json:"supportedRegions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ComputeClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComputeClass `json:"items"`
}
