package webproxy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandNiftyCreateWebProxyInput(d *schema.ResourceData) *computing.NiftyCreateWebProxyInput {
	return &computing.NiftyCreateWebProxyInput{
		RouterName:  nifcloud.String(d.Get("router_name").(string)),
		RouterId:    nifcloud.String(d.Get("router_id").(string)),
		Description: nifcloud.String(d.Get("description").(string)),
		ListenInterface: &types.RequestListenInterface{
			NetworkId:   nifcloud.String(d.Get("listen_interface_network_id").(string)),
			NetworkName: nifcloud.String(d.Get("listen_interface_network_name").(string)),
		},
		ListenPort: nifcloud.String(d.Get("listen_port").(string)),
		BypassInterface: &types.RequestBypassInterface{
			NetworkId:   nifcloud.String(d.Get("bypass_interface_network_id").(string)),
			NetworkName: nifcloud.String(d.Get("bypass_interface_network_name").(string)),
		},
		Option: &types.RequestOption{
			NameServer: nifcloud.String(d.Get("name_server").(string)),
		},
	}
}

func expandNiftyDescribeWebProxiesInput(d *schema.ResourceData) *computing.NiftyDescribeWebProxiesInput {
	return &computing.NiftyDescribeWebProxiesInput{
		RouterId: []string{d.Id()},
	}
}

func expandNiftyDescribeRoutersInput(d *schema.ResourceData) *computing.NiftyDescribeRoutersInput {
	return &computing.NiftyDescribeRoutersInput{
		RouterId: []string{d.Id()},
	}
}

func expandNiftyModifyWebProxyAttributeInputForDescription(d *schema.ResourceData) *computing.NiftyModifyWebProxyAttributeInput {
	return &computing.NiftyModifyWebProxyAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: types.AttributeOfNiftyModifyWebProxyAttributeRequestDescription,
		Value:     nifcloud.String(d.Get("description").(string)),
	}
}

func expandNiftyModifyWebProxyAttributeInputForNameServer(d *schema.ResourceData) *computing.NiftyModifyWebProxyAttributeInput {
	return &computing.NiftyModifyWebProxyAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: types.AttributeOfNiftyModifyWebProxyAttributeRequestOptionNameserver,
		Value:     nifcloud.String(d.Get("name_server").(string)),
	}
}

func expandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkID(d *schema.ResourceData) *computing.NiftyModifyWebProxyAttributeInput {
	return &computing.NiftyModifyWebProxyAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: types.AttributeOfNiftyModifyWebProxyAttributeRequestListeninterfaceNetworkid,
		Value:     nifcloud.String(d.Get("listen_interface_network_id").(string)),
	}
}

func expandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkName(d *schema.ResourceData) *computing.NiftyModifyWebProxyAttributeInput {
	return &computing.NiftyModifyWebProxyAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: types.AttributeOfNiftyModifyWebProxyAttributeRequestListeninterfaceNetworkname,
		Value:     nifcloud.String(d.Get("listen_interface_network_name").(string)),
	}
}

func expandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkID(d *schema.ResourceData) *computing.NiftyModifyWebProxyAttributeInput {
	return &computing.NiftyModifyWebProxyAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: types.AttributeOfNiftyModifyWebProxyAttributeRequestBypassinterfaceNetworkid,
		Value:     nifcloud.String(d.Get("bypass_interface_network_id").(string)),
	}
}

func expandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkName(d *schema.ResourceData) *computing.NiftyModifyWebProxyAttributeInput {
	return &computing.NiftyModifyWebProxyAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: types.AttributeOfNiftyModifyWebProxyAttributeRequestBypassinterfaceNetworkname,
		Value:     nifcloud.String(d.Get("bypass_interface_network_name").(string)),
	}
}

func expandNiftyModifyWebProxyAttributeInputForListenPort(d *schema.ResourceData) *computing.NiftyModifyWebProxyAttributeInput {
	return &computing.NiftyModifyWebProxyAttributeInput{
		RouterId:  nifcloud.String(d.Id()),
		Attribute: types.AttributeOfNiftyModifyWebProxyAttributeRequestListenPort,
		Value:     nifcloud.String(d.Get("listen_port").(string)),
	}
}

func expandNiftyDeleteWebProxyInput(d *schema.ResourceData) *computing.NiftyDeleteWebProxyInput {
	return &computing.NiftyDeleteWebProxyInput{
		RouterId: nifcloud.String(d.Id()),
	}
}
