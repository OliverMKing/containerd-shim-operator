/*
Copyright 2023.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	runtimev1alpha1 "github.com/olivermking/api/v1alpha1"
)

const (
	// finalizer is the name of the finalizer added to the shim CR
	finalizer = "shim.runtime.k8s.containerd.io/finalizer"
)

// ShimReconciler reconciles a Shim object
type ShimReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=runtime.k8s.containerd.io,resources=shims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=runtime.k8s.containerd.io,resources=shims/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=runtime.k8s.containerd.io,resources=shims/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Shim object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *ShimReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	lgr := log.FromContext(ctx).WithValues("shim", req.NamespacedName)

	shim := &runtimev1alpha1.Shim{}
	if err := r.Get(ctx, req.NamespacedName, shim); err != nil {
		// logged at Info because we can get things queued on deleted objects
		lgr.Info("unable to fetch Shim", "error", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if shim.ObjectMeta.DeletionTimestamp.IsZero() { // object is not being deleted, ensure finalizer
		if !controllerutil.ContainsFinalizer(shim, finalizer) {
			controllerutil.AddFinalizer(shim, finalizer)
			if err := r.Update(ctx, shim); err != nil {
				lgr.Error(err, "unable to update Shim with finalizer")
				return ctrl.Result{}, err
			}
		}
	} else { // object is being deleted, cleanup resources then remove finalizer
		if controllerutil.ContainsFinalizer(shim, finalizer) {
			// TODO: cleanup resources

			controllerutil.RemoveFinalizer(shim, finalizer)
			if err := r.Update(ctx, shim); err != nil {
				lgr.Error(err, "unable to remove finalizer from Shim")
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	// TODO: add reconcile logic

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ShimReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&runtimev1alpha1.Shim{}).
		Complete(r)
}
