package nassecuritygroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateNASSecurityGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name":        "test_group_name",
		"description":       "test_description",
		"availability_zone": "test_zone",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.CreateNASSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &nas.CreateNASSecurityGroupInput{
				AvailabilityZone:            nifcloud.String("test_zone"),
				NASSecurityGroupName:        nifcloud.String("test_group_name"),
				NASSecurityGroupDescription: nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateNASSecurityGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandAuthorizeNASSecurityGroupIngressInput(t *testing.T) {
	rule := map[string]interface{}{
		"cidr_ip":             "0.0.0.0/0",
		"security_group_name": "test_security_group_name",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"rule":       []interface{}{rule},
		"group_name": "test_group_name",
	})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.AuthorizeNASSecurityGroupIngressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &nas.AuthorizeNASSecurityGroupIngressInput{
				CIDRIP:               nifcloud.String("0.0.0.0/0"),
				SecurityGroupName:    nifcloud.String("test_security_group_name"),
				NASSecurityGroupName: nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandAuthorizeNASSecurityGroupIngressInput(tt.args, rule)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRevokeNASSecurityGroupIngressInput(t *testing.T) {
	rule := map[string]interface{}{
		"cidr_ip":             "0.0.0.0/0",
		"security_group_name": "test_security_group_name",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"rule":       []interface{}{rule},
		"group_name": "test_group_name",
	})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.RevokeNASSecurityGroupIngressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &nas.RevokeNASSecurityGroupIngressInput{
				CIDRIP:               nifcloud.String("0.0.0.0/0"),
				SecurityGroupName:    nifcloud.String("test_security_group_name"),
				NASSecurityGroupName: nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRevokeNASSecurityGroupIngressInput(tt.args, rule)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeNASSecurityGroupsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.DescribeNASSecurityGroupsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &nas.DescribeNASSecurityGroupsInput{
				NASSecurityGroupName: nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeNASSecurityGroupsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyNASSecurityGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name":  "test_new_group_name",
		"description": "test_description",
	})
	rd.SetId("test_old_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.ModifyNASSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &nas.ModifyNASSecurityGroupInput{
				NewNASSecurityGroupName:     nifcloud.String("test_new_group_name"),
				NASSecurityGroupName:        nifcloud.String("test_old_group_name"),
				NASSecurityGroupDescription: nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyNASSecurityGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteNASSecurityGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *nas.DeleteNASSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &nas.DeleteNASSecurityGroupInput{
				NASSecurityGroupName: nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteNASSecurityGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
