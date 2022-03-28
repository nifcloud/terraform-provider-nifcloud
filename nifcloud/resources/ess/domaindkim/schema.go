package domaindkim

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const description = "Provides an ESS domain DKIM generation resource."

// New returns the nifcloud_ess_domain_dkim resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: create,
		ReadContext:   read,
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
		"domain": {
			Type:        schema.TypeString,
			Description: "Verified domain name to generate DKIM tokens for.",
			Required:    true,
			ForceNew:    true,
		},
		"dkim_tokens": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "DKIM tokens generated by ESS. These tokens should be used to create CNAME records used to verify ESS Easy DKIM. ",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}
