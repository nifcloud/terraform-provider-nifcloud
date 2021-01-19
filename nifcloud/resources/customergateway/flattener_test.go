package customergateway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"customer_gateway_id":                "test_customer_gateway_id",
		"nifty_customer_gateway_name":        "test_nifty_customer_gateway_name",
		"ip_address":                         "test_ip_address",
		"nifty_lan_side_ip_address":          "test_nifty_lan_side_ip_address",
		"nifty_lan_side_cidr_block":          "test_nifty_lan_side_cidr_block",
		"nifty_customer_gateway_description": "test_nifty_customer_gateway_description",
	})
	rd.SetId("test_customer_gateway_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeCustomerGatewaysResponse
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
				res: &computing.DescribeCustomerGatewaysResponse{
					DescribeCustomerGatewaysOutput: &computing.DescribeCustomerGatewaysOutput{
						CustomerGatewaySet: []computing.CustomerGatewaySet{
							{
								CustomerGatewayId:               nifcloud.String("test_customer_gateway_id"),
								NiftyCustomerGatewayName:        nifcloud.String("test_nifty_customer_gateway_name"),
								IpAddress:                       nifcloud.String("test_ip_address"),
								NiftyLanSideIpAddress:           nifcloud.String("test_nifty_lan_side_ip_address"),
								NiftyLanSideCidrBlock:           nifcloud.String("test_nifty_lan_side_cidr_block"),
								NiftyCustomerGatewayDescription: nifcloud.String("test_nifty_customer_gateway_description"),
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
				res: &computing.DescribeCustomerGatewaysResponse{
					DescribeCustomerGatewaysOutput: &computing.DescribeCustomerGatewaysOutput{
						CustomerGatewaySet: []computing.CustomerGatewaySet{},
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
