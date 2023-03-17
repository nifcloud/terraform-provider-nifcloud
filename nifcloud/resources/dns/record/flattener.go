package record

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns/types"
)

func flatten(d *schema.ResourceData, res *dns.ListResourceRecordSetsOutput) error {
	if res == nil || len(res.ResourceRecordSets) == 0 {
		d.SetId("")
		return nil
	}

	var resourceRecordSet types.ResourceRecordSets

	for _, s := range res.ResourceRecordSets {
		if nifcloud.ToString(s.SetIdentifier) == d.Id() {
			resourceRecordSet = s
		}
	}

	if nifcloud.ToString(resourceRecordSet.SetIdentifier) != d.Id() {
		return fmt.Errorf("unable to find dns record within: %#v", resourceRecordSet)
	}

	if err := d.Set("set_identifier", resourceRecordSet.SetIdentifier); err != nil {
		return err
	}

	if err := d.Set("name", flattenName(d, &resourceRecordSet)); err != nil {
		return err
	}

	if err := d.Set("type", resourceRecordSet.Type); err != nil {
		return err
	}

	if err := d.Set("record", resourceRecordSet.ResourceRecords[0].Value); err != nil {
		return err
	}

	if err := d.Set("ttl", resourceRecordSet.TTL); err != nil {
		return err
	}

	if _, ok := d.GetOk("weighted_routing_policy"); ok {
		if err := d.Set("weighted_routing_policy", flattenWeight(&resourceRecordSet)); err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("failover_routing_policy"); ok {
		if err := d.Set("failover_routing_policy", flattenFailover(&resourceRecordSet)); err != nil {
			return err
		}
	}

	if err := d.Set("default_host", resourceRecordSet.XniftyDefaultHost); err != nil {
		return err
	}

	if err := d.Set("comment", resourceRecordSet.XniftyComment); err != nil {
		return err
	}

	return nil
}

func flattenName(d *schema.ResourceData, record *types.ResourceRecordSets) interface{} {
	if *record.Name == d.Get("zone_id") {
		return "@"
	}
	return record.Name
}

func flattenWeight(record *types.ResourceRecordSets) []map[string]interface{} {
	res := map[string]interface{}{}

	if record != nil && record.Weight != nil {
		res["weight"] = nifcloud.ToInt32(record.Weight)
	}

	return []map[string]interface{}{res}
}

func flattenFailover(record *types.ResourceRecordSets) []map[string]interface{} {
	res := map[string]interface{}{}

	if record != nil && record.Failover != nil {
		res["type"] = nifcloud.ToString(record.Failover)
		res["health_check"] = flattenHealthCheck(record.XniftyHealthCheckConfig)
	}

	return []map[string]interface{}{res}
}

func flattenHealthCheck(healthCheck *types.XniftyHealthCheckConfig) []map[string]interface{} {
	res := map[string]interface{}{}

	if healthCheck != nil && healthCheck.Protocol != nil &&
		healthCheck.IPAddress != nil && healthCheck.Port != nil {
		res["protocol"] = nifcloud.ToString(healthCheck.Protocol)
		res["ip_address"] = nifcloud.ToString(healthCheck.IPAddress)
		res["port"] = nifcloud.ToInt32(healthCheck.Port)
		res["resource_path"] = nifcloud.ToString(healthCheck.ResourcePath)
		res["resource_domain"] = nifcloud.ToString(healthCheck.FullyQualifiedDomainName)
	}

	return []map[string]interface{}{res}
}
