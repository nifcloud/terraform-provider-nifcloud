package securitygroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateSecurityGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name":        "test_group_name",
		"description":       "test_description",
		"availability_zone": "test_availability_zone",
	})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateSecurityGroupInput{
				GroupName:        nifcloud.String("test_group_name"),
				GroupDescription: nifcloud.String("test_description"),
				Placement: &computing.RequestPlacementOfCreateSecurityGroup{
					AvailabilityZone: nifcloud.String("test_availability_zone"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateSecurityGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateSecurityGroupInputForName(t *testing.T) {
	r := New()
	rd := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"group_name": "test_group_name",
	})
	rd.SetId("test_group_name")
	dn := r.Data(rd.State())

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.UpdateSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: dn,
			want: &computing.UpdateSecurityGroupInput{
				GroupName:       nifcloud.String("test_group_name"),
				GroupNameUpdate: nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateSecurityGroupInputForName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateSecurityGroupInputForDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name":  "test_group_name",
		"description": "test_description",
	})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.UpdateSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.UpdateSecurityGroupInput{
				GroupName:              nifcloud.String("test_group_name"),
				GroupDescriptionUpdate: nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateSecurityGroupInputForDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateSecurityGroupInputForLogLimit(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name": "test_group_name",
		"log_limit":  1000,
	})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.UpdateSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.UpdateSecurityGroupInput{
				GroupName:           nifcloud.String("test_group_name"),
				GroupLogLimitUpdate: nifcloud.Int64(1000),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateSecurityGroupInputForLogLimit(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeSecurityGroupsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name": "test_group_name",
	})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeSecurityGroupsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeSecurityGroupsInput{
				GroupName: []string{"test_group_name"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeSecurityGroupsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteSecurityGroupsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name": "test_group_name",
	})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteSecurityGroupInput{
				GroupName: nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteSecurityGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRevokeSecurityGroupIngressInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name": "test_group_name",
	})
	rd.SetId("test_group_name")

	ipPermissions := []computing.RequestIpPermissionsOfRevokeSecurityGroupIngress{}

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.RevokeSecurityGroupIngressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.RevokeSecurityGroupIngressInput{
				GroupName:     nifcloud.String("test_group_name"),
				IpPermissions: ipPermissions,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRevokeSecurityGroupIngressInput(tt.args, ipPermissions)
			assert.Equal(t, tt.want, got)
		})
	}
}
