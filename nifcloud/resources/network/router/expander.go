package router

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandNiftyCreateRouterInput(d *schema.ResourceData) *computing.NiftyCreateRouterInput {
	var networkInterface []computing.RequestNetworkInterfaceOfNiftyCreateRouter
	for _, ni := range d.Get("network_interface").(*schema.Set).List() {
		if v, ok := ni.(map[string]interface{}); ok {
			n := computing.RequestNetworkInterfaceOfNiftyCreateRouter{}
			if row, ok := v["network_id"]; ok {
				n.NetworkId = nifcloud.String(row.(string))
			}
			if row, ok := v["network_name"]; ok {
				n.NetworkName = nifcloud.String(row.(string))
			}
			if row, ok := v["ip_address"]; ok {
				n.IpAddress = nifcloud.String(row.(string))
			}
			if row, ok := v["dhcp"]; ok {
				n.Dhcp = nifcloud.Bool(row.(bool))
			}
			if row, ok := v["dhcp_options_id"]; ok {
				n.DhcpOptionsId = nifcloud.String(row.(string))
			}
			if row, ok := v["dhcp_config_id"]; ok {
				n.DhcpConfigId = nifcloud.String(row.(string))
			}
			networkInterface = append(networkInterface, n)
		}
	}

	var securityGroup []string
	if row, ok := d.GetOk("security_group"); ok {
		securityGroup = append(securityGroup, row.(string))
	}

	input := &computing.NiftyCreateRouterInput{
		RouterName:       nifcloud.String(d.Get("name").(string)),
		SecurityGroup:    securityGroup,
		Type:             computing.TypeOfNiftyCreateRouterRequest(d.Get("type").(string)),
		AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		AccountingType:   computing.AccountingTypeOfNiftyCreateRouterRequest(d.Get("accounting_type").(string)),
		Description:      nifcloud.String(d.Get("description").(string)),
		NetworkInterface: networkInterface,
	}

	return input
}

func expandAssociateRouteTableInput(d *schema.ResourceData) *computing.AssociateRouteTableInput {
	return &computing.AssociateRouteTableInput{
		RouteTableId: nifcloud.String(d.Get("route_table_id").(string)),
		RouterId:     nifcloud.String(d.Id()),
	}
}

func expandNiftyAssociateNatTableInput(d *schema.ResourceData) *computing.NiftyAssociateNatTableInput {
	return &computing.NiftyAssociateNatTableInput{
		NatTableId: nifcloud.String(d.Get("nat_table_id").(string)),
		RouterId:   nifcloud.String(d.Id()),
	}
}

func expandNiftyDescribeRoutersInput(d *schema.ResourceData) *computing.NiftyDescribeRoutersInput {
	return &computing.NiftyDescribeRoutersInput{
		RouterId: []string{d.Id()},
	}
}

func expandNiftyModifyRouterAttributeInputForRouterName(d *schema.ResourceData) *computing.NiftyModifyRouterAttributeInput {
	return &computing.NiftyModifyRouterAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestRouterName,
		Value:     nifcloud.String(d.Get("name").(string)),
	}
}

func expandNiftyModifyRouterAttributeInputForAccountingType(d *schema.ResourceData) *computing.NiftyModifyRouterAttributeInput {
	return &computing.NiftyModifyRouterAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestAccountingType,
		Value:     nifcloud.String(d.Get("accounting_type").(string)),
	}
}

func expandNiftyModifyRouterAttributeInputForDescription(d *schema.ResourceData) *computing.NiftyModifyRouterAttributeInput {
	return &computing.NiftyModifyRouterAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestDescription,
		Value:     nifcloud.String(d.Get("description").(string)),
	}
}

func expandNiftyModifyRouterAttributeInputForType(d *schema.ResourceData) *computing.NiftyModifyRouterAttributeInput {
	return &computing.NiftyModifyRouterAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestType,
		Value:     nifcloud.String(d.Get("type").(string)),
	}
}

