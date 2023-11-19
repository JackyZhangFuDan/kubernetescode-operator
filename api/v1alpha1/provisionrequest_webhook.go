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

package v1alpha1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var provisionrequestlog = logf.Log.WithName("provisionrequest-resource")
var manager ctrl.Manager

func (r *ProvisionRequest) SetupWebhookWithManager(mgr ctrl.Manager) error {
	manager = mgr
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-provision-mydomain-com-v1alpha1-provisionrequest,mutating=true,failurePolicy=fail,sideEffects=None,groups=provision.mydomain.com,resources=provisionrequests,verbs=create;update,versions=v1alpha1,name=mprovisionrequest.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ProvisionRequest{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ProvisionRequest) Default() {
	provisionrequestlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-provision-mydomain-com-v1alpha1-provisionrequest,mutating=false,failurePolicy=fail,sideEffects=None,groups=provision.mydomain.com,resources=provisionrequests,verbs=create;update,versions=v1alpha1,name=vprovisionrequest.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ProvisionRequest{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ProvisionRequest) ValidateCreate() (admission.Warnings, error) {
	provisionrequestlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	company := r.GetLabels()["company"]
	req, err := labels.NewRequirement("company", selection.Equals, []string{company})
	if err != nil {
		return nil, err
	}

	clt := manager.GetClient()
	prs := &ProvisionRequestList{}
	err = clt.List(context.TODO(), prs, &client.ListOptions{LabelSelector: labels.NewSelector().Add(*req)})
	if err != nil {
		return nil, fmt.Errorf("failed to list provision request")
	}
	if len(prs.Items) > 0 {
		return nil, fmt.Errorf("the company already has provision request")
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ProvisionRequest) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	provisionrequestlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ProvisionRequest) ValidateDelete() (admission.Warnings, error) {
	provisionrequestlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
