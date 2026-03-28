// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium
package v2alpha1

import (
	slimv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CiliumQoSPolicy defines QOS rules which should be
// applied to a pods egress traffic
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:categories={cilium},singular="ciliumqospolicy",path="ciliumqospolicies",scope="Namespaced",shortName={cqosp,ciliumqosp}
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.matchedEndpoints",name="Endpoints",type=integer
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",name="Age",type=date
// +kubebuilder:storageversion
type CiliumQoSPolicy struct {
	// +deepequal-gen=false
	metav1.TypeMeta `json:",inline"`
	// +deepequal-gen=false
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec   CiliumQoSPolicySpec   `json:"spec"`
	Status CiliumQoSPolicyStatus `json:"status,omitempty"`
}

type CiliumQoSPolicyStatus struct {
	// MatchedEndpoints is the number of local endpoints currently programmed
	// with this policy's BPF rules.
	// +kubebuilder:validation:Optional
	MatchedEndpoints int32 `json:"matchedEndpoints,omitempty"`
}

type CiliumQoSPolicySpec struct {
	// Selector to decide on which pods the
	// policy is going to be applied. To furhter
	// restirct the policy use the `to` field.
	// +kubebuilder:validation:Required
	PodSelector slimv1.LabelSelector `json:"podSelector"`

	// Restricts the policy to specific IP and Port
	// based destinations. So one pod could create
	// different dscp marks for different destinations.
	// +kubebuilder:validation:Optional
	To CiliumQoSPolicyTo `json:"to"`

	// A DSCP QoS handler which defines the fields which
	// are required for dscp.
	// +kubebuilder:validation:Optional
	DSCP CiliumQoSDSCP `json:"dscp"`
}

type CiliumQoSPolicyTo struct {
	// Supports both IPv4 and IPv6 CIDRs. If omitted, all
	// world-bound traffic from matching pods is marked. With
	// the value of the corresponding QoS method.
	// +kubebuilder:validation:Optional
	// +listType=set
	DestinationCIDRs []QoSToCIDR `json:"destinationCIDRs,omitempty"`

	// DestinationPorts restricts marking to specific L4 destination ports.
	// Applies to both TCP and UDP. If empty, any port matches.
	// Maximum 16 ports per rule.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=16
	// +listType=set
	DestinationPorts []QoSToPort `json:"destinationPorts,omitempty"`
}

// QoSToCIDR is an IP CIDR.
//
// +kubebuilder:validation:Format=cidr
type QoSToCIDR string

// +kubebuilder:validation:Minimum=1
// +kubebuilder:validation:Maximum=65535
type QoSToPort uint16

type CiliumQoSDSCP struct {
	// Value is the 6-bit DSCP codepoint (0-63).
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=63
	// +kubebuilder:validation:Required
	Value uint8 `json:"value"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +deepequal-gen=false

// CiliumQoSPolicyList is a list of
// CiliumQoSPolicy objects.
type CiliumQoSPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items is a list of CiliumQoSPolicy.
	Items []CiliumQoSPolicy `json:"items"`
}
