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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var sc = "default"

// GeneratePVC creates PVCs for the SevenDays server
func GeneratePVCs(sd v1alpha1.SevenDays) ([]v1.PersistentVolumeClaim, error) {
	pvcSizeRequest, err := resource.ParseQuantity("16Gi")
	if err != nil {
		return []v1.PersistentVolumeClaim{}, err
	}
	return []v1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: make(map[string]string),
				Labels:      standardLabels(sd),
				Name:        sd.Name + "-app",
				Namespace:   sd.Namespace,
			},
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				StorageClassName: &sc,
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: pvcSizeRequest,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: make(map[string]string),
				Labels:      standardLabels(sd),
				Name:        sd.Name + "-steamcmd",
				Namespace:   sd.Namespace,
			},
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				StorageClassName: &sc,
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: pvcSizeRequest,
					},
				},
			},
		},
	}, nil
}

// IndexPVC indexer func for controller-runtime
func IndexPVC(o client.Object) []string {
	pvc := o.(*v1.PersistentVolumeClaim)
	owner := metav1.GetControllerOf(pvc)
	if owner == nil {
		return nil
	}
	if owner.APIVersion != v1alpha1.GroupVersion.String() || owner.Kind != "SevenDays" {
		return nil
	}
	return []string{owner.Name}
}
