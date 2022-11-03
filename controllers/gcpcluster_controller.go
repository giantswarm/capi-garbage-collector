/*
Copyright 2022.

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
)

//counterfeiter:generate . ComputeClient
type ComputeClient interface {
	DeleteRoutes(ctx context.Context, cluster *capg.GCPCluster) error
}

// GCPClusterReconciler reconciles a GCPCluster object
type GCPClusterReconciler struct {
	client.Client

	computeClient ComputeClient
}

func NewGCPClusterReconciler(client client.Client, computeClient ComputeClient) *GCPClusterReconciler {
	return &GCPClusterReconciler{
		Client:        client,
		computeClient: computeClient,
	}
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io.giantswarm.io,resources=gcpclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io.giantswarm.io,resources=gcpclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io.giantswarm.io,resources=gcpclusters/finalizers,verbs=update
func (r *GCPClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling dangling GCP Cluster resources")

	defer logger.Info("Finished cleaning up dangling resources")

	cluster := &capg.GCPCluster{}
	err := r.Get(ctx, req.NamespacedName, cluster)
	if err != nil {
		logger.Error(err, "could not get the gcp cluster")
		return reconcile.Result{}, nil
	}

	if !cluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, cluster)
	}

	return ctrl.Result{}, nil
}

func (r *GCPClusterReconciler) reconcileDelete(ctx context.Context, cluster *capg.GCPCluster) (reconcile.Result, error) {
	logger := r.getLogger(ctx)

	err := r.computeClient.DeleteRoutes(ctx, cluster)
	if err != nil {
		logger.Error(err, "failed to delete routes")
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *GCPClusterReconciler) getLogger(ctx context.Context) logr.Logger {
	logger := log.FromContext(ctx)
	return logger.WithName("gcpcluster-reconciler")
}

// SetupWithManager sets up the controller with the Manager.
func (r *GCPClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&capg.GCPCluster{}).
		Complete(r)
}
