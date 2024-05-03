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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Replica struct {
	Count   int32  `json:"count"` // count of pods to be deployed
	Version string `json:"version"`
}

type Container struct {
	Image string `json:"image"`
	Port  int    `json:"port"`
}

// DeploySetSpec defines the desired state of DeploySet
type DeploySetSpec struct {
	Replica   Replica   `json:"replica"`
	Container Container `json:"container"`
}

// DeploySetStatus defines the observed state of DeploySet
type DeploySetStatus struct {
	ReadyReplicas     int                `json:"readyReplicas"`
	AvailableReplicas int                `json:"availableReplicas"`
	Condition         []metav1.Condition `json:"condition"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DeploySet is the Schema for the deploysets API
type DeploySet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploySetSpec   `json:"spec,omitempty"`
	Status DeploySetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeploySetList contains a list of DeploySet
type DeploySetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeploySet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeploySet{}, &DeploySetList{})
}
