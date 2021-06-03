package firewallgroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateFirewallGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"description": "test_description",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.CreateFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.CreateFirewallGroupInput{
				FirewallGroup: &hatoba.CreateFirewallGroupRequestFirewallGroup{
					Name:        nifcloud.String("test_name"),
					Description: nifcloud.String("test_description"),
				},
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

func TestExpandAuthorizeFirewallGroupInput(t *testing.T) {
	rule := map[string]interface{}{
		"protocol":    "TCP",
		"direction":   "IN",
		"from_port":   80,
		"to_port":     80,
		"cidr_ip":     "0.0.0.0/0",
		"description": "test_description",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name": "test_name",
		"rule": []interface{}{rule},
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.AuthorizeFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.AuthorizeFirewallGroupInput{
				FirewallGroupName: nifcloud.String("test_name"),
				Rules: []hatoba.AuthorizeFirewallGroupRequestFirewallRule{
					{
						Protocol:    nifcloud.String("TCP"),
						Direction:   nifcloud.String("IN"),
						FromPort:    nifcloud.Int64(80),
						ToPort:      nifcloud.Int64(80),
						CidrIp:      nifcloud.String("0.0.0.0/0"),
						Description: nifcloud.String("test_description"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandAuthorizeFirewallGroupInput(tt.args, rule)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRevokeFirewallGroupInput(t *testing.T) {
	rule := map[string]interface{}{
		"id": "test_rule_id",
	}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name": "test_name",
		"rule": []interface{}{rule},
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.RevokeFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.RevokeFirewallGroupInput{
				FirewallGroupName: nifcloud.String("test_name"),
				Ids:               nifcloud.String("test_rule_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRevokeFirewallGroupInput(tt.args, rule)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetFirewallGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.GetFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.GetFirewallGroupInput{
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

func TestExpandUpdateFirewallGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_updated_name",
		"description": "test_description",
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.UpdateFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.UpdateFirewallGroupInput{
				FirewallGroupName: nifcloud.String("test_name"),
				FirewallGroup: &hatoba.UpdateFirewallGroupRequestFirewallGroup{
					Name:        nifcloud.String("test_updated_name"),
					Description: nifcloud.String("test_description"),
				},
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

func TestExpandDeleteFirewallGroupInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *hatoba.DeleteFirewallGroupInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &hatoba.DeleteFirewallGroupInput{
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
