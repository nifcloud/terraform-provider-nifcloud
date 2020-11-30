package elasticip

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a elastic ip resource."

// New returns the nifcloud_elastic_ip resource schema.
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
		"ip_type": {
			Type:        schema.TypeBool,
			Description: "Choice of the private ip address(true) or public ip address(false).",
			Required:    true,
			ForceNew:    true,
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone.",
			Required:    true,
			ForceNew:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The elastic ip description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
		"private_ip": {
			Type:        schema.TypeString,
			Description: "The private ip address.",
			Computed:    true,
		},
		"public_ip": {
			Type:        schema.TypeString,
			Description: "The public ip address.",
			Computed:    true,
		},
	}
}
