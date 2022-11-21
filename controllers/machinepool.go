package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	capi "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/giantswarm/capi-garbage-collector/pkg/key"
)

var (
	CAPIMachinePoolFinalizer = "machinepool.cluster.x-k8s.io"
)

type GarbageCollectorController struct {
	client ctrlclient.Client
}

func NewGarbageCollectorController(client ctrlclient.Client) *GarbageCollectorController {
	return &GarbageCollectorController{
		client: client,
	}
}

func (r *GarbageCollectorController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&capi.MachinePool{}).
		Complete(r)
}

func (r *GarbageCollectorController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling")
	defer logger.Info("Done reconciling")

	var machinePool capi.MachinePool

	err := r.client.Get(ctx, ctrlclient.ObjectKey{Name: req.Name, Namespace: req.Namespace}, &machinePool)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			logger.Info("MachinePool no longer exists")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, microerror.Mask(err)
	}

	logger = logger.WithValues("machine pool", fmt.Sprintf("%s/%s", machinePool.Namespace, machinePool.Name))

	if !machinePool.DeletionTimestamp.IsZero() {
		logger.Info("Reconciling delete - cleanup")
		return r.reconcileDelete(ctx, &machinePool, logger)
	}

	return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 5}, nil
}

func (r *GarbageCollectorController) reconcileDelete(ctx context.Context, machinePool *capi.MachinePool, logger logr.Logger) (ctrl.Result, error) {
	clusterName := key.GetClusterIDFromLabels(machinePool.ObjectMeta)
	var kubeconfigSecret corev1.Secret
	err := r.client.Get(ctx, ctrlclient.ObjectKey{Name: key.KubeconfigSecretName(clusterName), Namespace: machinePool.Namespace}, &kubeconfigSecret)
	if k8sErrors.IsNotFound(err) {
		logger.Info("kubeconfig for the cluster no longer exists, cleaning machine pool")

		if len(machinePool.Finalizers) == 1 && machinePool.Finalizers[0] == CAPIMachinePoolFinalizer {
			machinePool.Finalizers = []string{}
			err = r.client.Update(ctx, machinePool)
			if err != nil {
				logger.Error(err, "failed to remove finalizer from machine pool")
				return ctrl.Result{}, err
			}
			logger.Info("cleanup up MachinePool")
		} else {
			logger.Info("object still has multiple finalizers, skipping")
		}
	} else if err != nil {
		logger.Error(err, "failed to get kubeconfig secret")
		return ctrl.Result{}, err
	} else {
		logger.Info("kubeconfig for the cluster still exists, skipping")
	}

	return ctrl.Result{Requeue: true, RequeueAfter: time.Minute * 3}, nil
}
