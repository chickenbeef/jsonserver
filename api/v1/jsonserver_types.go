/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JsonServerSpec defines the desired state of JsonServer.
type JsonServerSpec struct {
	// Replicas is the number of instances of the JsonServer to run
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas"`

	// JsonConfig is the JSON configuration to be served by the JsonServer
	JsonConfig string `json:"jsonConfig"`
}

// JsonServerStatus defines the observed state of JsonServer.
type JsonServerStatus struct {
	// +kubebuilder:validation:Enum=Synced;Error
	State string `json:"state,omitempty"`

	// Message provides additional information about the JsonServer state
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JsonServer is the Schema for the jsonservers API.
type JsonServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JsonServerSpec   `json:"spec,omitempty"`
	Status JsonServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="Number of replicas"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.state",description="Current status"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.message",description="Status message"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// JsonServerList contains a list of JsonServer.
type JsonServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JsonServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JsonServer{}, &JsonServerList{})
}
