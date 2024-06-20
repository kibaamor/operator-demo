/*
Copyright 2024.

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

package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var kdaemonlog = logf.Log.WithName("kdaemon-resource")

func (r *KDaemon) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-apps-kibazen-cn-v1alpha1-kdaemon,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps.kibazen.cn,resources=kdaemons,verbs=create;update,versions=v1alpha1,name=mkdaemon.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &KDaemon{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *KDaemon) Default() {
	kdaemonlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-apps-kibazen-cn-v1alpha1-kdaemon,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps.kibazen.cn,resources=kdaemons,verbs=create;update,versions=v1alpha1,name=vkdaemon.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &KDaemon{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *KDaemon) ValidateCreate() error {
	kdaemonlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	if r.Spec.Image == "" {
		return fmt.Errorf("image is required")
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KDaemon) ValidateUpdate(old runtime.Object) error {
	kdaemonlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	if r.Spec.Image == "" {
		return fmt.Errorf("image is required")
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KDaemon) ValidateDelete() error {
	kdaemonlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
