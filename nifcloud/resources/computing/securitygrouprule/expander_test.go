package securitygrouprule

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandAuthorizeSecurityGroupIngressInputList(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"type":                       "IN",
		"cidr_ip":                    "0.0.0.0/0",
		"from_port":                  1,
		"security_group_names":       []interface{}{"test_security_group_name"},
		"protocol":                   "TCP",
		"source_security_group_name": "test_source_security_group_name",
		"to_port":                    65535,
		"description":                "test_description",
	})
	rd.SetId("some")

	tests := []struct {
		name string
		args *schema.ResourceData
		want []*computing.AuthorizeSecurityGroupIngressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: []*computing.AuthorizeSecurityGroupIngressInput{
				{
					GroupName: nifcloud.String("test_security_group_name"),
					IpPermissions: []computing.RequestIpPermissions{
						{
							IpProtocol:  computing.IpProtocolOfIpPermissionsForAuthorizeSecurityGroupIngressTcp,
							InOut:       computing.InOutOfIpPermissionsForAuthorizeSecurityGroupIngressIn,
							Description: nifcloud.String("test_description"),
							FromPort:    nifcloud.Int64(1),
							ToPort:      nifcloud.Int64(65535),
							ListOfRequestIpRanges: []computing.RequestIpRanges{
								{CidrIp: nifcloud.String("0.0.0.0/0")},
							},
							ListOfRequestGroups: []computing.RequestGroups{
								{GroupName: nifcloud.String("test_source_security_group_name")},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandAuthorizeSecurityGroupIngressInputList(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRevokeSecurityGroupIngressInputList(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"type":                       "IN",
		"cidr_ip":                    "0.0.0.0/0",
		"from_port":                  1,
		"security_group_names":       []interface{}{"test_security_group_name"},
		"protocol":                   "TCP",
		"source_security_group_name": "test_source_security_group_name",
		"to_port":                    65535,
	})
	rd.SetId("some")

	tests := []struct {
		name string
		args *schema.ResourceData
		want []*computing.RevokeSecurityGroupIngressInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: []*computing.RevokeSecurityGroupIngressInput{
				{
					GroupName: nifcloud.String("test_security_group_name"),
					IpPermissions: []computing.RequestIpPermissionsOfRevokeSecurityGroupIngress{
						{
							IpProtocol: computing.IpProtocolOfIpPermissionsForRevokeSecurityGroupIngressTcp,
							InOut:      computing.InOutOfIpPermissionsForRevokeSecurityGroupIngressIn,
							FromPort:   nifcloud.Int64(1),
							ToPort:     nifcloud.Int64(65535),
							ListOfRequestIpRanges: []computing.RequestIpRanges{
								{CidrIp: nifcloud.String("0.0.0.0/0")},
							},
							ListOfRequestGroups: []computing.RequestGroups{
								{GroupName: nifcloud.String("test_source_security_group_name")},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRevokeSecurityGroupIngressInputList(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeSecurityGroupsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"security_group_names": []interface{}{"test_security_group_name"},
	})
	rd.SetId("some")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeSecurityGroupsInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeSecurityGroupsInput{
				GroupName: []string{"test_security_group_name"},
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
