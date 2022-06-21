// Copyright Contributors to the Open Cluster Management project

/*
Copyright 2021.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	mchlog = log.Log.WithName("multiclusterhub-resource")
	Client client.Client
)

func (r *MultiClusterHub) SetupWebhookWithManager(mgr ctrl.Manager) error {
	Client = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).For(r).Complete()
}

var _ webhook.Defaulter = &MultiClusterHub{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MultiClusterHub) Default() {
	mchlog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:name=multiclusterhub-operator-validating-webhook,path=/validate-v1-multiclusterhub,mutating=false,failurePolicy=fail,sideEffects=None,groups=operator.open-cluster-management.io,resources=multiclusterhubs,verbs=create;update;delete,versions=v1,name=multiclusterhub.validating-webhook.open-cluster-management.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &MultiClusterHub{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MultiClusterHub) ValidateCreate() error {
	mchlog.Info("validate create", "name", r.Name)
	multiClusterHubList := &MultiClusterHubList{}
	if err := Client.List(context.Background(), multiClusterHubList); err != nil {
		return fmt.Errorf("unable to list MultiClusterHubs: %s", err)
	}
	if len(multiClusterHubList.Items) == 0 {
		return nil
	}
	return fmt.Errorf("the MultiClusterHub CR already exists")
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MultiClusterHub) ValidateUpdate(old runtime.Object) error {
	mchlog.Info("validate update", "name", r.Name)
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MultiClusterHub) ValidateDelete() error {
	mchlog.Info("validate delete", "name", r.Name)
	return nil
}
