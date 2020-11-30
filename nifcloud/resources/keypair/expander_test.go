package keypair

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandImportKeyPairInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"key_name":    "test_key_name",
		"public_key":  "test_public_key",
		"description": "test_description",
	})
	rd.SetId("test_key_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ImportKeyPairInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ImportKeyPairInput{
				KeyName:           nifcloud.String("test_key_name"),
				PublicKeyMaterial: nifcloud.String("test_public_key"),
				Description:       nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandImportKeyPairInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyModifyKeyPairAttributeInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"key_name":    "test_key_name",
		"description": "test_description",
	})
	rd.SetId("test_key_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyModifyKeyPairAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyModifyKeyPairAttributeInput{
				KeyName:   nifcloud.String("test_key_name"),
				Attribute: "description",
				Value:     nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyModifyKeyPairAttributeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
