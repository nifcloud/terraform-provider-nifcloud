package record

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
)

func expandCreateChangeResourceRecordSetsInput(d *schema.ResourceData) *dns.ChangeResourceRecordSetsInput {
	return &dns.ChangeResourceRecordSetsInput{
		ZoneID: nifcloud.String(d.Get("zone_id").(string)),
		RequestChangeBatch: &dns.RequestChangeBatch{
			ListOfRequestChanges: []dns.RequestChanges{{
				RequestChange: &dns.RequestChange{
					Action:                   nifcloud.String("CREATE"),
					RequestResourceRecordSet: expandRequestResourceRecordSetInput(d),
				},
			}},
		},
	}
}

func expandListResourceRecordSets(d *schema.ResourceData) *dns.ListResourceRecordSetsInput {
	return &dns.ListResourceRecordSetsInput{
		Identifier: nifcloud.String(d.Id()),
		Name:       nifcloud.String(d.Get("name").(string)),
		Type:       nifcloud.String(d.Get("type").(string)),
		ZoneID:     nifcloud.String(d.Get("zone_id").(string)),
	}
}

func expandDeleteChangeResourceRecordSetsInput(d *schema.ResourceData) *dns.ChangeResourceRecordSetsInput {
	return &dns.ChangeResourceRecordSetsInput{
		ZoneID: nifcloud.String(d.Get("zone_id").(string)),
		RequestChangeBatch: &dns.RequestChangeBatch{
			ListOfRequestChanges: []dns.RequestChanges{{
				RequestChange: &dns.RequestChange{
					Action:                   nifcloud.String("DELETE"),
					RequestResourceRecordSet: expandRequestResourceRecordSetInput(d),
				},
			}},
		},
	}
}

func expandRequestResourceRecordSetInput(d *schema.ResourceData) *dns.RequestResourceRecordSet {
	input := &dns.RequestResourceRecordSet{
		Name:              nifcloud.String(d.Get("name").(string)),
		SetIdentifier:     nifcloud.String(d.Get("set_identifier").(string)),
		TTL:               nifcloud.Int64(int64(d.Get("ttl").(int))),
		Type:              nifcloud.String(d.Get("type").(string)),
		XniftyComment:     nifcloud.String(d.Get("comment").(string)),
		XniftyDefaultHost: nifcloud.String(d.Get("default_host").(string)),
		ListOfRequestResourceRecords: []dns.RequestResourceRecords{{
			RequestResourceRecord: &dns.RequestResourceRecord{
				Value: nifcloud.String(d.Get("record").(string)),
			},
		}},
	}

	weightSet := d.Get("weighted_routing_policy").([]interface{})
	if len(weightSet) != 0 {
		weight := weightSet[0].(map[string]interface{})

		if value, ok := weight["weight"]; ok {
			input.Weight = nifcloud.Int64(int64(value.(int)))
		}
	}

	failoverSet := d.Get("failover_routing_policy").([]interface{})
	if len(failoverSet) != 0 {
		failover := failoverSet[0].(map[string]interface{})

		if value, ok := failover["type"]; ok {
			input.Failover = nifcloud.String(value.(string))
		}

		if len(failover["health_check"].([]interface{})) != 0 {
			healthCheckSet := failover["health_check"].([]interface{})
			healthCheck := healthCheckSet[0].(map[string]interface{})
			input.RequestXniftyHealthCheckConfig = &dns.RequestXniftyHealthCheckConfig{
				FullyQualifiedDomainName: nifcloud.String(healthCheck["resource_domain"].(string)),
				IPAddress:                nifcloud.String(healthCheck["ip_address"].(string)),
				Port:                     nifcloud.Int64(int64(healthCheck["port"].(int))),
				Protocol:                 nifcloud.String(healthCheck["protocol"].(string)),
				ResourcePath:             nifcloud.String(healthCheck["resource_path"].(string)),
			}
		}
	}

	return input
}
