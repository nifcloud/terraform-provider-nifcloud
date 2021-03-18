package routetable

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateRouteInput(t *testing.T) {
	route := map[string]interface{}{
		"network_id":     "test_network_id",
		"network_name":   "test_network_name",
		"ip_address":     "test_ip_address",
		"cidr_block":     "test_cidr_block",
		"route_table_id": "test_route_table_id",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route": []interface{}{route},
	})
	rd.SetId("test_route_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateRouteInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateRouteInput{
				RouteTableId:         nifcloud.String("test_route_table_id"),
				DestinationCidrBlock: nifcloud.String("test_cidr_block"),
				NetworkId:            nifcloud.String("test_network_id"),
				NetworkName:          nifcloud.String("test_network_name"),
				IpAddress:            nifcloud.String("test_ip_address"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateRouteInput(tt.args, route)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestExpandDeleteRouteInput(t *testing.T) {
	route := map[string]interface{}{
		"network_id":     "test_network_id",
		"network_name":   "test_network_name",
		"ip_address":     "test_ip_address",
		"cidr_block":     "test_cidr_block",
		"route_table_id": "test_route_table_id",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route": []interface{}{route},
	})
	rd.SetId("test_route_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteRouteInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteRouteInput{
				RouteTableId:         nifcloud.String("test_route_table_id"),
				DestinationCidrBlock: nifcloud.String("test_cidr_block"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteRouteInput(tt.args, route)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteRouteTableInput(t *testing.T) {
	route := map[string]interface{}{
		"network_id":     "test_network_id",
		"network_name":   "test_network_name",
		"ip_address":     "test_ip_address",
		"cidr_block":     "test_cidr_block",
		"route_table_id": "test_route_table_id",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route": []interface{}{route},
	})
	rd.SetId("test_route_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteRouteTableInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteRouteTableInput{
				RouteTableId: nifcloud.String("test_route_table_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteRouteTableInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeRouteTablesInput(t *testing.T) {
	route := map[string]interface{}{
		"network_id":     "test_network_id",
		"network_name":   "test_network_name",
		"ip_address":     "test_ip_address",
		"cidr_block":     "test_cidr_block",
		"route_table_id": "test_route_table_id",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route": []interface{}{route},
	})
	rd.SetId("test_route_table_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeRouteTablesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeRouteTablesInput{
				RouteTableId: []string{"test_route_table_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeRouteTablesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
