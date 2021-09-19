/*
Copyright 2020 The Laputa Cloud Co.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gamev1alpha1 "github.com/laputacloudco/k8s-7dtd/api/v1alpha1"
)

// SevenDaysReconciler reconciles a SevenDays object
type SevenDaysReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=game.laputacloud.co,resources=sevendays,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=game.laputacloud.co,resources=sevendays/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=game.laputacloud.co,resources=sevendays/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=create;delete;get;list;update;watch
//+kubebuilder:rbac:groups="",resources=configmaps;persistentvolumeclaims;services,verbs=create;delete;get;list;update;watch

// Reconcile a SevenDays.
func (r *SevenDaysReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SevenDaysReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gamev1alpha1.SevenDays{}).
		Complete(r)
}
