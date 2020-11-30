package privatelan

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandNiftyCreatePrivateLanInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":   "1",
		"availability_zone": "test_availability_zone",
		"cidr_block":        "test_cidr_block",
		"description":       "test_description",
		"private_lan_name":  "test_privateLan_name",
	})
	rd.SetId("test_network_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreatePrivateLanInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreatePrivateLanInput{
				AccountingType:   "1",
				AvailabilityZone: nifcloud.String("test_availability_zone"),
				CidrBlock:        nifcloud.String("test_cidr_block"),
				Description:      nifcloud.String("test_description"),
				PrivateLanName:   nifcloud.String("test_privateLan_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreatePrivateLanInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDescribePrivateLansInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"netowork_id": "test_network_id",
	})
	rd.SetId("test_network_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDescribePrivateLansInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDescribePrivateLansInput{
				NetworkId: []string{"test_network_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDescribePrivateLansInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestExpandNiftyModifyPrivateLanAttributeInputForAccountingType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"network_id":      "test_network_id",
		"accounting_type": "test_accounting_type",
	})
	rd.SetId("test_network_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyPrivateLanAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyPrivateLanAttributeInput{
				NetworkId: nifcloud.String("test_network_id"),
				Attribute: "accountingType",
				Value:     nifcloud.String("test_accounting_type"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyPrivateLanAttributeInputForAccountingType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyPrivateLanAttributeInputForDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"network_id":  "test_network_id",
		"description": "test_description",
	})
	rd.SetId("test_network_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyPrivateLanAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyPrivateLanAttributeInput{
				NetworkId: nifcloud.String("test_network_id"),
				Attribute: "description",
				Value:     nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyPrivateLanAttributeInputForDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyPrivateLanAttributeInputForPrivateLanName(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"network_id":       "test_network_id",
		"private_lan_name": "test_private_lan_name",
	})
	rd.SetId("test_network_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyPrivateLanAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyPrivateLanAttributeInput{
				NetworkId: nifcloud.String("test_network_id"),
				Attribute: "privateLanName",
				Value:     nifcloud.String("test_private_lan_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyPrivateLanAttributeInputForPrivateLanName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyPrivateLanAttributeInputForCidrBlock(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"network_id": "test_network_id",
		"cidr_block": "test_cidr_block",
	})
	rd.SetId("test_network_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyPrivateLanAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyPrivateLanAttributeInput{
				NetworkId: nifcloud.String("test_network_id"),
				Attribute: "cidrBlock",
				Value:     nifcloud.String("test_cidr_block"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyPrivateLanAttributeInputForCidrBlock(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
