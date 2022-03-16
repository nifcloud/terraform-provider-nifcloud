package cluster

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba/types"
)

func expandCreateClusterInput(d *schema.ResourceData) *hatoba.CreateClusterInput {
	return &hatoba.CreateClusterInput{
		Cluster: &types.RequestCluster{
			Name:                   nifcloud.String(d.Get("name").(string)),
			Description:            nifcloud.String(d.Get("description").(string)),
			KubernetesVersion:      types.KubernetesVersionOfclusterForCreateCluster(d.Get("kubernetes_version").(string)),
			ListOfRequestLocations: expandLocations(d.Get("locations").([]interface{})),
			FirewallGroup:          nifcloud.String(d.Get("firewall_group").(string)),
			RequestAddonsConfig:    expandAddonsConfig(d.Get("addons_config").([]interface{})),
			RequestNetworkConfig:   expandNetworkConfig(d.Get("network_config").([]interface{})),
			ListOfRequestNodePools: expandNodePools(d.Get("node_pools").(*schema.Set).List()),
		},
	}
}

func expandLocations(raw []interface{}) []string {
	if len(raw) == 0 {
		return nil
	}

	res := make([]string, len(raw))
	for i, l := range raw {
		res[i] = l.(string)
	}

	return res
}

func expandAddonsConfig(raw []interface{}) *types.RequestAddonsConfig {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}

	config := raw[0].(map[string]interface{})
	res := &types.RequestAddonsConfig{}

	if value, ok := config["http_load_balancing"]; ok && len(value.([]interface{})) > 0 {
		addon := value.([]interface{})[0].(map[string]interface{})
		res.RequestHttpLoadBalancing = &types.RequestHttpLoadBalancing{
			Disabled: nifcloud.Bool(addon["disabled"].(bool)),
		}
	}

	return res
}

func expandNetworkConfig(raw []interface{}) *types.RequestNetworkConfig {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}

	config := raw[0].(map[string]interface{})
	res := &types.RequestNetworkConfig{}

	if value, ok := config["network_id"]; ok {
		res.NetworkId = nifcloud.String(value.(string))
	}

	return res
}

func expandNodePools(raw []interface{}) []types.RequestNodePools {
	if len(raw) == 0 {
		return nil
	}

	res := make([]types.RequestNodePools, len(raw))
	for i, r := range raw {
		np := r.(map[string]interface{})
		res[i] = types.RequestNodePools{
			Name:         nifcloud.String(np["name"].(string)),
			InstanceType: types.InstanceTypeOfclusterForCreateCluster(np["instance_type"].(string)),
			NodeCount:    nifcloud.Int32(int32(np["node_count"].(int))),
		}
	}

	return res
}

func expandGetClusterInput(d *schema.ResourceData) *hatoba.GetClusterInput {
	return &hatoba.GetClusterInput{
		ClusterName: nifcloud.String(d.Id()),
	}
}

func expandUpdateClusterInput(d *schema.ResourceData) *hatoba.UpdateClusterInput {
	input := &hatoba.UpdateClusterInput{
		ClusterName: nifcloud.String(d.Id()),
		Cluster: &types.RequestClusterOfUpdateCluster{
			Description:         nifcloud.String(d.Get("description").(string)),
			KubernetesVersion:   types.KubernetesVersionOfclusterForUpdateCluster(d.Get("kubernetes_version").(string)),
			RequestAddonsConfig: expandAddonsConfig(d.Get("addons_config").([]interface{})),
		},
	}

	if d.HasChange("name") && !d.IsNewResource() {
		input.Cluster.Name = nifcloud.String(d.Get("name").(string))
	}

	return input
}

func expandCreateNodePoolInput(d *schema.ResourceData, targetNodePool map[string]interface{}) *hatoba.CreateNodePoolInput {
	return &hatoba.CreateNodePoolInput{
		ClusterName: nifcloud.String(d.Id()),
		NodePool: &types.RequestNodePool{
			Name:         nifcloud.String(targetNodePool["name"].(string)),
			InstanceType: types.InstanceTypeOfnodePoolForCreateNodePool(targetNodePool["instance_type"].(string)),
			NodeCount:    nifcloud.Int32(int32(targetNodePool["node_count"].(int))),
		},
	}
}

func expandDeleteNodePoolsInput(d *schema.ResourceData, targets []string) *hatoba.DeleteNodePoolsInput {
	return &hatoba.DeleteNodePoolsInput{
		ClusterName: nifcloud.String(d.Id()),
		Names:       nifcloud.String(strings.Join(targets, ",")),
	}
}

func expandSetNodePoolSizeInput(d *schema.ResourceData, targetNodePool map[string]interface{}) *hatoba.SetNodePoolSizeInput {
	return &hatoba.SetNodePoolSizeInput{
		ClusterName:  nifcloud.String(d.Id()),
		NodePoolName: nifcloud.String(targetNodePool["name"].(string)),
		NodeCount:    nifcloud.Int32(int32(targetNodePool["node_count"].(int))),
	}
}

func expandDeleteClusterInput(d *schema.ResourceData) *hatoba.DeleteClusterInput {
	return &hatoba.DeleteClusterInput{
		ClusterName: nifcloud.String(d.Id()),
	}
}

func expandGetClusterCredentialsInput(d *schema.ResourceData) *hatoba.GetClusterCredentialsInput {
	return &hatoba.GetClusterCredentialsInput{
		ClusterName: nifcloud.String(d.Id()),
	}
}
