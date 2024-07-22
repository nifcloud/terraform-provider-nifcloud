package devopsfirewallgroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateFirewallGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_name",
		"availability_zone": "east-11",
		"description":       "test_description",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.CreateFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.CreateFirewallGroupInput{
				FirewallGroupName: nifcloud.String("test_name"),
				AvailabilityZone:  types.AvailabilityZoneOfCreateFirewallGroupRequest("east-11"),
				Description:       nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateFirewallGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateFirewallGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name_changed",
		"description": "test_description",
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.UpdateFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.UpdateFirewallGroupInput{
				FirewallGroupName:        nifcloud.String("test_name"),
				ChangedFirewallGroupName: nifcloud.String("test_name_changed"),
				Description:              nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateFirewallGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetFirewallGroupsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.GetFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.GetFirewallGroupInput{
				FirewallGroupName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetFirewallGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteFirewallGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devops.DeleteFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devops.DeleteFirewallGroupInput{
				FirewallGroupName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteFirewallGroupInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandAuthorizeFirewallRules(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"description": "test_description",
	})
	rd.SetId("test_name")

	type args struct {
		d     *schema.ResourceData
		rules []types.RequestRules
	}
	tests := []struct {
		name string
		args args
		want *devops.AuthorizeFirewallRulesInput
	}{
		{
			name: "expands the resource data",
			args: args{
				d: rd,
				rules: []types.RequestRules{
					{
						Protocol:    types.ProtocolOfrulesForAuthorizeFirewallRules("TCP"),
						Port:        nifcloud.Int32(443),
						CidrIp:      types.CidrIpOfrulesForAuthorizeFirewallRules("172.16.0.0/24"),
						Description: nifcloud.String("test_description"),
					},
				},
			},
			want: &devops.AuthorizeFirewallRulesInput{
				FirewallGroupName: nifcloud.String("test_name"),
				Rules: []types.RequestRules{
					{
						Protocol:    types.ProtocolOfrulesForAuthorizeFirewallRules("TCP"),
						Port:        nifcloud.Int32(443),
						CidrIp:      types.CidrIpOfrulesForAuthorizeFirewallRules("172.16.0.0/24"),
						Description: nifcloud.String("test_description"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandAuthorizeFirewallRulesInput(tt.args.d, tt.args.rules)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRevokeFirewallRules(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"description": "test_description",
	})
	rd.SetId("test_name")

	type args struct {
		d       *schema.ResourceData
		ruleIds []string
	}
	tests := []struct {
		name string
		args args
		want *devops.RevokeFirewallRulesInput
	}{
		{
			name: "expands the resource data",
			args: args{
				d:       rd,
				ruleIds: []string{"test_id_01", "test_id_02"},
			},
			want: &devops.RevokeFirewallRulesInput{
				FirewallGroupName: nifcloud.String("test_name"),
				Ids:               nifcloud.String("test_id_01,test_id_02"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRevokeFirewallRulesInput(tt.args.d, tt.args.ruleIds)
			assert.Equal(t, tt.want, got)
		})
	}
}
