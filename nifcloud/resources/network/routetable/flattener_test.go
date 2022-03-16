package routetable

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
		"route_table_id": "test_route_table_id",
		"route": []interface{}{map[string]interface{}{
			"cidr_block":   "test_cidr_block",
			"network_id":   "test_network_id",
			"network_name": "test_network_name",
			"ip_address":   "test_ip_address",
		}},
	})
	rd.SetId("test_route_table_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeRouteTablesOutput
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
				res: &computing.DescribeRouteTablesOutput{
					RouteTableSet: []types.RouteTableSet{
						{
							RouteTableId: nifcloud.String("test_route_table_id"),
							RouteSet: []types.RouteSet{
								{
									NetworkId:            nifcloud.String("test_network_id"),
									NetworkName:          nifcloud.String("test_network_name"),
									IpAddress:            nifcloud.String("test_ip_address"),
									DestinationCidrBlock: nifcloud.String("test_cidr_block"),
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
				d:   wantNotFoundRd,
				res: &computing.DescribeRouteTablesOutput{},
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
