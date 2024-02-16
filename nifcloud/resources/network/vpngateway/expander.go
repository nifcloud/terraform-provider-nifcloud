package vpngateway

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandCreateVpnGatewayInput(d *schema.ResourceData) *computing.CreateVpnGatewayInput {

	var niftyNetwork *types.RequestNiftyNetwork
	if row, ok := d.GetOk("network_name"); ok {
		niftyNetwork = &types.RequestNiftyNetwork{
			NetworkName: nifcloud.String(row.(string)),
			IpAddress:   nifcloud.String(d.Get("ip_address").(string)),
		}
	}
	if row, ok := d.GetOk("network_id"); ok {
		niftyNetwork = &types.RequestNiftyNetwork{
			NetworkId: nifcloud.String(row.(string)),
			IpAddress: nifcloud.String(d.Get("ip_address").(string)),
		}
	}

	var securityGroup []string
	if row, ok := d.GetOk("security_group"); ok {
		securityGroup = append(securityGroup, row.(string))
	}

	placement := &types.RequestPlacementOfCreateVpnGateway{
		AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
	}

	input := &computing.CreateVpnGatewayInput{
		AccountingType:             types.AccountingTypeOfCreateVpnGatewayRequest(d.Get("accounting_type").(string)),
		NiftyVpnGatewayDescription: nifcloud.String(d.Get("description").(string)),
		NiftyVpnGatewayName:        nifcloud.String(d.Get("name").(string)),
		NiftyVpnGatewayType:        types.NiftyVpnGatewayTypeOfCreateVpnGatewayRequest(d.Get("type").(string)),
		Placement:                  placement,
		NiftyNetwork:               niftyNetwork,
		SecurityGroup:              securityGroup,
	}
	return input
}

func expandDescribeVpnGatewaysInput(d *schema.ResourceData) *computing.DescribeVpnGatewaysInput {
	return &computing.DescribeVpnGatewaysInput{
		VpnGatewayId: []string{d.Id()},
	}
}

func expandDeleteVpnGatewayInput(d *schema.ResourceData) *computing.DeleteVpnGatewayInput {
	return &computing.DeleteVpnGatewayInput{
		VpnGatewayId: nifcloud.String(d.Id()),
	}
}

func expandNiftyModifyVpnGatewayAttributeInputForAccountingType(d *schema.ResourceData) *computing.NiftyModifyVpnGatewayAttributeInput {
	return &computing.NiftyModifyVpnGatewayAttributeInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		Attribute:    types.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayAccountingType,
		Value:        nifcloud.String(d.Get("accounting_type").(string)),
	}
}

func expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayDescription(d *schema.ResourceData) *computing.NiftyModifyVpnGatewayAttributeInput {
	return &computing.NiftyModifyVpnGatewayAttributeInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		Attribute:    types.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayDescription,
		Value:        nifcloud.String(d.Get("description").(string)),
	}
}

func expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayName(d *schema.ResourceData) *computing.NiftyModifyVpnGatewayAttributeInput {
	return &computing.NiftyModifyVpnGatewayAttributeInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		Attribute:    types.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayName,
		Value:        nifcloud.String(d.Get("name").(string)),
	}
}

func expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayType(d *schema.ResourceData) *computing.NiftyModifyVpnGatewayAttributeInput {
	return &computing.NiftyModifyVpnGatewayAttributeInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		Attribute:    types.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayType,
		Value:        nifcloud.String(d.Get("type").(string)),
	}
}

func expandNiftyUpdateVpnGatewayNetworkInterfacesInput(d *schema.ResourceData) *computing.NiftyUpdateVpnGatewayNetworkInterfacesInput {
	networkInterface := types.RequestNetworkInterfaceOfNiftyUpdateVpnGatewayNetworkInterfaces{
		IpAddress:        nifcloud.String(d.Get("ip_address").(string)),
		IsOutsideNetwork: nifcloud.Bool(false),
	}
	if row, ok := d.GetOk("network_name"); ok {
		networkInterface.NetworkName = nifcloud.String(row.(string))
	}
	if row, ok := d.GetOk("network_id"); ok {
		networkInterface.NetworkId = nifcloud.String(row.(string))
	}

	return &computing.NiftyUpdateVpnGatewayNetworkInterfacesInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		NetworkInterface: []types.RequestNetworkInterfaceOfNiftyUpdateVpnGatewayNetworkInterfaces{
			{
				NetworkId:        nifcloud.String("net-COMMON_GLOBAL"),
				IsOutsideNetwork: nifcloud.Bool(true),
			},
			networkInterface,
		},
	}
}

func expandNiftyModifyVpnGatewayAttributeInputForSecurityGroup(d *schema.ResourceData) *computing.NiftyModifyVpnGatewayAttributeInput {
	return &computing.NiftyModifyVpnGatewayAttributeInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		Attribute:    types.AttributeOfNiftyModifyVpnGatewayAttributeRequestGroupId,
		Value:        nifcloud.String(d.Get("security_group").(string)),
	}
}

func expandNiftyAssociateRouteTableWithVpnGatewayInput(d *schema.ResourceData) *computing.NiftyAssociateRouteTableWithVpnGatewayInput {
	return &computing.NiftyAssociateRouteTableWithVpnGatewayInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		RouteTableId: nifcloud.String(d.Get("route_table_id").(string)),
	}
}

func expandNiftyDisassociateRouteTableFromVpnGatewayInput(d *schema.ResourceData) *computing.NiftyDisassociateRouteTableFromVpnGatewayInput {
	return &computing.NiftyDisassociateRouteTableFromVpnGatewayInput{
		AssociationId: nifcloud.String(d.Get("route_table_association_id").(string)),
	}
}
func expandNiftyReplaceRouteTableAssociationWithVpnGatewayInput(d *schema.ResourceData) *computing.NiftyReplaceRouteTableAssociationWithVpnGatewayInput {
	return &computing.NiftyReplaceRouteTableAssociationWithVpnGatewayInput{
		AssociationId: nifcloud.String(d.Get("route_table_association_id").(string)),
		RouteTableId:  nifcloud.String(d.Get("route_table_id").(string)),
	}
}
