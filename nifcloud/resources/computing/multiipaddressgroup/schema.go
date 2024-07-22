package multiipaddressgroup

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a multi IP address group resource."

// New returns the nifcloud_multi_ip_address_group resource schema.
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
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone name.",
			Required:    true,
			ForceNew:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The description of the multi IP address group.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 500),
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the multi IP address group.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), "Enter the name within 1-15 characters [0-9a-zA-Z]."),
			),
		},
		"ip_address_count": {
			Type:         schema.TypeInt,
			Description:  "The number of IP addresses to secure.",
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"default_gateway": {
			Type:        schema.TypeString,
			Description: "The default gateway of the multi IP address network.",
			Computed:    true,
		},
		"subnet_mask": {
			Type:        schema.TypeString,
			Description: "The subnet mask of the multi IP address network.",
			Computed:    true,
		},
		"ip_addresses": {
			Type:        schema.TypeList,
			Description: "A list of the secured IP addresses for the multi IP address network.",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}
