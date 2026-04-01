package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	goproexamplecomv1alpha1 "basic/projects/kubernetes/api/v1alpha1"
)

// GoProReconciler reconciles a GoPro object
type GoProReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop
func (r *GoProReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Log the reconciliation request
	logger.Info("Reconciling GoPro", "request", req.NamespacedName)

	// Return early - this is a scaffold
	// Full implementation would fetch and reconcile the GoPro resource
	logger.Info("GoPro reconciliation complete", "request", req.NamespacedName)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *GoProReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&goproexamplecomv1alpha1.GoPro{}).
		Complete(r)
}

// InitSchemes initializes the schemes for the operator
func InitSchemes(s *runtime.Scheme) error {
	if err := goproexamplecomv1alpha1.AddToScheme(s); err != nil {
		return fmt.Errorf("failed to add gopro scheme: %w", err)
	}
	return nil
}
