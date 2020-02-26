package rendering

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	operatorsv1alpha1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1alpha1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/rendering/patching"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/rendering/templates"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/kustomize/v3/pkg/resource"
)

const (
	apiserviceName         = "mcm-apiserver"
	controllerName         = "mcm-controller"
	webhookName            = "mcm-webhook"
	clusterControllerName  = "multicloud-operators-cluster-controller"
	helmRepoName           = "multicloudhub-repo"
	topologyAggregatorName = "topology-aggregator"
	metadataErr            = "failed to find metadata field"
)

var log = logf.Log.WithName("renderer")

type renderFn func(*resource.Resource) (*unstructured.Unstructured, error)

type Renderer struct {
	cr        *operatorsv1alpha1.MultiCloudHub
	renderFns map[string]renderFn
}

func NewRenderer(multipleCloudHub *operatorsv1alpha1.MultiCloudHub) *Renderer {
	renderer := &Renderer{cr: multipleCloudHub}
	renderer.renderFns = map[string]renderFn{
		"APIService":                   renderer.renderAPIServices,
		"Deployment":                   renderer.renderDeployments,
		"Service":                      renderer.renderNamespace,
		"ServiceAccount":               renderer.renderNamespace,
		"ConfigMap":                    renderer.renderNamespace,
		"ClusterRoleBinding":           renderer.renderClusterRoleBinding,
		"MutatingWebhookConfiguration": renderer.renderMutatingWebhookConfiguration,
		"Secret":                       renderer.renderSecret,
		"Subscription":                 renderer.renderSubscription,
		"EtcdCluster":                  renderer.renderBaseMetadataNamespace,
		"StatefulSet":                  renderer.renderBaseMetadataNamespace,
		"ClusterServiceVersion":        renderer.renderBaseMetadataNamespace,
		"Channel":                      renderer.renderBaseMetadataNamespace,
		"HiveConfig":                   renderer.renderHiveConfig,
		"SecurityContextConstraints":   renderer.renderSecContextConstraints,
	}
	return renderer
}

func (r *Renderer) Render() ([]*unstructured.Unstructured, error) {
	templates, err := templates.GetTemplateRenderer().GetTemplates(r.cr)
	if err != nil {
		return nil, err
	}
	resources, err := r.renderTemplates(templates)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *Renderer) renderTemplates(templates []*resource.Resource) ([]*unstructured.Unstructured, error) {
	uobjs := []*unstructured.Unstructured{}
	for _, template := range templates {
		render, ok := r.renderFns[template.GetKind()]
		if !ok {
			uobjs = append(uobjs, &unstructured.Unstructured{Object: template.Map()})
			continue
		}
		uobj, err := render(template.DeepCopy())
		if err != nil {
			return []*unstructured.Unstructured{}, err
		}
		if uobj == nil {
			continue
		}
		uobjs = append(uobjs, uobj)

	}

	// mcm-mutating-webhook has a section `webhooks.clientConfig.caBundle` created in secret mcm-webhook-secrets.
	// Current render mechanism cannot handle this dependence scenario. So re-render template dependence in the end.
	uobjs, err := reRenderDependence(uobjs)
	return uobjs, err
}

func (r *Renderer) renderAPIServices(res *resource.Resource) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: res.Map()}
	spec, ok := u.Object["spec"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to find apiservices spec field")
	}
	spec["service"] = map[string]interface{}{
		"namespace": r.cr.Namespace,
		"name":      apiserviceName,
	}
	return u, nil
}

func (r *Renderer) renderNamespace(res *resource.Resource) (*unstructured.Unstructured, error) {
	res.SetNamespace(r.cr.Namespace)
	return &unstructured.Unstructured{Object: res.Map()}, nil
}

