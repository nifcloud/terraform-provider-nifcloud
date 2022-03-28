package networkinterface

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateNetworkInterfaceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_address":        "test_ip_address",
		"network_id":        "test_network_id",
		"description":       "test_description",
		"availability_zone": "test_availability_zone",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateNetworkInterfaceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateNetworkInterfaceInput{
				NiftyNetworkId: nifcloud.String("test_network_id"),
				IpAddress:      nifcloud.String("test_ip_address"),
				Description:    nifcloud.String("test_description"),
				Placement: &types.RequestPlacementOfCreateNetworkInterface{
					AvailabilityZone: nifcloud.String("test_availability_zone"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateNetworkInterfaceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeNetworkInterfacesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_network_interface_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeNetworkInterfacesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeNetworkInterfacesInput{
				NetworkInterfaceId: []string{"test_network_interface_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeNetworkInterfacesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyNetworkInterfaceAttributeInputForDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"description": "test_description",
	})
	rd.SetId("test_network_interface_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyNetworkInterfaceAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyNetworkInterfaceAttributeInput{
				NetworkInterfaceId: nifcloud.String("test_network_interface_id"),
				Description:        nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyNetworkInterfaceAttributeInputForDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyNetworkInterfaceAttributeInputForIPAddress(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_address": "test_ip_address",
	})
	rd.SetId("test_network_interface_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyNetworkInterfaceAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyNetworkInterfaceAttributeInput{
				NetworkInterfaceId: nifcloud.String("test_network_interface_id"),
				IpAddress:          nifcloud.String("test_ip_address"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyNetworkInterfaceAttributeInputForIPAddress(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteNetworkInterfaceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_network_interface_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteNetworkInterfaceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteNetworkInterfaceInput{
				NetworkInterfaceId: nifcloud.String("test_network_interface_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteNetworkInterfaceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
