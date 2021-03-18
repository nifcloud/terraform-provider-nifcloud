package dbsecuritygroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateDBSecurityGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"group_name":        "test_group_name",
		"description":       "test_description",
		"availability_zone": "test_zone",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.CreateDBSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.CreateDBSecurityGroupInput{
				NiftyAvailabilityZone:      nifcloud.String("test_zone"),
				DBSecurityGroupName:        nifcloud.String("test_group_name"),
				DBSecurityGroupDescription: nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateDBSecurityGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandAuthorizeDBSecurityGroupIngressInput(t *testing.T) {
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
		want *rdb.AuthorizeDBSecurityGroupIngressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.AuthorizeDBSecurityGroupIngressInput{
				CIDRIP:               nifcloud.String("0.0.0.0/0"),
				EC2SecurityGroupName: nifcloud.String("test_security_group_name"),
				DBSecurityGroupName:  nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandAuthorizeDBSecurityGroupIngressInput(tt.args, rule)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRevokeDBSecurityGroupIngressInput(t *testing.T) {
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
		want *rdb.RevokeDBSecurityGroupIngressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.RevokeDBSecurityGroupIngressInput{
				CIDRIP:               nifcloud.String("0.0.0.0/0"),
				EC2SecurityGroupName: nifcloud.String("test_security_group_name"),
				DBSecurityGroupName:  nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRevokeDBSecurityGroupIngressInput(tt.args, rule)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeDBSecurityGroupsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.DescribeDBSecurityGroupsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.DescribeDBSecurityGroupsInput{
				DBSecurityGroupName: nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeDBSecurityGroupsInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteDBSecurityGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_group_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.DeleteDBSecurityGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.DeleteDBSecurityGroupInput{
				DBSecurityGroupName: nifcloud.String("test_group_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteDBSecurityGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
