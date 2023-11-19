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
	"fmt"

	"k8s.io/apimachinery/pkg/types"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	provisionv1alpha1 "github.com/kubernetescode-operator/api/v1alpha1"
)

// ProvisionRequestReconciler reconciles a ProvisionRequest object
type ProvisionRequestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=provision.mydomain.com,resources=provisionrequests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=provision.mydomain.com,resources=provisionrequests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=provision.mydomain.com,resources=provisionrequests/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=create;get;list;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=create;get;list;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ProvisionRequest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *ProvisionRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog := log.FromContext(ctx)

	klog.Info(fmt.Sprintf("start to run sync logic for PR %s", req.NamespacedName))
	defer klog.Info(fmt.Sprintf("finish sync logic for PR %s", req.NamespacedName))

	pr := &provisionv1alpha1.ProvisionRequest{}
	err := r.Get(ctx, req.NamespacedName, pr)
	if errors.IsNotFound(err) {
		klog.Info(fmt.Sprintf("Provision Request %s has been deleted", req.NamespacedName))
		return ctrl.Result{}, nil
	}
	if err != nil {
		klog.Error(err, "Get provision request fail")
		return ctrl.Result{}, err
	}

	pr2 := pr.DeepCopy()

	custNameSpaceName := pr.Spec.NamespaceName
	err = r.Get(ctx, types.NamespacedName{Name: custNameSpaceName, Namespace: ""}, &v1.Namespace{})
	if errors.IsNotFound(err) {
		custNameSpace := v1.Namespace{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"},
			ObjectMeta: metav1.ObjectMeta{
				UID:         uuid.NewUUID(),
				Name:        custNameSpaceName,
				Annotations: make(map[string]string),
			},
			Spec: v1.NamespaceSpec{},
		}
		err = r.Create(ctx, &custNameSpace, &client.CreateOptions{})
		if err != nil {
			return ctrl.Result{}, &errors.StatusError{ErrStatus: metav1.Status{
				Status:  "Failure",
				Message: "fail to create customer namespace",
			}}
		}
	}

	if !pr.Status.DbReady {
		var replicas int32 = 1
		selector := map[string]string{}
		selector["type"] = "provisioinrequest"
		selector["company"] = pr.Labels["company"]

		d := apps.Deployment{
			TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
			ObjectMeta: metav1.ObjectMeta{
				UID:         uuid.NewUUID(),
				Name:        "cust-db",
				Namespace:   custNameSpaceName,
				Annotations: make(map[string]string),
			},
			Spec: apps.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{MatchLabels: selector},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: selector,
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:            "customer-db",
								Image:           "mysql:5.7",
								ImagePullPolicy: "IfNotPresent",
								Env:             []v1.EnvVar{{Name: "MYSQL_ROOT_PASSWORD", Value: "pleasechangetosecret"}},
								Ports:           []v1.ContainerPort{{ContainerPort: 3306, Name: "mysql"}},
							},
						},
					},
				},
			},
		}
		err = r.Get(ctx, types.NamespacedName{Name: d.Name, Namespace: custNameSpaceName}, &apps.Deployment{})
		if errors.IsNotFound(err) {
			err = r.Create(ctx, &d, &client.CreateOptions{})
			if err != nil {
				klog.Error(err, "Failed when creating DB deployment for Provision Request")
				return ctrl.Result{}, err
			}
		} else if err != nil {
			return ctrl.Result{}, &errors.StatusError{ErrStatus: metav1.Status{
				Status:  "Failure",
				Message: "fail to read DB deployment",
			}}
		}
	}

	if !pr.Status.IngressReady {
		// 这里省去配置Ingress的逻辑......
	}

	pr2.Status.IngressReady = true
	pr2.Status.DbReady = true
	pr2.Kind = "ProvisionRequest"
	err = r.Status().Update(context.TODO(), pr2, &client.SubResourceUpdateOptions{})
	if err != nil {
		klog.Error(err, "Fail to update request status")
		return ctrl.Result{}, &errors.StatusError{ErrStatus: metav1.Status{
			Status:  "Failure",
			Message: "fail to update provision request status",
		}}
	}

	klog.Info(fmt.Sprintf("Sucessfully fulfill provision request %s", req.NamespacedName))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProvisionRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&provisionv1alpha1.ProvisionRequest{}).
		Complete(r)
}
