/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DbVolume string

const DbVolumeBig DbVolume = "BIG"
const DbVolumeSmall DbVolume = "SMALL"
const DbVolumeMedium DbVolume = "MEDIUM"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProvisionRequestSpec defines the desired state of ProvisionRequest
type ProvisionRequestSpec struct {
	// +optional
	IngressEntrance string `json:"ingressEntrance" `
	// +optional
	BusinessDbVolume DbVolume `json:"businessDbVolume"`
	//+kubebuilder:validation:MinLength=1
	NamespaceName string `json:"namespaceName"`
}

// ProvisionRequestStatus defines the observed state of ProvisionRequest
type ProvisionRequestStatus struct {
	metav1.TypeMeta `json:",inline" `
	IngressReady    bool `json:"ingressReady" `
	DbReady         bool `json:"dbReady"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ProvisionRequest is the Schema for the provisionrequests API
type ProvisionRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProvisionRequestSpec   `json:"spec,omitempty"`
	Status ProvisionRequestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProvisionRequestList contains a list of ProvisionRequest
type ProvisionRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProvisionRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProvisionRequest{}, &ProvisionRequestList{})
}
