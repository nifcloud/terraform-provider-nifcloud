package cluster

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
)

func expandCreateClusterInput(d *schema.ResourceData) *hatoba.CreateClusterInput {
	return &hatoba.CreateClusterInput{
		Cluster: &hatoba.RequestCluster{
			Name:                   nifcloud.String(d.Get("name").(string)),
			Description:            nifcloud.String(d.Get("description").(string)),
			KubernetesVersion:      nifcloud.String(d.Get("kubernetes_version").(string)),
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

func expandAddonsConfig(raw []interface{}) *hatoba.RequestAddonsConfig {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}

	config := raw[0].(map[string]interface{})
	res := &hatoba.RequestAddonsConfig{}

	if value, ok := config["http_load_balancing"]; ok && len(value.([]interface{})) > 0 {
		addon := value.([]interface{})[0].(map[string]interface{})
		res.RequestHttpLoadBalancing = &hatoba.RequestHttpLoadBalancing{
			Disabled: nifcloud.Bool(addon["disabled"].(bool)),
		}
	}

	return res
}

func expandNetworkConfig(raw []interface{}) *hatoba.RequestNetworkConfig {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}

	config := raw[0].(map[string]interface{})
	res := &hatoba.RequestNetworkConfig{}

	if value, ok := config["network_id"]; ok {
		res.NetworkId = nifcloud.String(value.(string))
	}

	return res
}

func expandNodePools(raw []interface{}) []hatoba.RequestNodePools {
	if len(raw) == 0 {
		return nil
	}

	res := make([]hatoba.RequestNodePools, len(raw))
	for i, r := range raw {
		np := r.(map[string]interface{})
		res[i] = hatoba.RequestNodePools{
			Name:         nifcloud.String(np["name"].(string)),
			InstanceType: nifcloud.String(np["instance_type"].(string)),
			NodeCount:    nifcloud.Int64(int64(np["node_count"].(int))),
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
		Cluster: &hatoba.RequestClusterOfUpdateCluster{
			Description:         nifcloud.String(d.Get("description").(string)),
			KubernetesVersion:   nifcloud.String(d.Get("kubernetes_version").(string)),
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
		NodePool: &hatoba.RequestNodePool{
			Name:         nifcloud.String(targetNodePool["name"].(string)),
			InstanceType: nifcloud.String(targetNodePool["instance_type"].(string)),
			NodeCount:    nifcloud.Int64(int64(targetNodePool["node_count"].(int))),
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
		NodeCount:    nifcloud.Int64(int64(targetNodePool["node_count"].(int))),
	}
}

func expandDeleteClusterInput(d *schema.ResourceData) *hatoba.DeleteClusterInput {
	return &hatoba.DeleteClusterInput{
		ClusterName: nifcloud.String(d.Id()),
	}
}
