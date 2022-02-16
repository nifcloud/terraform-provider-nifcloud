package zone

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a dns zone resource."

// New returns the nifcloud_dns_zone resource schema.
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
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the hosted zone.",
			Required:    true,
			ForceNew:    true,
		},
		"comment": {
			Type:             schema.TypeString,
			Description:      "The comment of the hosted zone.",
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"name_servers": {
			Type:        schema.TypeList,
			Description: "The list of name servers.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Computed: true,
		},
	}
}
