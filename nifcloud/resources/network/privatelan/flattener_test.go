package privatelan

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":   "test_accounting_type",
		"availability_zone": "test_availability_zone",
		"cidr_block":        "test_cidr_block",
		"description":       "test_description",
		"private_lan_name":  "test_private_lan_name",
		"state":             "test_state",
		"network_id":        "test_network_id",
	})
	rd.SetId("test_network_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.NiftyDescribePrivateLansOutput
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
				res: &computing.NiftyDescribePrivateLansOutput{
					PrivateLanSet: []types.PrivateLanSet{
						{
							AvailabilityZone:        nifcloud.String("test_availability_zone"),
							CidrBlock:               nifcloud.String("test_cidr_block"),
							Description:             nifcloud.String("test_description"),
							PrivateLanName:          nifcloud.String("test_network_id"),
							State:                   nifcloud.String("test_state"),
							NetworkId:               nifcloud.String("test_network_id"),
							NextMonthAccountingType: nifcloud.String("test_accounting_type"),
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
				res: &computing.NiftyDescribePrivateLansOutput{
					PrivateLanSet: []types.PrivateLanSet{},
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