func (r *Renderer) renderDeployments(res *resource.Resource) (*unstructured.Unstructured, error) {
	err := patching.ApplyGlobalPatches(res, r.cr)
	if err != nil {
		return nil, err
	}

	res.SetNamespace(r.cr.Namespace)

	name := res.GetName()
	switch name {
	case apiserviceName:
		if err := patching.ApplyAPIServerPatches(res, r.cr); err != nil {
			return nil, err
		}
		return &unstructured.Unstructured{Object: res.Map()}, nil
	case controllerName:
		if err := patching.ApplyControllerPatches(res, r.cr); err != nil {
			return nil, err
		}
		return &unstructured.Unstructured{Object: res.Map()}, nil
	case webhookName:
		if err := patching.ApplyWebhookPatches(res, r.cr); err != nil {
			return nil, err
		}
		return &unstructured.Unstructured{Object: res.Map()}, nil
	case clusterControllerName:
		return &unstructured.Unstructured{Object: res.Map()}, nil
	case helmRepoName:
		return &unstructured.Unstructured{Object: res.Map()}, nil
	case topologyAggregatorName:
		if err := patching.ApplyTopologyAggregatorPatches(res, r.cr); err != nil {
			return nil, err
		}
		return &unstructured.Unstructured{Object: res.Map()}, nil
	default:
		return nil, fmt.Errorf("unknown MultipleCloudHub deployment component %s", name)
	}
}

func (r *Renderer) renderClusterRoleBinding(res *resource.Resource) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: res.Map()}
	subjects, ok := u.Object["subjects"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to find clusterrolebinding subjects field")
	}
	subject := subjects[0].(map[string]interface{})
	kind := subject["kind"]
	if kind == "Group" {
		return u, nil
	}
	subject["namespace"] = r.cr.Namespace
	return u, nil
}

func (r *Renderer) renderMutatingWebhookConfiguration(res *resource.Resource) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: res.Map()}
	webooks, ok := u.Object["webhooks"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to find webhooks spec field")
	}
	webhook := webooks[0].(map[string]interface{})
	clientConfig := webhook["clientConfig"].(map[string]interface{})
	service := clientConfig["service"].(map[string]interface{})

	service["namespace"] = r.cr.Namespace
	return u, nil
}

func (r *Renderer) renderBaseMetadataNamespace(res *resource.Resource) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: res.Map()}
	metadata, ok := u.Object["metadata"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(metadataErr)
	}

	metadata["namespace"] = r.cr.Namespace
	return u, nil
}

func (r *Renderer) renderSubscription(res *resource.Resource) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: res.Map()}

	// Default to base renderFunc if not an IBM subscription
	if u.GetAPIVersion() != "app.ibm.com/v1alpha1" {
		return r.renderBaseMetadataNamespace(res)
	}

	// IBM subscription handling
	metadata, ok := u.Object["metadata"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(metadataErr)
	}

	metadata["namespace"] = r.cr.Namespace

	// Update channel to prepend the CRs namespace
	spec, ok := u.Object["spec"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to find ibm subscription spec field")
	}
	spec["channel"] = fmt.Sprintf("%s/%s", r.cr.Namespace, spec["channel"])

	imageTagPostfix := r.cr.Spec.ImageTagPostfix
	if imageTagPostfix != "" {
		imageTagPostfix = "-" + imageTagPostfix
	}

	// Check if contains a packageOverrides
	packageOverrides, ok := spec["packageOverrides"].([]interface{})
	if ok {
		for i := 0; i < len(packageOverrides); i++ {
			packageOverride, ok := packageOverrides[i].(map[string]interface{})
			if ok {
				override := packageOverride["packageOverrides"].([]interface{})
				for j := 0; j < len(override); j++ {
					packageData, _ := override[j].(map[string]interface{})
					packageData["value"] = strings.ReplaceAll(packageData["value"].(string), "{{POSTFIX}}", imageTagPostfix)
					packageData["value"] = strings.ReplaceAll(packageData["value"].(string), "{{IMAGEREPO}}", r.cr.Spec.ImageRepository)
					packageData["value"] = strings.ReplaceAll(packageData["value"].(string), "{{PULLSECRET}}", r.cr.Spec.ImagePullSecret)
					packageData["value"] = strings.ReplaceAll(packageData["value"].(string), "{{NAMESPACE}}", r.cr.Namespace)
					packageData["value"] = strings.ReplaceAll(packageData["value"].(string), "{{PULLPOLICY}}", string(r.cr.Spec.ImagePullPolicy))
				}
			}
		}
	}
	return u, nil
}

