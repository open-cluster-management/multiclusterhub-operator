package v1

import (
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
	"k8s.io/utils/strings/slices"
)

func TestMultiClusterHub_Prune(t *testing.T) {
	tests := []struct {
		name      string
		have      []ComponentConfig
		component string
		want      bool
		want2     []ComponentConfig
	}{
		{
			name: "solo component pruned",
			have: []ComponentConfig{
				{Name: ManagementIngress, Enabled: true},
			},
			component: ManagementIngress,
			want:      true,
			want2:     []ComponentConfig{},
		},
		{
			name: "component pruned",
			have: []ComponentConfig{
				{Name: ClusterLifecycle, Enabled: true},
				{Name: ManagementIngress, Enabled: true},
				{Name: Search, Enabled: true},
			},
			component: ManagementIngress,
			want:      true,
			want2: []ComponentConfig{
				{Name: ClusterLifecycle, Enabled: true},
				{Name: Search, Enabled: true},
			},
		},
		{
			name: "nothing to prune",
			have: []ComponentConfig{
				{Name: ClusterLifecycle, Enabled: true},
				{Name: Search, Enabled: true},
			},
			component: ManagementIngress,
			want:      false,
			want2: []ComponentConfig{
				{Name: ClusterLifecycle, Enabled: true},
				{Name: Search, Enabled: true},
			},
		},
		{
			name:      "nil list",
			have:      nil,
			component: ManagementIngress,
			want:      false,
			want2:     nil,
		},
	}
	for _, tt := range tests {
		mch := &MultiClusterHub{
			Spec: MultiClusterHubSpec{
				Overrides: &Overrides{
					Components: tt.have,
				},
			},
		}
		t.Run(tt.name, func(t *testing.T) {
			got := mch.Prune(tt.component)
			if got != tt.want {
				t.Errorf("MultiClusterHub.Prune() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(mch.Spec.Overrides.Components, tt.want2) {
				t.Errorf("MultiClusterHub.Prune() = %v, want %v", mch.Spec.Overrides.Components, tt.want2)
			}
		})
	}
}

func TestGetDisabledComponents(t *testing.T) {
	tests := []struct {
		name      string
		component string
		want      bool
		want2     int
	}{
		{
			name:      "default disabled components",
			component: ClusterBackup,
			want:      true,
			want2:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disabledComponents, err := GetDefaultDisabledComponents()

			if err != nil {
				t.Errorf("GetDefaultDisabledComponents() = %v, want: %v", err.Error(), nil)
			}

			pass := false
			for _, c := range disabledComponents {
				if c == tt.component {
					pass = true
				}
			}

			if !pass {
				t.Errorf("GetDefaultDisabledComponents() = %v, want: %v", pass, tt.want)
			}

			if len(disabledComponents) != 1 {
				t.Errorf("GetDefaultDisabledComponents() = %v, want: %v", len(disabledComponents), tt.want2)
			}
		})
	}
}

func TestGetDefaultHostedComponents(t *testing.T) {
	tests := []struct {
		name      string
		component string
		want      int
		want2     bool
	}{
		{
			name:      "default hosted components",
			component: MultiClusterEngine,
			want:      1,
			want2:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hostedComponents := GetDefaultHostedComponents()

			if len(hostedComponents) != 1 {
				t.Errorf("GetDefaultHostedComponents() = %v, want: %v", len(hostedComponents), tt.want)
			}

			pass := false
			for _, c := range hostedComponents {
				if c == tt.component {
					pass = true
				}
			}

			if !pass {
				t.Errorf("GetDefaultHostedComponents() = %v, want: %v", pass, tt.want2)
			}

		})
	}
}

func TestGetClusterManagementAddonName(t *testing.T) {
	tests := []struct {
		name      string
		component string
		want      string
	}{
		{
			name:      "submariner ClusterManagementAddOn",
			component: SubmarinerAddon,
			want:      "submariner",
		},
		{
			name:      "unknown ClusterManagementAddOn",
			component: "unknown",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClusterManagementAddonName(tt.component)
			if err != nil && tt.component != "unknown" {
				t.Errorf("GetClusterManagementAddonName(%v) = %v, want: %v", tt.component, err.Error(), tt.want)
			}

			if got != tt.want {
				t.Errorf("GetClusterManagementAddonName(%v) = %v, want: %v", tt.component, got, tt.want)
			}
		})
	}
}

