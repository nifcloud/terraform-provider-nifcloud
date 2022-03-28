package customergateway

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

const (
	// AvailableWaiter(svc).Wait types.
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
	typeMapping = map[string]types.TypeOfCreateCustomerGatewayRequest{
		typeIpsec:       types.TypeOfCreateCustomerGatewayRequestIpsec,
		typeIpsecVti:    types.TypeOfCreateCustomerGatewayRequestIpsecVti,
		typeL2tpv3Ipsec: types.TypeOfCreateCustomerGatewayRequestL2tpv3Ipsec,
	}
)

func expandCreateCustomerGatewayInput(d *schema.ResourceData) *computing.CreateCustomerGatewayInput {
	return &computing.CreateCustomerGatewayInput{
		NiftyCustomerGatewayDescription: nifcloud.String(d.Get("description").(string)),
		NiftyCustomerGatewayName:        nifcloud.String(d.Get("name").(string)),
		Type:                            typeMapping[d.Get("type").(string)],
		IpAddress:                       nifcloud.String(d.Get("ip_address").(string)),
		NiftyLanSideCidrBlock:           nifcloud.String(d.Get("lan_side_cidr_block").(string)),
		NiftyLanSideIpAddress:           nifcloud.String(d.Get("lan_side_ip_address").(string)),
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
		Attribute:         types.AttributeOfNiftyModifyCustomerGatewayAttributeRequestNiftyCustomerGatewayName,
		Value:             nifcloud.String(d.Get("name").(string)),
	}
}

func expandNiftyModifyCustomerGatewayAttributeInputForNiftyCustomerGatewayDescription(d *schema.ResourceData) *computing.NiftyModifyCustomerGatewayAttributeInput {
	return &computing.NiftyModifyCustomerGatewayAttributeInput{
		CustomerGatewayId: nifcloud.String(d.Id()),
		Attribute:         types.AttributeOfNiftyModifyCustomerGatewayAttributeRequestNiftyCustomerGatewayDescription,
		Value:             nifcloud.String(d.Get("description").(string)),
	}
}

func expandDeleteCustomerGatewayInput(d *schema.ResourceData) *computing.DeleteCustomerGatewayInput {
	return &computing.DeleteCustomerGatewayInput{
		CustomerGatewayId: nifcloud.String(d.Id()),
	}
}
