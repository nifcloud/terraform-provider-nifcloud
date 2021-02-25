package vpngateway

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandCreateVpnGatewayInput(d *schema.ResourceData) *computing.CreateVpnGatewayInput {

	var niftyNetwork *computing.RequestNiftyNetwork
	if row, ok := d.GetOk("network_name"); ok {
		niftyNetwork = &computing.RequestNiftyNetwork{
			NetworkName: nifcloud.String(row.(string)),
			IpAddress:   nifcloud.String(d.Get("ip_address").(string)),
		}
	}
	if row, ok := d.GetOk("network_id"); ok {
		niftyNetwork = &computing.RequestNiftyNetwork{
			NetworkId: nifcloud.String(row.(string)),
			IpAddress: nifcloud.String(d.Get("ip_address").(string)),
		}
	}

	var securityGroup []string
	if row, ok := d.GetOk("security_group"); ok {
		securityGroup = append(securityGroup, row.(string))
	}

	placement := &computing.RequestPlacementOfCreateVpnGateway{
		AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
	}

	input := &computing.CreateVpnGatewayInput{
		AccountingType:             computing.AccountingTypeOfCreateVpnGatewayRequest(d.Get("accounting_type").(string)),
		NiftyVpnGatewayDescription: nifcloud.String(d.Get("description").(string)),
		NiftyVpnGatewayName:        nifcloud.String(d.Get("name").(string)),
		NiftyVpnGatewayType:        computing.NiftyVpnGatewayTypeOfCreateVpnGatewayRequest(d.Get("type").(string)),
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
		Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayAccountingType,
		Value:        computing.ValueOfNiftyModifyVpnGatewayAttributeRequest(d.Get("accounting_type").(string)),
	}
}

func expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayDescription(d *schema.ResourceData) *computing.NiftyModifyVpnGatewayAttributeInput {
	return &computing.NiftyModifyVpnGatewayAttributeInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayDescription,
		Value:        computing.ValueOfNiftyModifyVpnGatewayAttributeRequest(d.Get("description").(string)),
	}
}

func expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayName(d *schema.ResourceData) *computing.NiftyModifyVpnGatewayAttributeInput {
	return &computing.NiftyModifyVpnGatewayAttributeInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayName,
		Value:        computing.ValueOfNiftyModifyVpnGatewayAttributeRequest(d.Get("name").(string)),
	}
}

func expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayType(d *schema.ResourceData) *computing.NiftyModifyVpnGatewayAttributeInput {
	return &computing.NiftyModifyVpnGatewayAttributeInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayType,
		Value:        computing.ValueOfNiftyModifyVpnGatewayAttributeRequest(d.Get("type").(string)),
	}
}

func expandNiftyUpdateVpnGatewayNetworkInterfacesInput(d *schema.ResourceData) *computing.NiftyUpdateVpnGatewayNetworkInterfacesInput {
	return &computing.NiftyUpdateVpnGatewayNetworkInterfacesInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		NetworkInterface: &computing.RequestNetworkInterfaceOfNiftyUpdateVpnGatewayNetworkInterfaces{
			IpAddress: nifcloud.String(d.Get("ip_address").(string)),
		},
	}
}

func expandNiftyModifyVpnGatewayAttributeInputForSecurityGroup(d *schema.ResourceData) *computing.NiftyModifyVpnGatewayAttributeInput {
	return &computing.NiftyModifyVpnGatewayAttributeInput{
		VpnGatewayId: nifcloud.String(d.Id()),
		Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestGroupId,
		Value:        computing.ValueOfNiftyModifyVpnGatewayAttributeRequest(d.Get("security_group").(string)),
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
