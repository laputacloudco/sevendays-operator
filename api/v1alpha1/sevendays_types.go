/*
Copyright 2020 The Laputa Cloud Co.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SevenDaysSpec defines the desired state of SevenDays
type SevenDaysSpec struct {
	// ServerConfigXML is a 7 Days config file literal.
	ServerConfigXML string `json:"serverconfig.xml,omitempty"`
	// Serve tells the controller to run or stop this server.
	Serve bool `json:"serve,omitempty"`
}

// ServerStatus indicates the Server Status
//+kubebuilder:validation:Enum=Creating;Destroying;Running;Starting;Stopped;Stopping;Unknown;Updating
type ServerStatus string

const (
	// Creating status
	Creating ServerStatus = "Creating"
	// Destroying status
	Destroying ServerStatus = "Destroying"
	// Running status
	Running ServerStatus = "Running"
	// Starting status
	Starting ServerStatus = "Starting"
	// Stopped status
	Stopped ServerStatus = "Stopped"
	// Stopping status
	Stopping ServerStatus = "Stopping"
	// Unknown status
	Unknown ServerStatus = "Unknown"
	// Updating status
	Updating ServerStatus = "Updating"
)

// SevenDaysStatus defines the observed state of SevenDays
type SevenDaysStatus struct {
	// Status indicates the Server Status
	Status ServerStatus `json:"status,omitempty"`
	// Address the public server address
	Address string `json:"address,omitempty"`
	// Cost is the running cost the servec
	Cost float64 `json:"cost,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
//+kubebuilder:printcolumn:name="Address",type=string,JSONPath=`.status.address

// SevenDays is the Schema for the sevendays API
type SevenDays struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SevenDaysSpec   `json:"spec,omitempty"`
	Status SevenDaysStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SevenDaysList contains a list of SevenDays
type SevenDaysList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SevenDays `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SevenDays{}, &SevenDaysList{})
}
