package firewallgroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"description": "test_description",
		"rule": []interface{}{map[string]interface{}{
			"protocol": "ANY",
			"cidr_ip":  "0.0.0.0/0",
			"id":       "test_rule_id_01",
		}},
	})
	rd.SetId("test_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *hatoba.GetFirewallGroupResponse
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
				res: &hatoba.GetFirewallGroupResponse{
					GetFirewallGroupOutput: &hatoba.GetFirewallGroupOutput{
						FirewallGroup: &hatoba.FirewallGroupResponse{
							Name:        nifcloud.String("test_name"),
							Description: nifcloud.String("test_description"),
							Rules: []hatoba.FirewallRule{
								{
									Protocol: nifcloud.String("ANY"),
									CidrIp:   nifcloud.String("0.0.0.0/0"),
									Id:       nifcloud.String("test_rule_id_01"),
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
				res: &hatoba.GetFirewallGroupResponse{
					GetFirewallGroupOutput: &hatoba.GetFirewallGroupOutput{
						FirewallGroup: &hatoba.FirewallGroupResponse{},
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