func TestGetLegacyPrometheusKind(t *testing.T) {
	tests := []struct {
		name  string
		kind  string
		want  int
		want2 []string
	}{
		{
			name:  "legacy Prometheus Configuration Kind",
			kind:  "PrometheusRule",
			want:  3,
			want2: LegacyConfigKind,
		},
		{
			name:  "legacy Prometheus Configuration Kind",
			kind:  "ServiceMonitor",
			want:  3,
			want2: LegacyConfigKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLegacyConfigKind()
			if len(got) == 0 {
				t.Errorf("GetLegacyConfigKind() = %v, want: %v", len(got), tt.want)
			}

			if ok := slices.Contains(got, tt.kind); !ok {
				t.Errorf("GetLegacyConfigKind() = %v, want: %v", got, tt.want2)
			}
		})
	}
}

func TestGetLegacyPrometheusRulesName(t *testing.T) {
	tests := []struct {
		name      string
		component string
		want      string
	}{
		{
			name:      "console PrometheusRule",
			component: Console,
			want:      MCHLegacyPrometheusRules[Console],
		},
		{
			name:      "unknown PrometheusRule",
			component: "unknown",
			want:      MCHLegacyPrometheusRules["unknown"],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLegacyPrometheusRulesName(tt.component)
			if err != nil && tt.component != "unknown" {
				t.Errorf("GetLegacyPrometheusRulesName(%v) = %v, want: %v", tt.component, err.Error(), tt.want)
			}

			if got != tt.want {
				t.Errorf("GetLegacyPrometheusRulesName(%v) = %v, want: %v", tt.component, got, tt.want)
			}
		})
	}
}

func TestGetLegacyServiceMonitorName(t *testing.T) {
	tests := []struct {
		name      string
		component string
		want      string
	}{
		{
			name:      "console ServiceMonitor",
			component: Console,
			want:      MCHLegacyServiceMonitors[Console],
		},
		{
			name:      "unknown ServiceMonitor",
			component: "unknown",
			want:      MCHLegacyServiceMonitors["unknown"],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLegacyServiceMonitorName(tt.component)
			if err != nil && tt.component != "unknown" {
				t.Errorf("GetLegacyServiceMonitorName(%v) = %v, want: %v", tt.component, err.Error(), tt.want)
			}

			if got != tt.want {
				t.Errorf("GetLegacyServiceMonitorName(%v) = %v, want: %v", tt.component, got, tt.want)
			}
		})
	}
}

func TestHubSizeDefault(t *testing.T) {
	tests := []struct {
		name string
		spec MultiClusterHubSpec
		want HubSize
	}{
		{
			name: "Default is Medium",
			spec: MultiClusterHubSpec{},
			want: Medium,
		},
		{
			name: "Override Default with Large",
			spec: MultiClusterHubSpec{
				HubSize: Large,
			},
			want: Large,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hsize := tt.spec.HubSize
			if hsize != tt.want {
				t.Errorf("HubSize: %v, want: %v", hsize, tt.want)
			}
		})
	}
}

func (h HubSize) String() string {
	return HubSizeStrings[h]
}

func (h *HubSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	fmt.Println("Unmarshaling YAML is occuring")
	var hubsize string
	if err := unmarshal(&hubsize); err != nil {
		return err
	}

	var exists bool
	*h, exists = HubSizeFromString[hubsize]

	if !exists {
		return fmt.Errorf("key %v does not exist in map", hubsize)
	}
	return nil
}

func TestHubSizeMarshal(t *testing.T) {
	tests := []struct {
		name       string
		yamlstring string
		want       HubSize
	}{
		{
			name:       "Marshal defaults to M",
			yamlstring: ``,
			want:       Medium,
		},
		{
			name:       "Marshals when overriding default with large",
			yamlstring: `hubsize: L`, // For some reason, "hubSize" didn't work, but "hubsize" did. Go figure
			want:       Large,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out MultiClusterHubSpec
			err := yaml.Unmarshal([]byte(tt.yamlstring), &out)
			if err != nil {
				t.Errorf("Unable to unmarshal yaml string: %v. %v", tt.yamlstring, err)
			}
			if out.HubSize != tt.want {
				t.Errorf("Hubsize not desired. HubSize: %v, want: %v", out.HubSize, tt.want)
			}
		})
	}
}

func TestGetLegacyServiceName(t *testing.T) {
	tests := []struct {
		name      string
		component string
		want      string
	}{
		{
			name:      "unknown Service",
			component: "unknown",
			want:      MCHLegacyServices["unknown"],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLegacyServiceName(tt.component)
			if err != nil && tt.component != "unknown" {
				t.Errorf("GetLegacyServiceName(%v) = %v, want: %v", tt.component, err.Error(), tt.want)
			}

			if got != tt.want {
				t.Errorf("GetLegacyServiceName(%v) = %v, want: %v", tt.component, got, tt.want)
			}
		})
	}
}
