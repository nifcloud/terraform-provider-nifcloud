package webproxy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandNiftyCreateWebProxyInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"description":                   "test_description",
		"router_id":                     "test_router_id",
		"router_name":                   "test_router_name",
		"listen_port":                   "test_listen_port",
		"name_server":                   "test_name_server",
		"bypass_interface_network_id":   "test_bypass_interface_network_id",
		"bypass_interface_network_name": "test_bypass_interface_network_name",
		"listen_interface_network_id":   "test_listen_interface_network_id",
		"listen_interface_network_name": "test_listen_interface_network_name",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreateWebProxyInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreateWebProxyInput{
				Description: nifcloud.String("test_description"),
				RouterName:  nifcloud.String("test_router_name"),
				RouterId:    nifcloud.String("test_router_id"),
				ListenPort:  nifcloud.String("test_listen_port"),
				ListenInterface: &computing.RequestListenInterface{
					NetworkId:   nifcloud.String("test_listen_interface_network_id"),
					NetworkName: nifcloud.String("test_listen_interface_network_name"),
				},
				BypassInterface: &computing.RequestBypassInterface{
					NetworkId:   nifcloud.String("test_bypass_interface_network_id"),
					NetworkName: nifcloud.String("test_bypass_interface_network_name"),
				},
				Option: &computing.RequestOption{
					NameServer: nifcloud.String("test_name_server"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreateWebProxyInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDescribeRoutersInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDescribeRoutersInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDescribeRoutersInput{
				RouterId: []string{"test_router_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDescribeRoutersInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDescribeWebProxiesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDescribeWebProxiesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDescribeWebProxiesInput{
				RouterId: []string{"test_router_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDescribeWebProxiesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyWebProxyAttributeInputForDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"description": "test_description",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyWebProxyAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyWebProxyAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyWebProxyAttributeRequestDescription,
				Value:     nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyWebProxyAttributeInputForDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyWebProxyAttributeInputForNameServer(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name_server": "test_name_server",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyWebProxyAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyWebProxyAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyWebProxyAttributeRequestOptionNameServer,
				Value:     nifcloud.String("test_name_server"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyWebProxyAttributeInputForNameServer(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkID(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"listen_interface_network_id": "test_listen_interface_network_id",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyWebProxyAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyWebProxyAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyWebProxyAttributeRequestListenInterfaceNetworkId,
				Value:     nifcloud.String("test_listen_interface_network_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkID(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"listen_interface_network_name": "test_listen_interface_network_name",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyWebProxyAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyWebProxyAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyWebProxyAttributeRequestListenInterfaceNetworkName,
				Value:     nifcloud.String("test_listen_interface_network_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyWebProxyAttributeInputForListenInterfaceNetworkName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bypass_interface_network_name": "test_bypass_interface_network_name",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyWebProxyAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyWebProxyAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyWebProxyAttributeRequestBypassInterfaceNetworkName,
				Value:     nifcloud.String("test_bypass_interface_network_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkID(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bypass_interface_network_id": "test_bypass_interface_network_id",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyWebProxyAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyWebProxyAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyWebProxyAttributeRequestBypassInterfaceNetworkId,
				Value:     nifcloud.String("test_bypass_interface_network_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyWebProxyAttributeInputForBypassInterfaceNetworkID(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyWebProxyAttributeInputForListenPort(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"listen_port": "test_listen_port",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyWebProxyAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyWebProxyAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyWebProxyAttributeRequestListenPort,
				Value:     nifcloud.String("test_listen_port"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyWebProxyAttributeInputForListenPort(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeleteWebProxyInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteWebProxyInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteWebProxyInput{
				RouterId: nifcloud.String("test_router_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteWebProxyInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
