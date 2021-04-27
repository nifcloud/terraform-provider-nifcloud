package vpngateway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateVpnGatewayInput_netid_and_netname(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_name",
		"type":              "test_type",
		"availability_zone": "test_availability_zone",
		"accounting_type":   "test_accounting_type",
		"description":       "test_description",
		"network_id":        "test_network_id",
		"network_name":      "test_network_name",
		"ip_address":        "test_ip_address",
		"security_group":    "test_security_group",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateVpnGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateVpnGatewayInput{
				AccountingType:             computing.AccountingTypeOfCreateVpnGatewayRequest("test_accounting_type"),
				NiftyVpnGatewayDescription: nifcloud.String("test_description"),
				NiftyVpnGatewayName:        nifcloud.String("test_name"),
				NiftyVpnGatewayType:        computing.NiftyVpnGatewayTypeOfCreateVpnGatewayRequest("test_type"),
				Placement: &computing.RequestPlacementOfCreateVpnGateway{
					AvailabilityZone: nifcloud.String("test_availability_zone"),
				},
				NiftyNetwork: &computing.RequestNiftyNetwork{
					NetworkId:   nifcloud.String("test_network_id"),
					NetworkName: nil,
					IpAddress:   nifcloud.String("test_ip_address"),
				},
				SecurityGroup: []string{"test_security_group"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateVpnGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandCreateVpnGatewayInput_netid_only(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_name",
		"type":              "test_type",
		"availability_zone": "test_availability_zone",
		"accounting_type":   "test_accounting_type",
		"description":       "test_description",
		"network_id":        "test_network_id",
		"ip_address":        "test_ip_address",
		"security_group":    "test_security_group",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateVpnGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateVpnGatewayInput{
				AccountingType:             computing.AccountingTypeOfCreateVpnGatewayRequest("test_accounting_type"),
				NiftyVpnGatewayDescription: nifcloud.String("test_description"),
				NiftyVpnGatewayName:        nifcloud.String("test_name"),
				NiftyVpnGatewayType:        computing.NiftyVpnGatewayTypeOfCreateVpnGatewayRequest("test_type"),
				Placement: &computing.RequestPlacementOfCreateVpnGateway{
					AvailabilityZone: nifcloud.String("test_availability_zone"),
				},
				NiftyNetwork: &computing.RequestNiftyNetwork{
					NetworkId:   nifcloud.String("test_network_id"),
					NetworkName: nil,
					IpAddress:   nifcloud.String("test_ip_address"),
				},
				SecurityGroup: []string{"test_security_group"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateVpnGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandCreateVpnGatewayInput_netname_only(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_name",
		"type":              "test_type",
		"availability_zone": "test_availability_zone",
		"accounting_type":   "test_accounting_type",
		"description":       "test_description",
		"network_name":      "test_network_name",
		"ip_address":        "test_ip_address",
		"security_group":    "test_security_group",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateVpnGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateVpnGatewayInput{
				AccountingType:             computing.AccountingTypeOfCreateVpnGatewayRequest("test_accounting_type"),
				NiftyVpnGatewayDescription: nifcloud.String("test_description"),
				NiftyVpnGatewayName:        nifcloud.String("test_name"),
				NiftyVpnGatewayType:        computing.NiftyVpnGatewayTypeOfCreateVpnGatewayRequest("test_type"),
				Placement: &computing.RequestPlacementOfCreateVpnGateway{
					AvailabilityZone: nifcloud.String("test_availability_zone"),
				},
				NiftyNetwork: &computing.RequestNiftyNetwork{
					NetworkId:   nil,
					NetworkName: nifcloud.String("test_network_name"),
					IpAddress:   nifcloud.String("test_ip_address"),
				},
				SecurityGroup: []string{"test_security_group"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateVpnGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestExpandDescribeVpnGatewaysInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeVpnGatewaysInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeVpnGatewaysInput{
				VpnGatewayId: []string{"test_vpngateway_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeVpnGatewaysInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteVpnGatewaysInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteVpnGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteVpnGatewayInput{
				VpnGatewayId: nifcloud.String("test_vpngateway_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteVpnGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyVpnGatewayAttributeInputForAccountingType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type": "test_accounting_type",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyVpnGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyVpnGatewayAttributeInput{
				VpnGatewayId: nifcloud.String("test_vpngateway_id"),
				Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayAccountingType,
				Value:        nifcloud.String("test_accounting_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyVpnGatewayAttributeInputForAccountingType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyVpnGatewayAttributeInputForVpnGatewayDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"description": "test_description",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyVpnGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyVpnGatewayAttributeInput{
				VpnGatewayId: nifcloud.String("test_vpngateway_id"),
				Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayDescription,
				Value:        nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyVpnGatewayAttributeInputForVpnGatewayName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name": "test_name",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyVpnGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyVpnGatewayAttributeInput{
				VpnGatewayId: nifcloud.String("test_vpngateway_id"),
				Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayName,
				Value:        nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyVpnGatewayAttributeInputForVpnGatewayType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"type": "test_type",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyVpnGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyVpnGatewayAttributeInput{
				VpnGatewayId: nifcloud.String("test_vpngateway_id"),
				Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestNiftyVpnGatewayType,
				Value:        nifcloud.String("test_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyVpnGatewayAttributeInputForVpnGatewayType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyUpdateVpnGatewayNetworkInterfacesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_address": "test_ip_address",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyUpdateVpnGatewayNetworkInterfacesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyUpdateVpnGatewayNetworkInterfacesInput{
				VpnGatewayId: nifcloud.String("test_vpngateway_id"),
				NetworkInterface: &computing.RequestNetworkInterfaceOfNiftyUpdateVpnGatewayNetworkInterfaces{
					IpAddress: nifcloud.String("test_ip_address"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyUpdateVpnGatewayNetworkInterfacesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyVpnGatewayAttributeInputForSecurityGroup(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"security_group": "test_security_group",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyVpnGatewayAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyVpnGatewayAttributeInput{
				VpnGatewayId: nifcloud.String("test_vpngateway_id"),
				Attribute:    computing.AttributeOfNiftyModifyVpnGatewayAttributeRequestGroupId,
				Value:        nifcloud.String("test_security_group"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyVpnGatewayAttributeInputForSecurityGroup(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyAssociateRouteTableWithVpnGatewayInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route_table_id": "test_route_table_id",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyAssociateRouteTableWithVpnGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyAssociateRouteTableWithVpnGatewayInput{
				VpnGatewayId: nifcloud.String("test_vpngateway_id"),
				RouteTableId: nifcloud.String("test_route_table_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyAssociateRouteTableWithVpnGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDisassociateRouteTableFromVpnGatewayInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route_table_association_id": "test_route_table_association_id",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDisassociateRouteTableFromVpnGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDisassociateRouteTableFromVpnGatewayInput{
				AssociationId: nifcloud.String("test_route_table_association_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDisassociateRouteTableFromVpnGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyReplaceRouteTableAssociationWithVpnGatewayInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route_table_association_id": "test_route_table_association_id",
		"route_table_id":             "test_route_table_id",
	})
	rd.SetId("test_vpngateway_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyReplaceRouteTableAssociationWithVpnGatewayInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyReplaceRouteTableAssociationWithVpnGatewayInput{
				RouteTableId:  nifcloud.String("test_route_table_id"),
				AssociationId: nifcloud.String("test_route_table_association_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyReplaceRouteTableAssociationWithVpnGatewayInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
