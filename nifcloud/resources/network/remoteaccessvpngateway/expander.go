package remoteaccessvpngateway

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandCreateRemoteAccessVpnGatewayInput(d *schema.ResourceData) *computing.CreateRemoteAccessVpnGatewayInput {
	networkInterface := make([]types.RequestNetworkInterfaceOfCreateRemoteAccessVpnGateway, 1)

	niMap := d.Get("network_interface").([]interface{})[0].(map[string]interface{})
	if v, ok := niMap["network_id"].(string); ok && v != "" {
		networkInterface[0].NetworkId = nifcloud.String(v)
	}

	if v, ok := niMap["ip_address"].(string); ok && v != "" {
		networkInterface[0].IpAddress = nifcloud.String(v)
	}

	var cipherSuite []string
	for _, c := range d.Get("cipher_suite").([]interface{}) {
		cipherSuite = append(cipherSuite, c.(string))
	}

	accountingType, _ := strconv.ParseInt(d.Get("accounting_type").(string), 10, 32)

	input := &computing.CreateRemoteAccessVpnGatewayInput{
		RemoteAccessVpnGatewayName: nifcloud.String(d.Get("name").(string)),
		RemoteAccessVpnGatewayType: types.RemoteAccessVpnGatewayTypeOfCreateRemoteAccessVpnGatewayRequest(d.Get("type").(string)),
		Placement: &types.RequestPlacementOfCreateRemoteAccessVpnGateway{
			AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		},
		AccountingType:   nifcloud.Int32(int32(accountingType)),
		Description:      nifcloud.String(d.Get("description").(string)),
		PoolNetworkCidr:  nifcloud.String(d.Get("pool_network_cidr").(string)),
		NetworkInterface: networkInterface,
		SSLCertificateId: nifcloud.String(d.Get("ssl_certificate_id").(string)),
		CipherSuite:      cipherSuite,
	}

	if v, ok := d.GetOk("ca_certificate_id"); ok {
		input.CACertificateId = nifcloud.String(v.(string))
	}

	return input
}

func expandCreateRemoteAccessVpnGatewayUsersInput(d *schema.ResourceData, user map[string]interface{}) *computing.CreateRemoteAccessVpnGatewayUsersInput {
	return &computing.CreateRemoteAccessVpnGatewayUsersInput{
		RemoteAccessVpnGatewayId: nifcloud.String(d.Id()),
		RemoteUser: []types.RequestRemoteUser{{
			UserName:    nifcloud.String(user["name"].(string)),
			Password:    nifcloud.String(user["password"].(string)),
			Description: nifcloud.String(user["description"].(string)),
		}},
	}
}

func expandDeleteRemoteAccessVpnGatewayUsersInput(d *schema.ResourceData, user map[string]interface{}) *computing.DeleteRemoteAccessVpnGatewayUsersInput {
	return &computing.DeleteRemoteAccessVpnGatewayUsersInput{
		RemoteAccessVpnGatewayId: nifcloud.String(d.Id()),
		RemoteUser: []types.RequestRemoteUserOfDeleteRemoteAccessVpnGatewayUsers{{
			UserName: nifcloud.String(user["name"].(string)),
		}},
	}
}

func expandDescribeRemoteAccessVpnGatewaysInput(d *schema.ResourceData) *computing.DescribeRemoteAccessVpnGatewaysInput {
	return &computing.DescribeRemoteAccessVpnGatewaysInput{
		RemoteAccessVpnGatewayId: []string{d.Id()},
	}
}

func expandDescribeRemoteAccessVpnGatewayClientConfigInput(d *schema.ResourceData) *computing.DescribeRemoteAccessVpnGatewayClientConfigInput {
	return &computing.DescribeRemoteAccessVpnGatewayClientConfigInput{
		RemoteAccessVpnGatewayId: nifcloud.String(d.Id()),
	}
}

func expandModifyRemoteAccessVpnGatewayAttributeInputForRemoteAccessVpnGatewayName(d *schema.ResourceData) *computing.ModifyRemoteAccessVpnGatewayAttributeInput {
	return &computing.ModifyRemoteAccessVpnGatewayAttributeInput{
		RemoteAccessVpnGatewayId:   nifcloud.String(d.Id()),
		RemoteAccessVpnGatewayName: nifcloud.String(d.Get("name").(string)),
	}
}

func expandModifyRemoteAccessVpnGatewayAttributeInputForAccountingType(d *schema.ResourceData) *computing.ModifyRemoteAccessVpnGatewayAttributeInput {
	return &computing.ModifyRemoteAccessVpnGatewayAttributeInput{
		RemoteAccessVpnGatewayId: nifcloud.String(d.Id()),
		AccountingType:           types.AccountingTypeOfModifyRemoteAccessVpnGatewayAttributeRequest(d.Get("accounting_type").(string)),
	}
}

func expandModifyRemoteAccessVpnGatewayAttributeInputForDescription(d *schema.ResourceData) *computing.ModifyRemoteAccessVpnGatewayAttributeInput {
	return &computing.ModifyRemoteAccessVpnGatewayAttributeInput{
		RemoteAccessVpnGatewayId: nifcloud.String(d.Id()),
		Description:              nifcloud.String(d.Get("description").(string)),
	}
}

func expandModifyRemoteAccessVpnGatewayAttributeInputForType(d *schema.ResourceData) *computing.ModifyRemoteAccessVpnGatewayAttributeInput {
	return &computing.ModifyRemoteAccessVpnGatewayAttributeInput{
		RemoteAccessVpnGatewayId:   nifcloud.String(d.Id()),
		RemoteAccessVpnGatewayType: types.RemoteAccessVpnGatewayTypeOfModifyRemoteAccessVpnGatewayAttributeRequest(d.Get("type").(string)),
	}
}

func expandSetRemoteAccessVpnGatewayCACertificateInput(d *schema.ResourceData) *computing.SetRemoteAccessVpnGatewayCACertificateInput {
	return &computing.SetRemoteAccessVpnGatewayCACertificateInput{
		CACertificateId:          nifcloud.String(d.Get("ca_certificate_id").(string)),
		RemoteAccessVpnGatewayId: nifcloud.String(d.Id()),
	}
}

func expandUnsetRemoteAccessVpnGatewayCACertificateInput(d *schema.ResourceData) *computing.UnsetRemoteAccessVpnGatewayCACertificateInput {
	return &computing.UnsetRemoteAccessVpnGatewayCACertificateInput{
		RemoteAccessVpnGatewayId: nifcloud.String(d.Id()),
	}
}

func expandSetRemoteAccessVpnGatewaySSLCertificateInput(d *schema.ResourceData) *computing.SetRemoteAccessVpnGatewaySSLCertificateInput {
	return &computing.SetRemoteAccessVpnGatewaySSLCertificateInput{
		SSLCertificateId:         nifcloud.String(d.Get("ssl_certificate_id").(string)),
		RemoteAccessVpnGatewayId: nifcloud.String(d.Id()),
	}
}

func expandDeleteRemoteAccessVpnGatewayInput(d *schema.ResourceData) *computing.DeleteRemoteAccessVpnGatewayInput {
	return &computing.DeleteRemoteAccessVpnGatewayInput{
		RemoteAccessVpnGatewayId: nifcloud.String(d.Id()),
	}
}
