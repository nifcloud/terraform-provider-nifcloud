package routetable

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandDescribeRouteTablesInput(d *schema.ResourceData) *computing.DescribeRouteTablesInput {
	return &computing.DescribeRouteTablesInput{
		RouteTableId: []string{d.Id()},
	}
}

func expandDeleteRouteTableInput(d *schema.ResourceData) *computing.DeleteRouteTableInput {
	return &computing.DeleteRouteTableInput{
		RouteTableId: nifcloud.String(d.Id()),
	}
}

func expandDeleteRouteInput(d *schema.ResourceData, route map[string]interface{}) *computing.DeleteRouteInput {
	return &computing.DeleteRouteInput{
		RouteTableId:         nifcloud.String(d.Id()),
		DestinationCidrBlock: nifcloud.String(route["cidr_block"].(string)),
	}
}

func expandCreateRouteInput(d *schema.ResourceData, route map[string]interface{}) *computing.CreateRouteInput {
	return &computing.CreateRouteInput{
		RouteTableId:         nifcloud.String(d.Id()),
		DestinationCidrBlock: nifcloud.String(route["cidr_block"].(string)),
		IpAddress:            nifcloud.String(route["ip_address"].(string)),
		NetworkId:            nifcloud.String(route["network_id"].(string)),
		NetworkName:          nifcloud.String(route["network_name"].(string)),
	}
}
