package record

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
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
		"default_host":   "test_default_host",
		"comment":        "test_comment",
		"set_identifier": "test_set_identifier",
	})
	rd.SetId("test_set_identifier")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *dns.ListResourceRecordSetsResponse
		d   *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want *schema.ResourceData
	}{
		{
			name: "flattens the response",
			args: args{
				d: rd,
				res: &dns.ListResourceRecordSetsResponse{
					ListResourceRecordSetsOutput: &dns.ListResourceRecordSetsOutput{
						ResourceRecordSets: []dns.ResourceRecordSets{
							{
								Failover:          nifcloud.String("PRIMARY"),
								Name:              nifcloud.String("test_name"),
								SetIdentifier:     nifcloud.String("test_set_identifier"),
								TTL:               nifcloud.Int64(60),
								Type:              nifcloud.String("A"),
								Weight:            nifcloud.Int64(60),
								XniftyComment:     nifcloud.String("test_comment"),
								XniftyDefaultHost: nifcloud.String("test_default_host"),
								ResourceRecords: []dns.ResourceRecords{{
									Value: nifcloud.String("192.0.2.1"),
								}},
								XniftyHealthCheckConfig: &dns.XniftyHealthCheckConfig{
									FullyQualifiedDomainName: nifcloud.String("test_resource_domain"),
									IPAddress:                nifcloud.String("192.0.2.1"),
									Port:                     nifcloud.Int64(8080),
									Protocol:                 nifcloud.String("HTTP"),
									ResourcePath:             nifcloud.String("test_resource_path"),
								},
							},
						},
					},
				},
			},
			want: rd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d: wantNotFoundRd,
				res: &dns.ListResourceRecordSetsResponse{
					ListResourceRecordSetsOutput: &dns.ListResourceRecordSetsOutput{
						ResourceRecordSets: []dns.ResourceRecordSets{},
					},
				},
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			assert.NoError(t, err)

			if tt.args.res == nil {
				return
			}

			wantState := tt.want.State()
			if wantState == nil {
				tt.want.SetId("some")
				wantState = tt.want.State()
			}

			gotState := tt.args.d.State()
			if gotState == nil {
				tt.args.d.SetId("some")
				gotState = tt.args.d.State()
			}

			assert.Equal(t, wantState.Attributes, gotState.Attributes)
		})
	}
}
