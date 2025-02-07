/*
Copyright 2024.

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

// KDaemonSpec defines the desired state of KDaemon
type KDaemonSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Pod image
	Image string `json:"image,omitempty"`
}

// KDaemonStatus defines the observed state of KDaemon
type KDaemonStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// available replicas number
	AvailableReplicas int `json:"availableReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KDaemon is the Schema for the kdaemons API
type KDaemon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KDaemonSpec   `json:"spec,omitempty"`
	Status KDaemonStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KDaemonList contains a list of KDaemon
type KDaemonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDaemon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDaemon{}, &KDaemonList{})
}
