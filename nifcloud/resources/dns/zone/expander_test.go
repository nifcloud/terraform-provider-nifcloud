package zone

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateHostedZoneInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":    "test_name",
		"comment": "test_comment",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *dns.CreateHostedZoneInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &dns.CreateHostedZoneInput{
				Name: nifcloud.String("test_name"),
				RequestHostedZoneConfig: &dns.RequestHostedZoneConfig{
					Comment: nifcloud.String("test_comment"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateHostedZoneInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetHostedZoneInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_zone")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *dns.GetHostedZoneInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &dns.GetHostedZoneInput{
				ZoneID: nifcloud.String("test_zone"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetHostedZoneInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteHostedZoneInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_zone")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *dns.DeleteHostedZoneInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &dns.DeleteHostedZoneInput{
				ZoneID: nifcloud.String("test_zone"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteHostedZoneInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
