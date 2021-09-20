/*
Copyright 2020 The Laputa Cloud Co.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package component

import (
	"errors"

	"github.com/laputacloudco/sevendays-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GenerateTCPService creates a TCP Service for the SevenDays server.
func GenerateTCPService(sd v1alpha1.SevenDays) v1.Service {
	return v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: make(map[string]string),
			Labels:      standardLabels(sd),
			Name:        sd.Name + "-tcp",
			Namespace:   sd.Namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "26900tcp",
					Protocol:   v1.ProtocolTCP,
					Port:       26900,
					TargetPort: intstr.FromString("26900tcp"),
				},
			},
			Selector: standardLabels(sd),
			Type:     v1.ServiceTypeLoadBalancer,
		},
	}
}

// GenerateUDPService creates a UDP Service for the SevenDays server.
func GenerateUDPService(sd v1alpha1.SevenDays, lbip string) v1.Service {
	return v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: make(map[string]string),
			Labels:      standardLabels(sd),
			Name:        sd.Name + "-udp",
			Namespace:   sd.Namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "26900udp",
					Protocol:   v1.ProtocolUDP,
					Port:       26900,
					TargetPort: intstr.FromString("26900udp"),
				},
				{
					Name:       "26901udp",
					Protocol:   v1.ProtocolUDP,
					Port:       26901,
					TargetPort: intstr.FromString("26901udp"),
				},
				{
					Name:       "26902udp",
					Protocol:   v1.ProtocolUDP,
					Port:       26902,
					TargetPort: intstr.FromString("26902udp"),
				},
			},
			Selector:       standardLabels(sd),
			Type:           v1.ServiceTypeLoadBalancer,
			LoadBalancerIP: lbip,
		},
	}
}

// ExtractExistingLoadbalancerIP searches the passed slice of Services for
// existing loadbalancers and returns an IP address of those loadbalancers.
// If the Service slice is empty, returns empty string and no error.
// If there are Services of type LoadBalancer but they do not (yet?) have IP
// addresses assigned, return an error to indicate that this should be aborted
// and tried again later.
func ExtractExistingLoadbalancerIP(services []v1.Service) (string, error) {
	foundLB := false
	for _, svc := range services {
		if svc.Spec.Type != v1.ServiceTypeLoadBalancer {
			continue
		}
		foundLB = true
		if svc.Spec.LoadBalancerIP != "" {
			return svc.Spec.LoadBalancerIP, nil
		}
		if len(svc.Status.LoadBalancer.Ingress) > 0 && svc.Status.LoadBalancer.Ingress[0].IP != "" {
			return svc.Status.LoadBalancer.Ingress[0].IP, nil
		}
	}
	if foundLB {
		return "", errors.New("loadbalancers exist but do not yet have IPs")
	}
	return "", nil
}

// IndexService owner indexer func for controller-runtime.
func IndexService(o client.Object) []string {
	svc := o.(*v1.Service)
	owner := metav1.GetControllerOf(svc)
	if owner == nil {
		return nil
	}
	if owner.APIVersion != v1alpha1.GroupVersion.String() || owner.Kind != "SevenDays" {
		return nil
	}
	return []string{owner.Name}
}
