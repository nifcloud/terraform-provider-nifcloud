package elasticip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandAllocateAddressInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_type":           false,
		"availability_zone": "east-21",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.AllocateAddressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.AllocateAddressInput{
				NiftyPrivateIp: nifcloud.Bool(false),
				Placement: &computing.RequestPlacementOfAllocateAddress{
					AvailabilityZone: nifcloud.String("east-21"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandAllocateAddressInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPrivateNiftyModifyAddressAttributeInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_type":     true,
		"private_ip":  "192.168.0.1",
		"description": "test_description",
	})
	rd.SetId("192.168.0.1")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyAddressAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyAddressAttributeInput{
				PrivateIpAddress: nifcloud.String("192.168.0.1"),
				Attribute:        "description",
				Value:            nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyAddressAttributeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPublicNiftyModifyAddressAttributeInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_type":     false,
		"public_ip":   "192.0.2.1",
		"description": "test_description",
	})
	rd.SetId("192.0.2.1")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyAddressAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyAddressAttributeInput{
				PublicIp:  nifcloud.String("192.0.2.1"),
				Attribute: "description",
				Value:     nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyAddressAttributeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPrivateDescribeAddressesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_type":    true,
		"private_ip": "192.168.0.1",
	})
	rd.SetId("192.168.0.1")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeAddressesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeAddressesInput{
				PrivateIpAddress: []string{"192.168.0.1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeAddressesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPublicDescribeAddressesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_type":   false,
		"public_ip": "192.0.2.1",
	})
	rd.SetId("192.0.2.1")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeAddressesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeAddressesInput{
				PublicIp: []string{"192.0.2.1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeAddressesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPrivateReleaseAddressInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_type":    true,
		"private_ip": "192.168.0.1",
	})
	rd.SetId("192.168.0.1")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ReleaseAddressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ReleaseAddressInput{
				PrivateIpAddress: nifcloud.String("192.168.0.1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandReleaseAddressInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPublicReleaseAddressInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"ip_type":   false,
		"public_ip": "192.0.2.1",
	})
	rd.SetId("192.0.2.1")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ReleaseAddressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ReleaseAddressInput{
				PublicIp: nifcloud.String("192.0.2.1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandReleaseAddressInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
