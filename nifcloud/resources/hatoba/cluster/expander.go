package cluster

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
)

func expandCreateClusterInput(d *schema.ResourceData) *hatoba.CreateClusterInput {
	return &hatoba.CreateClusterInput{
		Cluster: &hatoba.CreateClusterRequestCluster{
			Name:              nifcloud.String(d.Get("name").(string)),
			Description:       nifcloud.String(d.Get("description").(string)),
			KubernetesVersion: nifcloud.String(d.Get("kubernetes_version").(string)),
			Locations:         expandLocations(d.Get("locations").([]interface{})),
			FirewallGroup:     nifcloud.String(d.Get("firewall_group").(string)),
			AddonsConfig:      expandAddonsConfig(d.Get("addons_config").([]interface{})),
			NetworkConfig:     expandNetworkConfig(d.Get("network_config").([]interface{})),
			NodePools:         expandNodePools(d.Get("node_pools").(*schema.Set).List()),
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

func expandAddonsConfig(raw []interface{}) *hatoba.AddonsConfig {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}

	config := raw[0].(map[string]interface{})
	res := &hatoba.AddonsConfig{}

	if value, ok := config["http_load_balancing"]; ok && len(value.([]interface{})) > 0 {
		addon := value.([]interface{})[0].(map[string]interface{})
		res.HttpLoadBalancing = &hatoba.HttpLoadBalancing{
			Disabled: nifcloud.Bool(addon["disabled"].(bool)),
		}
	}

	return res
}

func expandNetworkConfig(raw []interface{}) *hatoba.NetworkConfig {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}

	config := raw[0].(map[string]interface{})
	res := &hatoba.NetworkConfig{}

	if value, ok := config["network_id"]; ok {
		res.NetworkId = nifcloud.String(value.(string))
	}

	return res
}

func expandNodePools(raw []interface{}) []hatoba.CreateClusterRequestNodePool {
	if len(raw) == 0 {
		return nil
	}

	res := make([]hatoba.CreateClusterRequestNodePool, len(raw))
	for i, r := range raw {
		np := r.(map[string]interface{})
		res[i] = hatoba.CreateClusterRequestNodePool{
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
		Cluster: &hatoba.UpdateClusterRequestCluster{
			Description:       nifcloud.String(d.Get("description").(string)),
			KubernetesVersion: nifcloud.String(d.Get("kubernetes_version").(string)),
			AddonsConfig:      expandAddonsConfig(d.Get("addons_config").([]interface{})),
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
		NodePool: &hatoba.CreateClusterRequestNodePool{
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
