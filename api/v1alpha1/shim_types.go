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
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ShimSpec defines the desired state of Shim
type ShimSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// NodeSelector is a selector which should be true for nodes you want this Shim to run on.
	// This should match a node's labels when you want the shim to run on that node. If this is
	// empty, the shim will apply to all nodes.
	// +required
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Source is the source of the shim
	// +required
	Source ShimSource `json:"source"`

	// RuntimeClass is the name of the RuntimeClass to use for the shim
	// +required
	RuntimeClass string `json:"runtimeClass,omitempty"`

	// RolloutStrategy is the strategy to use when rolling out the shim
	// +required
	RolloutStrategy ShimRolloutStrategy `json:"rolloutStrategy,omitempty"`

	// TODO: how do upgrades work? we should have some way of upgrading the shim??
}

// Only one of its members may be specified
type ShimSource struct {
	// AnonymousHttp is a shim hosted at a URL which can be downloaded anonymously
	AnonymousHttp *AnonymousHttpSource `json:"anonymousHttp,omitempty"`
}

type AnonymousHttpSource struct {
	// Location is the URL of the shim's .tar.gz file
	// +required
	Location string `json:"location"`
	// File is the name of the file inside the .tar.gz which contains the shim
	// +required
	File string `json:"file"`
}

// Only one of its members may be specified
type ShimRolloutStrategy struct {
	RollingRolloutStrategy *RollingRolloutStrategy `json:"rolling,omitempty"`
}

// Only one of its members may be specified
type RollingRolloutStrategy struct {
	// MaxUnavailable is the maximum number of nodes which may be unavailable at any given time
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// +required
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable"`
}

// ShimStatus defines the observed state of Shim
type ShimStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The generation observed by the controller
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions is an array of current observed conditions
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions"`

	// TODO: need to figure out a more deterministic way of doing this
	// TODO: should we do this or is a label on the node enough?

	// UpgradedNodes is a list of nodes which have been upgraded to the latest version of the shim
	UpgradedNodes []NodeReference `json:"upgradedNodes,omitempty"`
	// UpgradingNodes is a list of nodes which are currently being upgraded to the latest version of the shim
	UpgradingNodes []NodeReference `json:"upgradingNodes,omitempty"`
	// QueuedNodes is a list of nodes which are queued to be upgraded to the latest version of the shim
	QueuedNodes []NodeReference `json:"queuedNodes,omitempty"`
}

type NodeReference struct {
	// Name is the name of the node
	// +required
	Name string `json:"name"`
	// UID is the UID of the node
	// +required
	UID string `json:"uid"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Shim is the Schema for the shims API
type Shim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ShimSpec   `json:"spec,omitempty"`
	Status ShimStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ShimList contains a list of Shim
type ShimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Shim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Shim{}, &ShimList{})
}
