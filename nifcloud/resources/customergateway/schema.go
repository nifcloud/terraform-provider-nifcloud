package customergateway

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a customer gateway resource."

// New returns the nifcloud_customer_gateway resource schema.
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
		"customer_gateway_id": {
			Type:        schema.TypeString,
			Description: "The customer gateway id.",
			Computed:    true,
		},
		"customer_gateway_name": {
			Type:        schema.TypeString,
			Description: "The customer gateway name.",
			Optional:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"customer_gateway_description": {
			Type:             schema.TypeString,
			Description:      "The customer gateway description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(1, 500),
		},
		"ip_address": {
			Type:             schema.TypeString,
			Description:      "The IP address.",
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validator.IPAddress,
		},
		"lan_side_ip_address": {
			Type:             schema.TypeString,
			Description:      "The lan side IP address.",
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validator.IPAddress,
		},
		"lan_side_cidr_block": {
			Type:             schema.TypeString,
			Description:      "The lan side CIDR block.",
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validator.CIDRNetworkAddress,
		},
		"type": {
			Type:         schema.TypeString,
			Description:  "The type.",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"IPsec", "IPsec VTI", "L2TPv3 / IPsec"}, false),
		},
	}
}
