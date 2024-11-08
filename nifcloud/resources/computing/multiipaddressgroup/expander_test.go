package multiipaddressgroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandDescribeMultiIPAddressGroupsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]any{})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeMultiIpAddressGroupsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeMultiIpAddressGroupsInput{
				MultiIpAddressGroupId: []string{"test_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeMultiIPAddressGroupsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandCreateMultiIPAddressGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]any{
		"name":              "test_name",
		"description":       "test_description",
		"availability_zone": "test_availability_zone",
		"ip_address_count":  3,
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateMultiIpAddressGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateMultiIpAddressGroupInput{
				MultiIpAddressGroupName: nifcloud.String("test_name"),
				Description:             nifcloud.String("test_description"),
				Placement: &types.RequestPlacementOfCreateMultiIpAddressGroup{
					AvailabilityZone: nifcloud.String("test_availability_zone"),
				},
				IpAddressCount: nifcloud.Int32(3),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateMultiIPAddressGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyMultiIPAddressGroupAttributeForNameInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]any{
		"name": "test_name",
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyMultiIpAddressGroupAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyMultiIpAddressGroupAttributeInput{
				MultiIpAddressGroupId:   nifcloud.String("test_id"),
				MultiIpAddressGroupName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyMultiIPAddressGroupAttributeForNameInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyMultiIPAddressGroupAttributeForDescriptionInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]any{
		"description": "test_description",
	})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyMultiIpAddressGroupAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyMultiIpAddressGroupAttributeInput{
				MultiIpAddressGroupId: nifcloud.String("test_id"),
				Description:           nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyMultiIPAddressGroupAttributeForDescriptionInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandIncreaseMultiIpAddressCountInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]any{
		"ip_address_count": 3,
	})
	rd.SetId("test_id")
	_ = rd.Set("ip_address_count", 5)

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.IncreaseMultiIpAddressCountInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.IncreaseMultiIpAddressCountInput{
				MultiIpAddressGroupId: nifcloud.String("test_id"),
				IpAddressCount:        nifcloud.Int32(2),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandIncreaseMultiIpAddressCountInput(tt.args, 2)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandReleaseMultiIpAddressesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]any{})
	rd.SetId("test_id")

	ipAddressesToRelease := []string{
		"192.168.0.4",
		"192.168.0.5",
	}

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ReleaseMultiIpAddressesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ReleaseMultiIpAddressesInput{
				MultiIpAddressGroupId: nifcloud.String("test_id"),
				IpAddress: []string{
					"192.168.0.4",
					"192.168.0.5",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandReleaseMultiIpAddressesInput(tt.args, ipAddressesToRelease)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteMultiIPAddressGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]any{})
	rd.SetId("test_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteMultiIpAddressGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteMultiIpAddressGroupInput{
				MultiIpAddressGroupId: nifcloud.String("test_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteMultiIPAddressGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
