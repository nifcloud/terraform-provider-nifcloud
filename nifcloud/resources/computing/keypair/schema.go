package keypair

import (
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Upload and register the specified SSH public key."

// New returns the nifcloud_key_pair resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key_name": {
			Type:        schema.TypeString,
			Description: "The name for the key pair.",
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(6, 32),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The key pair description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
		"public_key": {
			Type:         schema.TypeString,
			Description:  "The public key material.",
			ValidateFunc: validation.StringIsBase64,
			Required:     true,
			ForceNew:     true,
			StateFunc: func(v interface{}) string {
				switch v := v.(type) {
				case string:
					return strings.TrimSpace(v)
				default:
					return ""
				}
			},
		},
		"fingerprint": {
			Type:        schema.TypeString,
			Description: "The MD5 public key fingerprint.",
			Computed:    true,
		},
	}
}
