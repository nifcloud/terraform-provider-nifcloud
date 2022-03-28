package emailidentity

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
	"github.com/stretchr/testify/assert"
)

func TestExpandVerifyEmailIdentityInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"email": "test_email",
	})
	rd.SetId("test_email")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *ess.VerifyEmailIdentityInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &ess.VerifyEmailIdentityInput{
				EmailAddress: nifcloud.String("test_email"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandVerifyEmailIdentityInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetIdentityVerificationAttributesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"email": "test_email",
	})
	rd.SetId("test_email")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *ess.GetIdentityVerificationAttributesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &ess.GetIdentityVerificationAttributesInput{
				Identities: []string{"test_email"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetIdentityVerificationAttributesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteIdentityInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"email": "test_email",
	})
	rd.SetId("test_email")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *ess.DeleteIdentityInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &ess.DeleteIdentityInput{
				Identity: nifcloud.String("test_email"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteIdentityInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