func expandNiftyModifyRouterAttributeInputForSecurityGroup(d *schema.ResourceData) *computing.NiftyModifyRouterAttributeInput {
	return &computing.NiftyModifyRouterAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestGroupId,
		Value:     nifcloud.String(d.Get("security_group").(string)),
	}
}

func expandNiftyDeregisterRoutersFromSecurityGroupInput(d *schema.ResourceData) *computing.NiftyDeregisterRoutersFromSecurityGroupInput {
	securityGroup, _ := d.GetChange("security_group")
	return &computing.NiftyDeregisterRoutersFromSecurityGroupInput{
		GroupName: nifcloud.String(securityGroup.(string)),
		RouterSet: []computing.RequestRouterSetOfNiftyDeregisterRoutersFromSecurityGroup{
			{
				RouterId: nifcloud.String(d.Id()),
			},
		},
	}
}

func expandNiftyUpdateRouterNetworkInterfacesInput(d *schema.ResourceData) *computing.NiftyUpdateRouterNetworkInterfacesInput {
	var networkInterface []computing.RequestNetworkInterfaceOfNiftyUpdateRouterNetworkInterfaces
	for _, ni := range d.Get("network_interface").(*schema.Set).List() {
		if v, ok := ni.(map[string]interface{}); ok {
			n := computing.RequestNetworkInterfaceOfNiftyUpdateRouterNetworkInterfaces{}
			if row, ok := v["network_id"]; ok {
				n.NetworkId = nifcloud.String(row.(string))
			}
			if row, ok := v["network_name"]; ok {
				n.NetworkName = nifcloud.String(row.(string))
			}
			if row, ok := v["ip_address"]; ok {
				n.IpAddress = nifcloud.String(row.(string))
			}
			if row, ok := v["dhcp"]; ok {
				n.Dhcp = nifcloud.Bool(row.(bool))
			}
			if row, ok := v["dhcp_options_id"]; ok {
				n.DhcpOptionsId = nifcloud.String(row.(string))
			}
			if row, ok := v["dhcp_config_id"]; ok {
				n.DhcpConfigId = nifcloud.String(row.(string))
			}
			networkInterface = append(networkInterface, n)
		}
	}

	return &computing.NiftyUpdateRouterNetworkInterfacesInput{
		RouterId:         nifcloud.String(d.Id()),
		NetworkInterface: networkInterface,
	}
}

func expandNiftyDisassociateNatTableInput(d *schema.ResourceData) *computing.NiftyDisassociateNatTableInput {
	return &computing.NiftyDisassociateNatTableInput{
		AssociationId: nifcloud.String(d.Get("nat_table_association_id").(string)),
	}
}

func expandNiftyReplaceNatTableAssociationInput(d *schema.ResourceData) *computing.NiftyReplaceNatTableAssociationInput {
	return &computing.NiftyReplaceNatTableAssociationInput{
		AssociationId: nifcloud.String(d.Get("nat_table_association_id").(string)),
		NatTableId:    nifcloud.String(d.Get("nat_table_id").(string)),
	}
}

func expandDisassociateRouteTableInput(d *schema.ResourceData) *computing.DisassociateRouteTableInput {
	return &computing.DisassociateRouteTableInput{
		AssociationId: nifcloud.String(d.Get("route_table_association_id").(string)),
	}
}

func expandReplaceRouteTableAssociation(d *schema.ResourceData) *computing.ReplaceRouteTableAssociationInput {
	return &computing.ReplaceRouteTableAssociationInput{
		AssociationId: nifcloud.String(d.Get("route_table_association_id").(string)),
		RouteTableId:  nifcloud.String(d.Get("route_table_id").(string)),
	}
}

func expandNiftyDeleteRouterInput(d *schema.ResourceData) *computing.NiftyDeleteRouterInput {
	return &computing.NiftyDeleteRouterInput{
		RouterId: nifcloud.String(d.Id()),
	}
}
