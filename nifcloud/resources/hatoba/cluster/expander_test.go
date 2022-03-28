package cluster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateClusterInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":               "test_name",
		"description":        "test_description",
		"kubernetes_version": "test_kubernetes_version",
		"locations":          []interface{}{"test_location"},
		"firewall_group":     "test_firewall_group",
		"addons_config": []interface{}{map[string]interface{}{
			"http_load_balancing": []interface{}{map[string]interface{}{
				"disabled": true,
			}},
		}},
		"network_config": []interface{}{map[string]interface{}{
			"network_id": "test_network_id",
		}},
		"node_pools": []interface{}{map[string]interface{}{
			"name":          "test_node_pool_name",
			"instance_type": "test_instance_type",
			"node_count":    3,
		}},
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.CreateClusterInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.CreateClusterInput{
				Cluster: &types.RequestCluster{
					Name:                   nifcloud.String("test_name"),
					Description:            nifcloud.String("test_description"),
					ListOfRequestLocations: []string{"test_location"},
					KubernetesVersion:      types.KubernetesVersionOfclusterForCreateCluster("test_kubernetes_version"),
					FirewallGroup:          nifcloud.String("test_firewall_group"),
					RequestAddonsConfig: &types.RequestAddonsConfig{
						RequestHttpLoadBalancing: &types.RequestHttpLoadBalancing{
							Disabled: nifcloud.Bool(true),
						},
					},
					RequestNetworkConfig: &types.RequestNetworkConfig{
						NetworkId: nifcloud.String("test_network_id"),
					},
					ListOfRequestNodePools: []types.RequestNodePools{
						{
							Name:         nifcloud.String("test_node_pool_name"),
							InstanceType: types.InstanceTypeOfclusterForCreateCluster("test_instance_type"),
							NodeCount:    nifcloud.Int32(3),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateClusterInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetClusterInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.GetClusterInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.GetClusterInput{
				ClusterName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetClusterInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateClusterInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":               "test_new_name",
		"description":        "test_description",
		"kubernetes_version": "test_kubernetes_version",
		"addons_config": []interface{}{map[string]interface{}{
			"http_load_balancing": []interface{}{map[string]interface{}{
				"disabled": false,
			}},
		}},
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.UpdateClusterInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.UpdateClusterInput{
				ClusterName: nifcloud.String("test_name"),
				Cluster: &types.RequestClusterOfUpdateCluster{
					Name:              nifcloud.String("test_new_name"),
					Description:       nifcloud.String("test_description"),
					KubernetesVersion: types.KubernetesVersionOfclusterForUpdateCluster("test_kubernetes_version"),
					RequestAddonsConfig: &types.RequestAddonsConfig{
						RequestHttpLoadBalancing: &types.RequestHttpLoadBalancing{
							Disabled: nifcloud.Bool(false),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateClusterInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandCreateNodePoolInput(t *testing.T) {
	nodePool := map[string]interface{}{
		"name":          "test_node_pool_name",
		"instance_type": "test_instance_type",
		"node_count":    5,
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"node_pools": []interface{}{nodePool},
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.CreateNodePoolInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.CreateNodePoolInput{
				ClusterName: nifcloud.String("test_name"),
				NodePool: &types.RequestNodePool{
					Name:         nifcloud.String("test_node_pool_name"),
					InstanceType: types.InstanceTypeOfnodePoolForCreateNodePool("test_instance_type"),
					NodeCount:    nifcloud.Int32(5),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateNodePoolInput(tt.args, nodePool)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteNodePoolsInput(t *testing.T) {
	nodePools := []string{"test_node_pool01", "test_node_pool02"}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.DeleteNodePoolsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.DeleteNodePoolsInput{
				ClusterName: nifcloud.String("test_name"),
				Names:       nifcloud.String("test_node_pool01,test_node_pool02"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteNodePoolsInput(tt.args, nodePools)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandSetNodePoolSizeInput(t *testing.T) {
	nodePool := map[string]interface{}{
		"name":       "test_node_pool_name",
		"node_count": 1,
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.SetNodePoolSizeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.SetNodePoolSizeInput{
				ClusterName:  nifcloud.String("test_name"),
				NodePoolName: nifcloud.String("test_node_pool_name"),
				NodeCount:    nifcloud.Int32(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandSetNodePoolSizeInput(tt.args, nodePool)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteClusterInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.DeleteClusterInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.DeleteClusterInput{
				ClusterName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteClusterInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
