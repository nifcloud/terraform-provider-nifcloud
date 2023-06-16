package record

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns/types"
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
	raw := map[string]interface{}{
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
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), raw)
	rd.SetId("test_set_identifier")
	wantRd := schema.TestResourceDataRaw(t, newSchema(), raw)
	wantRd.SetId("test_set_identifier")

	rawWithAtSignAsName := map[string]interface{}{
		"zone_id": "test_zone_id",
		"name":    "@",
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
	}
	rdWithAtSignAsName := schema.TestResourceDataRaw(t, newSchema(), rawWithAtSignAsName)
	rdWithAtSignAsName.SetId("test_set_identifier")
	wantRdWithAtSignAsName := schema.TestResourceDataRaw(t, newSchema(), rawWithAtSignAsName)
	wantRdWithAtSignAsName.SetId("test_set_identifier")

	notFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *dns.ListResourceRecordSetsOutput
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
				res: &dns.ListResourceRecordSetsOutput{
					ResourceRecordSets: []types.ResourceRecordSets{
						{
							Failover:      nifcloud.String("PRIMARY"),
							Name:          nifcloud.String("test_name"),
							SetIdentifier: nifcloud.String("test_set_identifier"),
							TTL:           nifcloud.Int32(60),
							Type:          nifcloud.String("A"),
							Weight:        nifcloud.Int32(60),
							XniftyComment: nifcloud.String("test_comment"),
							ResourceRecords: []types.ResourceRecords{{
								Value: nifcloud.String("192.0.2.1"),
							}},
							XniftyHealthCheckConfig: &types.XniftyHealthCheckConfig{
								FullyQualifiedDomainName: nifcloud.String("test_resource_domain"),
								IPAddress:                nifcloud.String("192.0.2.1"),
								Port:                     nifcloud.Int32(8080),
								Protocol:                 nifcloud.String("HTTP"),
								ResourcePath:             nifcloud.String("test_resource_path"),
							},
						},
					},
				},
			},
			want: wantRd,
		},
		{
			name: "flattens the response whose name value is equal to the zone_id",
			args: args{
				d: rdWithAtSignAsName,
				res: &dns.ListResourceRecordSetsOutput{
					ResourceRecordSets: []types.ResourceRecordSets{
						{
							Failover:      nifcloud.String("PRIMARY"),
							Name:          nifcloud.String("test_zone_id"),
							SetIdentifier: nifcloud.String("test_set_identifier"),
							TTL:           nifcloud.Int32(60),
							Type:          nifcloud.String("A"),
							Weight:        nifcloud.Int32(60),
							XniftyComment: nifcloud.String("test_comment"),
							ResourceRecords: []types.ResourceRecords{{
								Value: nifcloud.String("192.0.2.1"),
							}},
							XniftyHealthCheckConfig: &types.XniftyHealthCheckConfig{
								FullyQualifiedDomainName: nifcloud.String("test_resource_domain"),
								IPAddress:                nifcloud.String("192.0.2.1"),
								Port:                     nifcloud.Int32(8080),
								Protocol:                 nifcloud.String("HTTP"),
								ResourcePath:             nifcloud.String("test_resource_path"),
							},
						},
					},
				},
			},
			want: wantRdWithAtSignAsName,
		},
		{
			name: "flattens the response when we requested with shorthand name",
			args: args{
				d: rd,
				res: &dns.ListResourceRecordSetsOutput{
					ResourceRecordSets: []types.ResourceRecordSets{
						{
							Failover:      nifcloud.String("PRIMARY"),
							Name:          nifcloud.String("test_name.test_zone_id"),
							SetIdentifier: nifcloud.String("test_set_identifier"),
							TTL:           nifcloud.Int32(60),
							Type:          nifcloud.String("A"),
							Weight:        nifcloud.Int32(60),
							XniftyComment: nifcloud.String("test_comment"),
							ResourceRecords: []types.ResourceRecords{{
								Value: nifcloud.String("192.0.2.1"),
							}},
							XniftyHealthCheckConfig: &types.XniftyHealthCheckConfig{
								FullyQualifiedDomainName: nifcloud.String("test_resource_domain"),
								IPAddress:                nifcloud.String("192.0.2.1"),
								Port:                     nifcloud.Int32(8080),
								Protocol:                 nifcloud.String("HTTP"),
								ResourcePath:             nifcloud.String("test_resource_path"),
							},
						},
					},
				},
			},
			want: wantRd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d: notFoundRd,
				res: &dns.ListResourceRecordSetsOutput{
					ResourceRecordSets: []types.ResourceRecordSets{},
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