func (r *Renderer) renderSecret(res *resource.Resource) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: res.Map()}
	metadata, ok := u.Object["metadata"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to find metadata field")
	}
	data, ok := u.Object["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(metadataErr)
	}

	metadata["namespace"] = r.cr.Namespace

	name := res.GetName()

	switch name {
	case "mcm-apiserver-self-signed-secrets":
		ca, err := utils.GenerateSelfSignedCACert("multicloudhub-api")
		if err != nil {
			return nil, err
		}
		alternateDNS := []string{
			fmt.Sprintf("%s.%s", apiserviceName, r.cr.Namespace),
			fmt.Sprintf("%s.%s.svc", apiserviceName, r.cr.Namespace),
		}
		cert, err := utils.GenerateSignedCert(apiserviceName, alternateDNS, ca)
		if err != nil {
			return nil, err
		}
		data["ca.crt"] = []byte(ca.Cert)
		data["tls.crt"] = []byte(cert.Cert)
		data["tls.key"] = []byte(cert.Key)

		return u, nil
	case "mcm-klusterlet-self-signed-secrets":
		ca, err := utils.GenerateSelfSignedCACert("multicloudhub-klusterlet")
		if err != nil {
			return nil, err
		}
		cert, err := utils.GenerateSignedCert("multicloudhub-klusterlet", []string{}, ca)
		if err != nil {
			return nil, err
		}
		data["ca.crt"] = []byte(ca.Cert)
		data["tls.crt"] = []byte(cert.Cert)
		data["tls.key"] = []byte(cert.Key)
		return u, nil
	case "topology-aggregator-secret":
		ca, err := utils.GenerateSelfSignedCACert("topology-aggregator")
		if err != nil {
			return nil, err
		}
		alternateDNS := []string{
			fmt.Sprintf("%s.%s", topologyAggregatorName, r.cr.Namespace),
			fmt.Sprintf("%s.%s.svc", topologyAggregatorName, r.cr.Namespace),
		}
		cert, err := utils.GenerateSignedCert(topologyAggregatorName, alternateDNS, ca)
		if err != nil {
			return nil, err
		}
		data["ca.crt"] = []byte(ca.Cert)
		data["tls.crt"] = []byte(cert.Cert)
		data["tls.key"] = []byte(cert.Key)
		return u, nil
	case "mcm-webhook-secret":
		ca, err := utils.GenerateSelfSignedCACert("mcm-webhook")
		if err != nil {
			return nil, err
		}
		cert, err := utils.GenerateSignedCert(webhookName, []string{}, ca)
		if err != nil {
			return nil, err
		}
		data["ca.crt"] = []byte(ca.Cert)
		data["tls.crt"] = []byte(cert.Cert)
		data["tls.key"] = []byte(cert.Key)
		return u, nil
	}

	return u, nil
}

func (r *Renderer) renderHiveConfig(res *resource.Resource) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: res.Map()}
	metadata, ok := u.Object["metadata"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(metadataErr)
	}

	metadata["namespace"] = r.cr.Namespace
	u.Object["spec"] = structs.Map(r.cr.Spec.Hive)

	return u, nil
}

func reRenderDependence(objs []*unstructured.Unstructured) ([]*unstructured.Unstructured, error) {
	var ca interface{}
	var mutatingConfig *unstructured.Unstructured
	for _, obj := range objs {
		if obj.GetKind() == "Secret" && obj.GetName() == "mcm-webhook-secret" {
			data, ok := obj.Object["data"].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("failed to get ca in mcm-webhook-secrets")
			}
			ca = data["ca.crt"]
		}

		if obj.GetKind() == "MutatingWebhookConfiguration" && obj.GetName() == "mcm-mutating-webhook" {
			mutatingConfig = obj
		}
	}

	if ca != nil && mutatingConfig != nil {
		webooks, ok := mutatingConfig.Object["webhooks"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to find webhooks spec field")
		}
		webhook := webooks[0].(map[string]interface{})
		clientConfig := webhook["clientConfig"].(map[string]interface{})
		clientConfig["caBundle"] = ca
	}

	return objs, nil
}

func (r *Renderer) renderSecContextConstraints(res *resource.Resource) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: res.Map()}
	users, ok := u.Object["users"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to find users field")
	}
	ns := r.cr.Namespace
	users[0] = fmt.Sprintf("system:serviceaccount:%s:default", ns)
	return u, nil
}
