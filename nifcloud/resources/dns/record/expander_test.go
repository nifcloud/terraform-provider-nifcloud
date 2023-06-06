package record

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateChangeResourceRecordSetsInput(t *testing.T) {
	healthCheck := map[string]interface{}{
		"protocol":        "HTTP",
		"ip_address":      "192.0.2.1",
		"port":            8080,
		"resource_path":   "test_resource_path",
		"resource_domain": "test_resource_domain",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"zone_id": "test_zone_id",
		"name":    "test_name",
		"type":    "A",
		"record":  "192.0.2.1",
		"ttl":     60,
		"weighted_routing_policy": []interface{}{map[string]interface{}{
			"weight": 60,
		}},
		"failover_routing_policy": []interface{}{map[string]interface{}{
			"type":         "PRIMARY",
			"health_check": []interface{}{healthCheck},
		}},
		"comment":        "test_comment",
		"set_identifier": "test_set_identifier",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *dns.ChangeResourceRecordSetsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &dns.ChangeResourceRecordSetsInput{
				ZoneID: nifcloud.String("test_zone_id"),
				RequestChangeBatch: &types.RequestChangeBatch{
					ListOfRequestChanges: []types.RequestChanges{{
						RequestChange: &types.RequestChange{
							Action: types.ActionOfChangeResourceRecordSetsRequestForChangeResourceRecordSetsCreate,
							RequestResourceRecordSet: &types.RequestResourceRecordSet{
								Failover:      types.FailoverOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("PRIMARY"),
								Name:          nifcloud.String("test_name"),
								SetIdentifier: nifcloud.String("test_set_identifier"),
								TTL:           nifcloud.Int32(60),
								Type:          types.TypeOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("A"),
								Weight:        nifcloud.Int32(60),
								XniftyComment: nifcloud.String("test_comment"),
								ListOfRequestResourceRecords: []types.RequestResourceRecords{{
									RequestResourceRecord: &types.RequestResourceRecord{
										Value: nifcloud.String("192.0.2.1"),
									},
								}},
								RequestXniftyHealthCheckConfig: &types.RequestXniftyHealthCheckConfig{
									FullyQualifiedDomainName: nifcloud.String("test_resource_domain"),
									IPAddress:                nifcloud.String("192.0.2.1"),
									Port:                     nifcloud.Int32(8080),
									Protocol:                 types.ProtocolOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("HTTP"),
									ResourcePath:             nifcloud.String("test_resource_path"),
								},
							},
						},
					}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateChangeResourceRecordSetsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandListResourceRecordSets(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"zone_id": "test_zone_id",
		"type":    "A",
		"name":    "test_name",
	})
	rd.SetId("test_set_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *dns.ListResourceRecordSetsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &dns.ListResourceRecordSetsInput{
				Identifier: nifcloud.String("test_set_identifier"),
				Name:       nifcloud.String("test_name"),
				Type:       types.TypeOfListResourceRecordSetsRequest("A"),
				ZoneID:     nifcloud.String("test_zone_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandListResourceRecordSets(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteChangeResourceRecordSetsInput(t *testing.T) {
	healthCheck := map[string]interface{}{
		"protocol":        "HTTP",
		"ip_address":      "192.0.2.1",
		"port":            8080,
		"resource_path":   "test_resource_path",
		"resource_domain": "test_resource_domain",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"zone_id": "test_zone_id",
		"name":    "test_name",
		"type":    "A",
		"record":  "192.0.2.1",
		"ttl":     60,
		"weighted_routing_policy": []interface{}{map[string]interface{}{
			"weight": 60,
		}},
		"failover_routing_policy": []interface{}{map[string]interface{}{
			"type":         "PRIMARY",
			"health_check": []interface{}{healthCheck},
		}},
		"comment":        "test_comment",
		"set_identifier": "test_set_identifier",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *dns.ChangeResourceRecordSetsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &dns.ChangeResourceRecordSetsInput{
				ZoneID: nifcloud.String("test_zone_id"),
				RequestChangeBatch: &types.RequestChangeBatch{
					ListOfRequestChanges: []types.RequestChanges{{
						RequestChange: &types.RequestChange{
							Action: types.ActionOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("DELETE"),
							RequestResourceRecordSet: &types.RequestResourceRecordSet{
								Failover:      types.FailoverOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("PRIMARY"),
								Name:          nifcloud.String("test_name"),
								SetIdentifier: nifcloud.String("test_set_identifier"),
								TTL:           nifcloud.Int32(60),
								Type:          types.TypeOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("A"),
								Weight:        nifcloud.Int32(60),
								XniftyComment: nifcloud.String("test_comment"),
								ListOfRequestResourceRecords: []types.RequestResourceRecords{{
									RequestResourceRecord: &types.RequestResourceRecord{
										Value: nifcloud.String("192.0.2.1"),
									},
								}},
								RequestXniftyHealthCheckConfig: &types.RequestXniftyHealthCheckConfig{
									FullyQualifiedDomainName: nifcloud.String("test_resource_domain"),
									IPAddress:                nifcloud.String("192.0.2.1"),
									Port:                     nifcloud.Int32(8080),
									Protocol:                 types.ProtocolOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("HTTP"),
									ResourcePath:             nifcloud.String("test_resource_path"),
								},
							},
						},
					}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteChangeResourceRecordSetsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRequestResourceRecordSetInput(t *testing.T) {
	healthCheck := map[string]interface{}{
		"protocol":        "HTTP",
		"ip_address":      "192.0.2.1",
		"port":            8080,
		"resource_path":   "test_resource_path",
		"resource_domain": "test_resource_domain",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":   "test_name",
		"type":   "A",
		"record": "192.0.2.1",
		"ttl":    60,
		"weighted_routing_policy": []interface{}{map[string]interface{}{
			"weight": 60,
		}},
		"failover_routing_policy": []interface{}{map[string]interface{}{
			"type":         "PRIMARY",
			"health_check": []interface{}{healthCheck},
		}},
		"comment":        "test_comment",
		"set_identifier": "test_set_identifier",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *types.RequestResourceRecordSet
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &types.RequestResourceRecordSet{
				Failover:      types.FailoverOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("PRIMARY"),
				Name:          nifcloud.String("test_name"),
				SetIdentifier: nifcloud.String("test_set_identifier"),
				TTL:           nifcloud.Int32(60),
				Type:          types.TypeOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("A"),
				Weight:        nifcloud.Int32(60),
				XniftyComment: nifcloud.String("test_comment"),
				ListOfRequestResourceRecords: []types.RequestResourceRecords{{
					RequestResourceRecord: &types.RequestResourceRecord{
						Value: nifcloud.String("192.0.2.1"),
					},
				}},
				RequestXniftyHealthCheckConfig: &types.RequestXniftyHealthCheckConfig{
					FullyQualifiedDomainName: nifcloud.String("test_resource_domain"),
					IPAddress:                nifcloud.String("192.0.2.1"),
					Port:                     nifcloud.Int32(8080),
					Protocol:                 types.ProtocolOfChangeResourceRecordSetsRequestForChangeResourceRecordSets("HTTP"),
					ResourcePath:             nifcloud.String("test_resource_path"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRequestResourceRecordSetInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
