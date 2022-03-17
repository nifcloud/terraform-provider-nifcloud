package record

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
)

func flatten(d *schema.ResourceData, res *dns.ListResourceRecordSetsResponse) error {
	if res == nil || len(res.ResourceRecordSets) == 0 {
		d.SetId("")
		return nil
	}

	var resourceRecordSet dns.ResourceRecordSets

	for _, s := range res.ResourceRecordSets {
		if nifcloud.StringValue(s.SetIdentifier) == d.Id() {
			resourceRecordSet = s
		}
	}

	if nifcloud.StringValue(resourceRecordSet.SetIdentifier) != d.Id() {
		return fmt.Errorf("unable to find dns record within: %#v", resourceRecordSet)
	}

	if err := d.Set("set_identifier", resourceRecordSet.SetIdentifier); err != nil {
		return err
	}

	if err := d.Set("name", resourceRecordSet.Name); err != nil {
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

func flattenWeight(record *dns.ResourceRecordSets) []map[string]interface{} {
	res := map[string]interface{}{}

	if record != nil && record.Weight != nil {
		res["weight"] = nifcloud.Int64Value(record.Weight)
	}

	return []map[string]interface{}{res}
}

func flattenFailover(record *dns.ResourceRecordSets) []map[string]interface{} {
	res := map[string]interface{}{}

	if record != nil && record.Failover != nil {
		res["type"] = nifcloud.StringValue(record.Failover)
		res["health_check"] = flattenHealthCheck(record.XniftyHealthCheckConfig)
	}

	return []map[string]interface{}{res}
}

func flattenHealthCheck(healthCheck *dns.XniftyHealthCheckConfig) []map[string]interface{} {
	res := map[string]interface{}{}

	if healthCheck != nil && healthCheck.Protocol != nil &&
		healthCheck.IPAddress != nil && healthCheck.Port != nil {
		res["protocol"] = nifcloud.StringValue(healthCheck.Protocol)
		res["ip_address"] = nifcloud.StringValue(healthCheck.IPAddress)
		res["port"] = nifcloud.Int64Value(healthCheck.Port)
		res["resource_path"] = nifcloud.StringValue(healthCheck.ResourcePath)
		res["resource_domain"] = nifcloud.StringValue(healthCheck.FullyQualifiedDomainName)
	}

	return []map[string]interface{}{res}
}
