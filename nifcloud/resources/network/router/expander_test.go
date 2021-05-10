package router

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandNiftyCreateRouterInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":   "test_accounting_type",
		"availability_zone": "test_availability_zone",
		"description":       "test_description",
		"nat_table_id":      "test_nat_table_id",
		"network_interface": []interface{}{map[string]interface{}{
			"dhcp":            true,
			"dhcp_config_id":  "test_dhcp_config_id",
			"dhcp_options_id": "test_dhcp_options_id",
			"ip_address":      "test_ip_address",
			"network_id":      "test_network_id",
			"network_name":    "test_network_name",
		}},
		"name":           "test_router_name",
		"route_table_id": "test_route_table_id",
		"security_group": "test_security_group",
		"type":           "test_type",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreateRouterInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreateRouterInput{
				AccountingType:   computing.AccountingTypeOfNiftyCreateRouterRequest("test_accounting_type"),
				AvailabilityZone: nifcloud.String("test_availability_zone"),
				Description:      nifcloud.String("test_description"),
				NetworkInterface: []computing.RequestNetworkInterfaceOfNiftyCreateRouter{
					{
						Dhcp:          nifcloud.Bool(true),
						DhcpConfigId:  nifcloud.String("test_dhcp_config_id"),
						DhcpOptionsId: nifcloud.String("test_dhcp_options_id"),
						IpAddress:     nifcloud.String("test_ip_address"),
						NetworkId:     nifcloud.String("test_network_id"),
						NetworkName:   nifcloud.String("test_network_name"),
					},
				},
				RouterName:    nifcloud.String("test_router_name"),
				SecurityGroup: []string{"test_security_group"},
				Type:          computing.TypeOfNiftyCreateRouterRequest("test_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreateRouterInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandAssociateRouteTableInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route_table_id": "test_route_table_id",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.AssociateRouteTableInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.AssociateRouteTableInput{
				RouteTableId: nifcloud.String("test_route_table_id"),
				RouterId:     nifcloud.String("test_router_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandAssociateRouteTableInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyAssociateNatTableInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"nat_table_id": "test_nat_table_id",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyAssociateNatTableInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyAssociateNatTableInput{
				NatTableId: nifcloud.String("test_nat_table_id"),
				RouterId:   nifcloud.String("test_router_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyAssociateNatTableInput(tt.args)
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

func TestExpandNiftyModifyRouterAttributeInputForRouterName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name": "test_router_name",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyRouterAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyRouterAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestRouterName,
				Value:     nifcloud.String("test_router_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyRouterAttributeInputForRouterName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyRouterAttributeInputForAccountingType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type": "test_accounting_type",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyRouterAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyRouterAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestAccountingType,
				Value:     nifcloud.String("test_accounting_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyRouterAttributeInputForAccountingType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyRouterAttributeInputForDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"description": "test_description",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyRouterAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyRouterAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestDescription,
				Value:     nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyRouterAttributeInputForDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyRouterAttributeInputForType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"type": "test_type",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyRouterAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyRouterAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestType,
				Value:     nifcloud.String("test_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyRouterAttributeInputForType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyRouterAttributeInputForSecurityGroup(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"security_group": "test_security_group",
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyRouterAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyRouterAttributeInput{
				RouterId:  nifcloud.String("test_router_id"),
				Attribute: computing.AttributeOfNiftyModifyRouterAttributeRequestGroupId,
				Value:     nifcloud.String("test_security_group"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyRouterAttributeInputForSecurityGroup(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeregisterRoutersFromSecurityGroupInput(t *testing.T) {
	state := &terraform.InstanceState{
		Attributes: map[string]string{
			"security_group": "test_security_group",
		},
	}
	diff := &terraform.InstanceDiff{
		Attributes: map[string]*terraform.ResourceAttrDiff{
			"security_group": {
				Old:         "test_security_group",
				New:         "",
				RequiresNew: true,
			},
		},
	}
	rd, err := schema.InternalMap(newSchema()).Data(state, diff)
	if err != nil {
		t.Fatal(err)
	}
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeregisterRoutersFromSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeregisterRoutersFromSecurityGroupInput{
				GroupName: nifcloud.String("test_security_group"),
				RouterSet: []computing.RequestRouterSetOfNiftyDeregisterRoutersFromSecurityGroup{
					{
						RouterId: nifcloud.String("test_router_id"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeregisterRoutersFromSecurityGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyUpdateRouterNetworkInterfacesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"network_interface": []interface{}{map[string]interface{}{
			"dhcp":            true,
			"dhcp_config_id":  "test_dhcp_config_id",
			"dhcp_options_id": "test_dhcp_options_id",
			"ip_address":      "test_ip_address",
			"network_id":      "test_network_id",
			"network_name":    "test_network_name",
		}},
	})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyUpdateRouterNetworkInterfacesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyUpdateRouterNetworkInterfacesInput{
				RouterId: nifcloud.String("test_router_id"),
				NetworkInterface: []computing.RequestNetworkInterfaceOfNiftyUpdateRouterNetworkInterfaces{
					{
						Dhcp:          nifcloud.Bool(true),
						DhcpConfigId:  nifcloud.String("test_dhcp_config_id"),
						DhcpOptionsId: nifcloud.String("test_dhcp_options_id"),
						IpAddress:     nifcloud.String("test_ip_address"),
						NetworkId:     nifcloud.String("test_network_id"),
						NetworkName:   nifcloud.String("test_network_name"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyUpdateRouterNetworkInterfacesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDisassociateNatTableInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"nat_table_association_id": "test_nat_table_association_id",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDisassociateNatTableInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDisassociateNatTableInput{
				AssociationId: nifcloud.String("test_nat_table_association_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDisassociateNatTableInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_expandNiftyReplaceNatTableAssociationInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"nat_table_id":             "test_nat_table_id",
		"nat_table_association_id": "test_nat_table_association_id",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyReplaceNatTableAssociationInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyReplaceNatTableAssociationInput{
				NatTableId:    nifcloud.String("test_nat_table_id"),
				AssociationId: nifcloud.String("test_nat_table_association_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyReplaceNatTableAssociationInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDisassociateRouteTableInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route_table_association_id": "test_route_table_association_id",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DisassociateRouteTableInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DisassociateRouteTableInput{
				AssociationId: nifcloud.String("test_route_table_association_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDisassociateRouteTableInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandReplaceRouteTableAssociation(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"route_table_id":             "test_route_table_id",
		"route_table_association_id": "test_route_table_association_id",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ReplaceRouteTableAssociationInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ReplaceRouteTableAssociationInput{
				RouteTableId:  nifcloud.String("test_route_table_id"),
				AssociationId: nifcloud.String("test_route_table_association_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandReplaceRouteTableAssociation(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeleteRouterInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_router_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteRouterInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteRouterInput{
				RouterId: nifcloud.String("test_router_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteRouterInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
