/*
Copyright 2020 The Laputa Cloud Co.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package component

import (
	"github.com/laputacloudco/sevendays-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NeedsUpdateConfigMap returns if the passed Services are out of sync
func NeedsUpdateConfigMap(want, have v1.ConfigMap) bool {
	if len(want.Data) != len(have.Data) {
		return true
	}
	for k, v := range want.Data {
		if vv, ok := have.Data[k]; !ok {
			return true
		} else if v != vv {
			return true
		}
	}
	return false
}

// GenerateConfigMap creates a configmap from a Minecraft server config struct
func GenerateConfigMap(sd v1alpha1.SevenDays) v1.ConfigMap {
	cm := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: make(map[string]string),
			Labels:      standardLabels(sd),
			Name:        sd.Name,
			Namespace:   sd.Namespace,
		},
		Data: map[string]string{
			"serverconfig.xml": sd.Spec.ServerConfigXML,
		},
	}
	return cm
}

// IndexConfigMap indexer func for controller-runtime
func IndexConfigMap(o client.Object) []string {
	cm := o.(*v1.ConfigMap)
	owner := metav1.GetControllerOf(cm)
	if owner == nil {
		return nil
	}
	if owner.APIVersion != v1alpha1.GroupVersion.String() || owner.Kind != "SevenDays" {
		return nil
	}
	return []string{owner.Name}
}
