package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	goproexamplecomv1alpha1 "basic/projects/kubernetes/api/v1alpha1"
	"basic/projects/kubernetes/controllers"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	klog.InitFlags(nil)
}

func main() {
	// Parse command line flags
	flag.Parse()

	// Set up logging
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Set up signals for graceful shutdown
	ctx := ctrl.SetupSignalHandler()

	// Create manager
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to create manager")
		os.Exit(1)
	}

	// Add schemes to manager
	if err := addSchemes(mgr); err != nil {
		setupLog.Error(err, "unable to add schemes")
		os.Exit(1)
	}

	// Create and register the reconciler
	reconciler := &controllers.GoProReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}

	if err := reconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	// Add healthz and readyz checks
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to add healthz check")
		os.Exit(1)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to add readyz check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")

	// Start the manager
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func addSchemes(mgr ctrl.Manager) error {
	if err := goproexamplecomv1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		return fmt.Errorf("failed to add gopro scheme: %w", err)
	}
	return nil
}
