package privatelan

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

//const description = "Upload and register the specified SSH public key."
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
		"network_id": {
			Type:        schema.TypeString,
			Description: "The id for the private lan.",
			Computed:    true,
		},
		"private_lan_name": {
			Type:        schema.TypeString,
			Description: "The name for the private lan.",
			Optional:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"cidr_block": {
			Type:        schema.TypeString,
			Description: "The CIDR IP Address.",
			Required:    true,
			ValidateDiagFunc: validator.Any(
				validator.CIDRNetworkAddress,
			),
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "availability zone",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"accounting_type": {
			Type:         schema.TypeString,
			Description:  "accounting type",
			Optional:     true,
			Default:      "2",
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The private lan description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
		"state": {
			Type:        schema.TypeString,
			Description: "The state of the private lan.",
			Computed:    true,
		},
	}
}
