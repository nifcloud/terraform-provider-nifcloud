package record

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns/types"
)

func expandCreateChangeResourceRecordSetsInput(d *schema.ResourceData) *dns.ChangeResourceRecordSetsInput {
	return &dns.ChangeResourceRecordSetsInput{
		ZoneID: nifcloud.String(d.Get("zone_id").(string)),
		RequestChangeBatch: &types.RequestChangeBatch{
			ListOfRequestChanges: []types.RequestChanges{{
				RequestChange: &types.RequestChange{
					Action:                   types.ActionOfChangeResourceRecordSetsRequestForChangeResourceRecordSetsCreate,
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
		Type:       types.TypeOfListResourceRecordSetsRequest(d.Get("type").(string)),
		ZoneID:     nifcloud.String(d.Get("zone_id").(string)),
	}
}

func expandDeleteChangeResourceRecordSetsInput(d *schema.ResourceData) *dns.ChangeResourceRecordSetsInput {
	return &dns.ChangeResourceRecordSetsInput{
		ZoneID: nifcloud.String(d.Get("zone_id").(string)),
		RequestChangeBatch: &types.RequestChangeBatch{
			ListOfRequestChanges: []types.RequestChanges{{
				RequestChange: &types.RequestChange{
					Action:                   types.ActionOfChangeResourceRecordSetsRequestForChangeResourceRecordSetsDelete,
					RequestResourceRecordSet: expandRequestResourceRecordSetInput(d),
				},
			}},
		},
	}
}

func expandRequestResourceRecordSetInput(d *schema.ResourceData) *types.RequestResourceRecordSet {
	input := &types.RequestResourceRecordSet{
		Name:              nifcloud.String(d.Get("name").(string)),
		SetIdentifier:     nifcloud.String(d.Get("set_identifier").(string)),
		TTL:               nifcloud.Int32(int32(d.Get("ttl").(int))),
		Type:              types.TypeOfChangeResourceRecordSetsRequestForChangeResourceRecordSets(d.Get("type").(string)),
		XniftyComment:     nifcloud.String(d.Get("comment").(string)),
		XniftyDefaultHost: nifcloud.String(d.Get("default_host").(string)),
		ListOfRequestResourceRecords: []types.RequestResourceRecords{{
			RequestResourceRecord: &types.RequestResourceRecord{
				Value: nifcloud.String(d.Get("record").(string)),
			},
		}},
	}

	weightSet := d.Get("weighted_routing_policy").([]interface{})
	if len(weightSet) != 0 {
		weight := weightSet[0].(map[string]interface{})

		if value, ok := weight["weight"]; ok {
			input.Weight = nifcloud.Int32(int32(value.(int)))
		}
	}

	failoverSet := d.Get("failover_routing_policy").([]interface{})
	if len(failoverSet) != 0 {
		failover := failoverSet[0].(map[string]interface{})

		if value, ok := failover["type"]; ok {
			input.Failover = types.FailoverOfChangeResourceRecordSetsRequestForChangeResourceRecordSets(value.(string))
		}

		if len(failover["health_check"].([]interface{})) != 0 {
			healthCheckSet := failover["health_check"].([]interface{})
			healthCheck := healthCheckSet[0].(map[string]interface{})
			input.RequestXniftyHealthCheckConfig = &types.RequestXniftyHealthCheckConfig{
				FullyQualifiedDomainName: nifcloud.String(healthCheck["resource_domain"].(string)),
				IPAddress:                nifcloud.String(healthCheck["ip_address"].(string)),
				Port:                     nifcloud.Int32(int32(healthCheck["port"].(int))),
				Protocol:                 types.ProtocolOfChangeResourceRecordSetsRequestForChangeResourceRecordSets(healthCheck["protocol"].(string)),
				ResourcePath:             nifcloud.String(healthCheck["resource_path"].(string)),
			}
		}
	}

	return input
}
