package bucket

import (
	"strings"
	"time"

	awspolicy "github.com/hashicorp/awspolicyequivalence"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const description = "Provides a storage bucket resource."

// New returns the nifcloud_storage_bucket resource schema.
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
		"bucket": {
			Type:        schema.TypeString,
			Description: "The name of the bucket.",
			Required:    true,
			ForceNew:    true,
		},
		"versioning": {
			Type:        schema.TypeList,
			Description: "A configuration of the bucket versioning state.",
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Description: "Enable versioning.",
						Required:    true,
					},
				},
			},
		},
		"policy": {
			Type:             schema.TypeString,
			Description:      "A bucket policy JSON document.",
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: suppressEquivalentPolicyDiffs,
		},
	}
}

func suppressEquivalentPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	if strings.TrimSpace(old) == "" && strings.TrimSpace(new) == "" {
		return true
	}

	if strings.TrimSpace(old) == "{}" && strings.TrimSpace(new) == "" {
		return true
	}

	if strings.TrimSpace(old) == "" && strings.TrimSpace(new) == "{}" {
		return true
	}

	if strings.TrimSpace(old) == "{}" && strings.TrimSpace(new) == "{}" {
		return true
	}

	equivalent, err := awspolicy.PoliciesAreEquivalent(old, new)
	if err != nil {
		return false
	}

	return equivalent
}
