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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

	// Create resources
	// ConfigMap for JSON data
	configMap, err := r.reconcileConfigMap(ctx, jsonServer)
	if err != nil {
		return r.updateStatus(ctx, jsonServer, "Error", "Error: unexpected failure")
	}

	// Deployment
	if err := r.reconcileDeployment(ctx, jsonServer, configMap); err != nil {
		return r.updateStatus(ctx, jsonServer, "Error", "Error: unexpected failure")
	}

	// Service
	if err := r.reconcileService(ctx, jsonServer); err != nil {
		return r.updateStatus(ctx, jsonServer, "Error", "Error: unexpected failure")
	}

	// Set Synced state
	return r.updateStatus(ctx, jsonServer, "Synced", "Synced succesfully!")
}

// SetupWithManager sets up the controller with the Manager.
func (r *JsonServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.JsonServer{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Named("jsonserver").
		Complete(r)
}

// Helper functions

// getResourceLabels returns the labels to be applied to resources owned by the JsonServer
func getResourceLabels(jsonServer *examplev1.JsonServer) map[string]string {
	return map[string]string{
		"app":        jsonServer.Name,
		"managed-by": "jsonserver-operator",
	}
}

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

// reconcileConfigMap ensures the ConfigMap exists
func (r *JsonServerReconciler) reconcileConfigMap(ctx context.Context, jsonServer *examplev1.JsonServer) (*corev1.ConfigMap, error) {
	log := logf.FromContext(ctx)

	// Define ConfigMap
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jsonServer.Name,
			Namespace: jsonServer.Namespace,
			Labels:    getResourceLabels(jsonServer),
		},
	}

	// Create or update ConfigMap
	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, configMap, func() error {
		// Set the owner reference
		if err := controllerutil.SetControllerReference(jsonServer, configMap, r.Scheme); err != nil {
			return err
		}

		// Set data
		if configMap.Data == nil {
			configMap.Data = make(map[string]string)
		}
		configMap.Data["db.json"] = jsonServer.Spec.JsonConfig

		return nil
	})

	if err != nil {
		log.Error(err, "Failed to create or update ConfigMap")
		return nil, err
	}

	log.Info("ConfigMap reconciled", "operation", op)
	return configMap, nil
}

func (r *JsonServerReconciler) reconcileDeployment(ctx context.Context, jsonServer *examplev1.JsonServer, configMap *corev1.ConfigMap) error {
	log := logf.FromContext(ctx)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jsonServer.Name,
			Namespace: jsonServer.Namespace,
			Labels:    getResourceLabels(jsonServer),
		},
	}

	// Create or update Deployment
	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		// Set the owner reference
		if err := controllerutil.SetControllerReference(jsonServer, deployment, r.Scheme); err != nil {
			return err
		}

		labels := getResourceLabels(jsonServer)

		deployment.Spec.Replicas = &jsonServer.Spec.Replicas
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
		deployment.Spec.Template = corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "json-server",
						Image: "backplane/json-server",
						Args:  []string{"/data/db.json"},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 3000,
								Name:          "http",
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "json-config",
								MountPath: "/data",
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "json-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: configMap.Name,
								},
							},
						},
					},
				},
			},
		}

		return nil
	})

	if err != nil {
		log.Error(err, "Failed to create or update Deployment")
		return err
	}

	log.Info("Deployment reconciled", "operation", op)
	return nil
}

func (r *JsonServerReconciler) reconcileService(ctx context.Context, jsonServer *examplev1.JsonServer) error {
	log := logf.FromContext(ctx)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jsonServer.Name,
			Namespace: jsonServer.Namespace,
			Labels:    getResourceLabels(jsonServer),
		},
	}

	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		if err := controllerutil.SetControllerReference(jsonServer, service, r.Scheme); err != nil {
			return err
		}

		labels := getResourceLabels(jsonServer)

		service.Spec.Selector = labels
		service.Spec.Ports = []corev1.ServicePort{
			{
				Port:       3000,
				TargetPort: intstr.FromInt(3000),
				Protocol:   corev1.ProtocolTCP,
			},
		}

		return nil
	})

	if err != nil {
		log.Error(err, "Failed to create or update Service")
		return err
	}

	log.Info("Service reconciled", "operation", op)
	return nil
}
