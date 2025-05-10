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

package v1

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	examplev1 "jsonserver-operator/api/v1"
)

// nolint:unused
// log is for logging in this package.
var jsonserverlog = logf.Log.WithName("jsonserver-resource")

// SetupJsonServerWebhookWithManager registers the webhook for JsonServer in the manager.
func SetupJsonServerWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&examplev1.JsonServer{}).
		WithValidator(&JsonServerCustomValidator{}).
		WithDefaulter(&JsonServerCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-example-example-com-v1-jsonserver,mutating=true,failurePolicy=fail,sideEffects=None,groups=example.example.com,resources=jsonservers,verbs=create;update,versions=v1,name=mjsonserver-v1.kb.io,admissionReviewVersions=v1

// JsonServerCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind JsonServer when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type JsonServerCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &JsonServerCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind JsonServer.
func (d *JsonServerCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	jsonserver, ok := obj.(*examplev1.JsonServer)

	if !ok {
		return fmt.Errorf("expected an JsonServer object but got %T", obj)
	}
	jsonserverlog.Info("Defaulting for JsonServer", "name", jsonserver.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-example-example-com-v1-jsonserver,mutating=false,failurePolicy=fail,sideEffects=None,groups=example.example.com,resources=jsonservers,verbs=create;update,versions=v1,name=vjsonserver-v1.kb.io,admissionReviewVersions=v1

// JsonServerCustomValidator struct is responsible for validating the JsonServer resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type JsonServerCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &JsonServerCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type JsonServer.
func (v *JsonServerCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	jsonserver, ok := obj.(*examplev1.JsonServer)
	if !ok {
		return nil, fmt.Errorf("expected a JsonServer object but got %T", obj)
	}
	jsonserverlog.Info("Validation for JsonServer upon creation", "name", jsonserver.GetName())

	// Validating webhook to block creation of objects not following naming convention
	if !strings.HasPrefix(jsonserver.GetName(), "app-") {
		return nil, fmt.Errorf("JsonServer name must follow the convention 'app-${name}'")
	}

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type JsonServer.
func (v *JsonServerCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	jsonserver, ok := newObj.(*examplev1.JsonServer)
	if !ok {
		return nil, fmt.Errorf("expected a JsonServer object for the newObj but got %T", newObj)
	}
	jsonserverlog.Info("Validation for JsonServer upon update", "name", jsonserver.GetName())

	// Validating webhook to block update of objects not following naming convention
	if !strings.HasPrefix(jsonserver.GetName(), "app-") {
		return nil, fmt.Errorf("JsonServer name must follow the convention 'app-${name}'")
	}

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type JsonServer.
func (v *JsonServerCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	jsonserver, ok := obj.(*examplev1.JsonServer)
	if !ok {
		return nil, fmt.Errorf("expected a JsonServer object but got %T", obj)
	}
	jsonserverlog.Info("Validation for JsonServer upon deletion", "name", jsonserver.GetName())

	return nil, nil
}
