package cluster

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba/types"
)

func flatten(d *schema.ResourceData, res *hatoba.GetClusterOutput) error {
	if res == nil {
		d.SetId("")
		return nil
	}

	cluster := res.Cluster

	if nifcloud.ToString(cluster.Name) != d.Id() {
		return fmt.Errorf("unable to find Hatoba cluster within: %#v", cluster)
	}

	if err := d.Set("name", cluster.Name); err != nil {
		return err
	}

	if err := d.Set("description", cluster.Description); err != nil {
		return err
	}

	if err := d.Set("locations", cluster.Locations); err != nil {
		return err
	}

	if err := d.Set("kubernetes_version", cluster.KubernetesVersion); err != nil {
		return err
	}

	if err := d.Set("firewall_group", cluster.FirewallGroup); err != nil {
		return err
	}

	if err := d.Set("addons_config", flattenAddonsConfig(cluster.AddonsConfig)); err != nil {
		return err
	}

	if err := d.Set("network_config", flattenNetworkConfig(cluster.NetworkConfig)); err != nil {
		return err
	}

	if err := d.Set("node_pools", flattenNodePools(cluster.NodePools)); err != nil {
		return err
	}

	return nil
}

func flattenAddonsConfig(c *types.AddonsConfig) []map[string]interface{} {
	res := map[string]interface{}{}

	if c != nil && c.HttpLoadBalancing != nil {
		res["http_load_balancing"] = []map[string]interface{}{
			{
				"disabled": nifcloud.ToBool(c.HttpLoadBalancing.Disabled),
			},
		}
	}

	return []map[string]interface{}{res}
}

func flattenNetworkConfig(c *types.NetworkConfig) []map[string]interface{} {
	res := map[string]interface{}{}

	if c != nil && c.NetworkId != nil {
		res["network_id"] = nifcloud.ToString(c.NetworkId)
	}

	return []map[string]interface{}{res}
}

func flattenNodePools(nodePools []types.NodePools) []map[string]interface{} {
	res := make([]map[string]interface{}, len(nodePools))

	for i, nodePool := range nodePools {
		res[i] = map[string]interface{}{
			"name":          nifcloud.ToString(nodePool.Name),
			"instance_type": nifcloud.ToString(nodePool.InstanceType),
			"node_count":    nifcloud.ToInt32(nodePool.NodeCount),
			"nodes":         flattenNodes(nodePool.Nodes),
		}
	}

	return res
}

func flattenNodes(nodes []types.Nodes) []map[string]interface{} {
	res := make([]map[string]interface{}, len(nodes))

	for i, n := range nodes {
		res[i] = map[string]interface{}{
			"name":               nifcloud.ToString(n.Name),
			"availability_zone":  nifcloud.ToString(n.AvailabilityZone),
			"public_ip_address":  nifcloud.ToString(n.PublicIpAddress),
			"private_ip_address": nifcloud.ToString(n.PrivateIpAddress),
		}
	}

	return res
}

func flattenCredentials(d *schema.ResourceData, res *hatoba.GetClusterCredentialsOutput) error {
	return d.Set("kube_config_raw", res.Credentials)
}
