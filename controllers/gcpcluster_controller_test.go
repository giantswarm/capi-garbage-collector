package controllers_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/giantswarm/capi-garbage-collector/controllers"
	"github.com/giantswarm/capi-garbage-collector/controllers/controllersfakes"
)

var _ = Describe("GCPCluster Reconciliation", func() {
	const (
		clusterName = "vegeta"
		gcpProject  = "multiverse"
		finalizer   = "final.count/down"

		timeout  = time.Second * 5
		interval = time.Millisecond * 250
	)

	var (
		ctx        context.Context
		gcpCluster *capg.GCPCluster

		fakeComputeClient *controllersfakes.FakeComputeClient
		clusterReconciler *controllers.GCPClusterReconciler

		result       reconcile.Result
		reconcileErr error
	)

	BeforeEach(func() {
		SetDefaultConsistentlyDuration(timeout)
		SetDefaultConsistentlyPollingInterval(interval)
		SetDefaultEventuallyPollingInterval(interval)
		SetDefaultEventuallyTimeout(timeout)
		ctx = context.Background()

		gcpCluster = &capg.GCPCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterName,
				Namespace: namespace,
			},
			Spec: capg.GCPClusterSpec{
				Project: gcpProject,
			},
		}

		Expect(k8sClient.Create(ctx, gcpCluster)).To(Succeed())

		fakeComputeClient = new(controllersfakes.FakeComputeClient)
		fakeComputeClient.DeleteRoutesReturns(nil)

		clusterReconciler = controllers.NewGCPClusterReconciler(k8sClient, fakeComputeClient)
	})

	JustBeforeEach(func() {
		req := reconcile.Request{
			NamespacedName: k8stypes.NamespacedName{
				Name:      gcpCluster.Name,
				Namespace: namespace,
			},
		}

		result, reconcileErr = clusterReconciler.Reconcile(ctx, req)
	})

	When("the cluster is not marked for deletion", func() {
		It("should not delete routes", func() {
			Expect(reconcileErr).NotTo(HaveOccurred())
			Expect(result.Requeue).To(BeFalse())

			Expect(fakeComputeClient.DeleteRoutesCallCount()).To(Equal(0))
		})
	})

	When("the cluster is marked for deletion", func() {
		BeforeEach(func() {
			originalCluster := gcpCluster.DeepCopy()
			controllerutil.AddFinalizer(gcpCluster, finalizer)
			err := k8sClient.Patch(ctx, gcpCluster, client.MergeFrom(originalCluster))
			Expect(err).NotTo(HaveOccurred())

			Expect(k8sClient.Delete(ctx, gcpCluster)).To(Succeed())
		})

		It("Reconciles Sucessfully", func() {
			Expect(reconcileErr).NotTo(HaveOccurred())
			Expect(result.Requeue).To(BeFalse())
		})

		It("Deletes routes", func() {
			Expect(fakeComputeClient.DeleteRoutesCallCount()).To(Equal(1))
		})
	})

	When("the client fails to delete", func() {
		BeforeEach(func() {
			originalCluster := gcpCluster.DeepCopy()
			controllerutil.AddFinalizer(gcpCluster, finalizer)
			err := k8sClient.Patch(ctx, gcpCluster, client.MergeFrom(originalCluster))
			Expect(err).NotTo(HaveOccurred())

			Expect(k8sClient.Delete(ctx, gcpCluster)).To(Succeed())
			fakeComputeClient.DeleteRoutesReturns(errors.New("Ki Blast!"))
		})

		It("returns an error", func() {
			Expect(reconcileErr).To(MatchError(ContainSubstring("Ki Blast!")))
		})
	})

	When("the cluster is not found", func() {
		BeforeEach(func() {
			Expect(k8sClient.Delete(ctx, gcpCluster)).To(Succeed())
		})

		It("should reconcile without any error", func() {
			Expect(reconcileErr).NotTo(HaveOccurred())
			Expect(result.Requeue).To(BeFalse())

			namespacedName := k8stypes.NamespacedName{
				Name:      gcpCluster.Name,
				Namespace: namespace,
			}
			cluster := &capg.GCPCluster{}
			err := k8sClient.Get(ctx, namespacedName, cluster)

			Expect(err).To(HaveOccurred())
			Expect(k8serrors.IsNotFound(err)).To(BeTrue())
		})

	})
})
