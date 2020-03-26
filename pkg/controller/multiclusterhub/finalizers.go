package multiclusterhub

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	operatorsv1alpha1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func (r *ReconcileMultiClusterHub) cleanupHiveConfigs(reqLogger logr.Logger, m *operatorsv1alpha1.MultiClusterHub) error {
	hiveConfigRes := schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "hiveconfigs"}

	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	listOptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("installer.name=%s", m.GetName()),
	}

	hiveResList, err := dc.Resource(hiveConfigRes).List(listOptions)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Hiveconfig resource not found. Continuing.")
			return nil
		}

		reqLogger.Error(err, "Error while listing hiveconfig instances")
		return err
	}

	for _, hiveRes := range hiveResList.Items {
		reqLogger.Info("Deleting hiveconfig", "Resource.Name", hiveRes.GetName())
		if err := dc.Resource(hiveConfigRes).Delete(hiveRes.GetName(), &deleteOptions); err != nil {
			reqLogger.Error(err, "Error while deleting hiveconfig instances")
			return err
		}
	}

	reqLogger.Info("Hiveconfigs finalized")
	return nil
}

func (r *ReconcileMultiClusterHub) cleanupAPIServices(reqLogger logr.Logger, m *operatorsv1alpha1.MultiClusterHub) error {
	err := r.client.DeleteAllOf(
		context.TODO(),
		&apiregistrationv1.APIService{},
		client.MatchingLabels{
			"installer.name":      m.GetName(),
			"installer.namespace": m.GetNamespace(),
		},
	)

	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("No matching API services to finalize. Continuing.")
			return nil
		}

		reqLogger.Error(err, "Error while deleting API services")
		return err
	}

	reqLogger.Info("API services finalized")
	return nil
}
