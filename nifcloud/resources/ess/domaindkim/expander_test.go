package domaindkim

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
	"github.com/stretchr/testify/assert"
)

func TestExpandVerifyDomainDkimInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"domain": "test_domain",
	})
	rd.SetId("test_domain")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *ess.VerifyDomainDkimInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &ess.VerifyDomainDkimInput{
				Domain: nifcloud.String("test_domain"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandVerifyDomainDkimInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetIdentityDkimAttributesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"domain": "test_domain",
	})
	rd.SetId("test_domain")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *ess.GetIdentityDkimAttributesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &ess.GetIdentityDkimAttributesInput{
				Identities: []string{"test_domain"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetIdentityDkimAttributesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
