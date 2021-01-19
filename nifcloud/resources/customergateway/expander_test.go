package customergateway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateCustomerGatwayInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"nifty_customer_gateway_name":        "test_nifty_customer_gateway_name",
		"ip_address":                         "test_ip_address",
		"nifty_lan_side_ip_address":          "test_nifty_lan_side_ip_address",
		"nifty_lan_side_cidr_block":          "test_nifty_lan_side_cidr_block",
		"nifty_customer_gateway_description": "test_nifty_customer_gateway_description",
		"bgp_asn":                            36000,
	})
	rd.SetId("test_customer_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateCustomerGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateCustomerGatewayInput{
				NiftyCustomerGatewayName:        nifcloud.String("test_nifty_customer_gateway_name"),
				IpAddress:                       nifcloud.String("test_ip_address"),
				NiftyLanSideIpAddress:           nifcloud.String("test_nifty_lan_side_ip_address"),
				NiftyLanSideCidrBlock:           nifcloud.String("test_nifty_lan_side_cidr_block"),
				NiftyCustomerGatewayDescription: nifcloud.String("test_nifty_customer_gateway_description"),
				BgpAsn:                          nifcloud.Int64(int64(36000)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateCustomerGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeCustomerGatewaysInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_customer_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeCustomerGatewaysInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeCustomerGatewaysInput{
				CustomerGatewayId: []string{"test_customer_gateway_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeCustomerGatewaysInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"nifty_customer_gateway_name": "test_nifty_customer_gateway_name",
	})
	rd.SetId("test_customer_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyCustomerGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyCustomerGatewayAttributeInput{
				CustomerGatewayId: nifcloud.String("test_customer_gateway_id"),
				Attribute:         computing.AttributeOfNiftyModifyCustomerGatewayAttributeRequestNiftyCustomerGatewayName,
				Value:             nifcloud.String("test_nifty_customer_gateway_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"nifty_customer_gateway_description": "test_nifty_customer_gateway_description",
	})
	rd.SetId("test_customer_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyCustomerGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyCustomerGatewayAttributeInput{
				CustomerGatewayId: nifcloud.String("test_customer_gateway_id"),
				Attribute:         computing.AttributeOfNiftyModifyCustomerGatewayAttributeRequestNiftyCustomerGatewayDescription,
				Value:             nifcloud.String("test_nifty_customer_gateway_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteCustomerGatewayInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_customer_gateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteCustomerGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteCustomerGatewayInput{
				CustomerGatewayId: nifcloud.String("test_customer_gateway_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteCustomerGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
