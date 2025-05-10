/*
Copyright 2025.

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
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "jsonserver-operator/api/v1"
)

// JsonServerReconciler reconciles a JsonServer object
type JsonServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.example.com,resources=jsonservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.example.com,resources=jsonservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.example.com,resources=jsonservers/finalizers,verbs=update

// RBAC to manage the custom resources (including delete so that it can cleanup the resources when the CRD is deleted)
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JsonServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *JsonServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling JsonServer", "name", req.NamespacedName)

	// Fetch the JsonServer instance
	jsonServer := &examplev1.JsonServer{}
	err := r.Get(ctx, req.NamespacedName, jsonServer)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("JsonServer resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - retry in next cycle
		log.Error(err, "Failed to get JsonServer")
		return ctrl.Result{}, err
	}

	if err := validateJSON(jsonServer.Spec.JsonConfig); err != nil {
		log.Error(err, "Invalid JSON configuration")
		return r.updateStatus(ctx, jsonServer, "Error", "Error: spec.jsonConfig is not a valid json object")
	}

	// Set Synced state
	return r.updateStatus(ctx, jsonServer, "Synced", "Synced succesfully!")
}

// SetupWithManager sets up the controller with the Manager.
func (r *JsonServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.JsonServer{}).
		Named("jsonserver").
		Complete(r)
}

// Helper functions

// validateJSON checks if the input string is a valid JSON
func validateJSON(input string) error {
	var js json.RawMessage
	return json.Unmarshal([]byte(input), &js)
}

// updateStatus updates the status of the JsonServer resource
func (r *JsonServerReconciler) updateStatus(ctx context.Context, jsonServer *examplev1.JsonServer, state, message string) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Update status if needed
	if jsonServer.Status.State != state || jsonServer.Status.Message != message {
		jsonServer.Status.State = state
		jsonServer.Status.Message = message
		if err := r.Status().Update(ctx, jsonServer); err != nil {
			log.Error(err, "Failed to update JsonServer status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

