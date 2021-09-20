/*
Copyright 2020 The Laputa Cloud Co.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package component

import (
	"github.com/laputacloudco/sevendays-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NeedsUpdateDeployment returns if the passed Services are out of sync
func NeedsUpdateDeployment(want, have appsv1.Deployment) bool {
	return *want.Spec.Replicas != *have.Spec.Replicas ||
		want.Annotations["configmap-revision"] != have.Annotations["configmap-revision"]
}

// GenerateDeployment creates a Minecraft server deployment
func GenerateDeployment(sd v1alpha1.SevenDays) (appsv1.Deployment, error) {
	r := int32(0)
	if sd.Spec.Serve {
		r = 1
	}

	deploy := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: make(map[string]string),
			Labels:      standardLabels(sd),
			Name:        sd.Name,
			Namespace:   sd.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &r,
			Selector: &metav1.LabelSelector{
				MatchLabels: standardLabels(sd),
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: make(map[string]string),
					Labels:      standardLabels(sd),
					Name:        sd.Name,
				},
				Spec: corev1.PodSpec{
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "agentpool",
												Operator: corev1.NodeSelectorOpIn,
												Values: []string{
													"alpha",
												},
											},
										},
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "server",
							Image:           "docker.io/didstopia/7dtd-server",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "26900tcp",
									ContainerPort: 26900,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "26900udp",
									ContainerPort: 26900,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "26901udp",
									ContainerPort: 26901,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "26902udp",
									ContainerPort: 26902,
									Protocol:      corev1.ProtocolUDP,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "SEVEN_DAYS_TO_DIE_TELNET_PORT",
									Value: "8081",
								},
								{
									Name:  "SEVEN_DAYS_TO_DIE_TELNET_PASSWORD",
									Value: "password",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "app",
									MountPath: "/app/.local/share/7DaysToDie",
								},
								{
									Name:      "steamcmd",
									MountPath: "/steamcmd/7dtd",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "app",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: sd.Name + "-app",
								},
							},
						},
						{
							Name: "steamcmd",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: sd.Name + "-steamcmd",
								},
							},
						},
						{
							Name: "serverconfig",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: sd.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if sd.Spec.ServerConfigXML == "" {
		return deploy, nil
	}

	// if the serverconfig xml is set, we need to include an init container
	// to copy the contents to a writeable path
	deploy.Spec.Template.Spec.InitContainers = []corev1.Container{
		{
			Name:  "serverconfig",
			Image: "docker.io/busybox:latest",
			Command: []string{
				"sh",
				"-c",
				"cp /serverconfig.xml /app/.local/share/7DaysToDie/serverconfig.xml",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "app",
					MountPath: "/app/.local/share/7DaysToDie",
				},
				{
					Name:      "serverconfig",
					MountPath: "/serverconfig.xml",
					SubPath:   "serverconfig.xml",
				},
			},
		},
	}

	return deploy, nil
}

// IndexDeployment indexer func for controller-runtime
func IndexDeployment(o client.Object) []string {
	deploy := o.(*appsv1.Deployment)
	owner := metav1.GetControllerOf(deploy)
	if owner == nil {
		return nil
	}
	if owner.APIVersion != v1alpha1.GroupVersion.String() || owner.Kind != "SevenDays" {
		return nil
	}
	return []string{owner.Name}
}
