package devopsfirewallgroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	wantRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_name",
		"availability_zone": "test_zone",
		"description":       "test_description",
		"rule": []interface{}{
			map[string]interface{}{
				"id":          "test_id_01",
				"protocol":    "TCP",
				"port":        443,
				"cidr_ip":     "172.16.0.0/24",
				"description": "test_description_01",
			},
			map[string]interface{}{
				"id":          "test_id_02",
				"protocol":    "ICMP",
				"port":        0,
				"cidr_ip":     "172.16.0.0/24",
				"description": "test_description_02",
			},
		},
	})
	wantRd.SetId("test_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *devops.GetFirewallGroupOutput
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
				res: &devops.GetFirewallGroupOutput{
					FirewallGroup: &types.FirewallGroup{
						FirewallGroupName: nifcloud.String("test_name"),
						AvailabilityZone:  nifcloud.String("test_zone"),
						Description:       nifcloud.String("test_description"),
						Rules: []types.Rules{
							{
								Id:          nifcloud.String("test_id_01"),
								Protocol:    nifcloud.String("TCP"),
								Port:        nifcloud.Int32(int32(443)),
								CidrIp:      nifcloud.String("172.16.0.0/24"),
								Description: nifcloud.String("test_description_01"),
							},
							{
								Id:          nifcloud.String("test_id_02"),
								Protocol:    nifcloud.String("ICMP"),
								Port:        nil,
								CidrIp:      nifcloud.String("172.16.0.0/24"),
								Description: nifcloud.String("test_description_02"),
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
				d:   wantNotFoundRd,
				res: nil,
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			assert.NoError(t, err)

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
