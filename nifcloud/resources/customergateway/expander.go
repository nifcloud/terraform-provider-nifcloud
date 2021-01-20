package customergateway

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

const (
	// Available types.
	// More info: https://pfs.nifcloud.com/api/rest/CreateCustomerGateway.htm

	// TypeIpsec represents IPsec
	typeIpsec = "IPsec"
	// TypeIpsecVti represents IPsec VTI
	typeIpsecVti = "IPsec VTI"
	// Type represents L2TPv3 / IPsec
	typeL2tpv3Ipsec = "L2TPv3 / IPsec"
)

var (
	// typeMapping converts string to customer gateway type.
	// More info: https://pfs.nifcloud.com/api/rest/CreateCustomerGateway.htm
	typeMapping = map[string]computing.TypeOfCreateCustomerGatewayRequest{
		typeIpsec:       computing.TypeOfCreateCustomerGatewayRequestIpsec,
		typeIpsecVti:    computing.TypeOfCreateCustomerGatewayRequestIpsecVti,
		typeL2tpv3Ipsec: computing.TypeOfCreateCustomerGatewayRequestL2tpv3Ipsec,
	}
)

func expandCreateCustomerGatewayInput(d *schema.ResourceData) *computing.CreateCustomerGatewayInput {
	return &computing.CreateCustomerGatewayInput{
		IpAddress:                       nifcloud.String(d.Get("ip_address").(string)),
		NiftyCustomerGatewayDescription: nifcloud.String(d.Get("nifty_customer_gateway_description").(string)),
		NiftyCustomerGatewayName:        nifcloud.String(d.Get("nifty_customer_gateway_name").(string)),
		NiftyLanSideCidrBlock:           nifcloud.String(d.Get("nifty_lan_side_cidr_block").(string)),
		NiftyLanSideIpAddress:           nifcloud.String(d.Get("nifty_lan_side_ip_address").(string)),
		Type:                            typeMapping[d.Get("type").(string)],
	}
}

func expandDescribeCustomerGatewaysInput(d *schema.ResourceData) *computing.DescribeCustomerGatewaysInput {
	return &computing.DescribeCustomerGatewaysInput{
		CustomerGatewayId: []string{d.Id()},
	}
}

func expandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayName(d *schema.ResourceData) *computing.NiftyModifyCustomerGatewayAttributeInput {
	return &computing.NiftyModifyCustomerGatewayAttributeInput{
		CustomerGatewayId: nifcloud.String(d.Id()),
		Attribute:         computing.AttributeOfNiftyModifyCustomerGatewayAttributeRequestNiftyCustomerGatewayName,
		Value:             nifcloud.String(d.Get("nifty_customer_gateway_name").(string)),
	}
}

func expandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayDescription(d *schema.ResourceData) *computing.NiftyModifyCustomerGatewayAttributeInput {
	return &computing.NiftyModifyCustomerGatewayAttributeInput{
		CustomerGatewayId: nifcloud.String(d.Id()),
		Attribute:         computing.AttributeOfNiftyModifyCustomerGatewayAttributeRequestNiftyCustomerGatewayDescription,
		Value:             nifcloud.String(d.Get("nifty_customer_gateway_description").(string)),
	}
}

func expandDeleteCustomerGatewayInput(d *schema.ResourceData) *computing.DeleteCustomerGatewayInput {
	return &computing.DeleteCustomerGatewayInput{
		CustomerGatewayId: nifcloud.String(d.Id()),
	}
}
