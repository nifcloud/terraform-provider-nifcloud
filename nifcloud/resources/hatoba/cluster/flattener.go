package cluster

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
)

func flatten(d *schema.ResourceData, res *hatoba.GetClusterResponse) error {
	if res == nil {
		d.SetId("")
		return nil
	}

	cluster := res.Cluster

	if nifcloud.StringValue(cluster.Name) != d.Id() {
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

func flattenAddonsConfig(c *hatoba.AddonsConfig) []map[string]interface{} {
	res := map[string]interface{}{}

	if c != nil && c.HttpLoadBalancing != nil {
		res["http_load_balancing"] = []map[string]interface{}{
			{
				"disabled": nifcloud.BoolValue(c.HttpLoadBalancing.Disabled),
			},
		}
	}

	return []map[string]interface{}{res}
}

func flattenNetworkConfig(c *hatoba.NetworkConfig) []map[string]interface{} {
	res := map[string]interface{}{}

	if c != nil && c.NetworkId != nil {
		res["network_id"] = nifcloud.StringValue(c.NetworkId)
	}

	return []map[string]interface{}{res}
}

func flattenNodePools(nodePools []hatoba.NodePool) []map[string]interface{} {
	res := make([]map[string]interface{}, len(nodePools))

	for i, nodePool := range nodePools {
		res[i] = map[string]interface{}{
			"name":          nifcloud.StringValue(nodePool.Name),
			"instance_type": nifcloud.StringValue(nodePool.InstanceType),
			"node_count":    nifcloud.Int64Value(nodePool.NodeCount),
			"nodes":         flattenNodes(nodePool.Nodes),
		}
	}

	return res
}

func flattenNodes(nodes []hatoba.Node) []map[string]interface{} {
	res := make([]map[string]interface{}, len(nodes))

	for i, n := range nodes {
		res[i] = map[string]interface{}{
			"name":               nifcloud.StringValue(n.Name),
			"availability_zone":  nifcloud.StringValue(n.AvailabilityZone),
			"public_ip_address":  nifcloud.StringValue(n.PublicIpAddress),
			"private_ip_address": nifcloud.StringValue(n.PrivateIpAddress),
		}
	}

	return res
}
