package vpngateway

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a vpn gateway resource."

// New returns the nifcloud_vpn_gateway resource schema.
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
			Create:  schema.DefaultTimeout(30 * time.Minute),
			Update:  schema.DefaultTimeout(20 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

// TODO schemaの定義順をexpanderに合わせる

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vpn_gateway_id": {
			Type:        schema.TypeString,
			Description: "The id for the vpn gateway.",
			Computed:    true,
		},
		"nifty_vpn_gateway_name": {
			Type:        schema.TypeString,
			Description: "The name for the vpn gateway.",
			Optional:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"nifty_vpn_gateway_type": {
			Type:         schema.TypeString,
			Description:  "The type of vpn gateway.",
			Optional:     true,
			Default:      "small",
			ValidateFunc: validation.StringInSlice([]string{"small", "medium", "large"}, false),
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone.",
			Optional:    true,
		},
		"accounting_type": {
			Type:         schema.TypeString,
			Description:  "The accounting type.",
			Optional:     true,
			Default:      "2",
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
		},
		"nifty_vpn_gateway_description": {
			Type:             schema.TypeString,
			Description:      "The vpn gateway description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
		"network_id": {
			Type:        schema.TypeString,
			Description: "The id for the network.",
			Optional:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"network_name": {
			Type:        schema.TypeString,
			Description: "The name for the network.",
			Optional:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"ip_address": {
			Type:             schema.TypeString,
			Description:      "The private ip address.",
			Optional:         true,
			ValidateDiagFunc: validator.IPAddress,
		},
		"security_group": {
			Type:        schema.TypeString,
			Description: "The name of firewall group.",
			Optional:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"route_table_association_id": {
			Type:        schema.TypeString,
			Description: "The ID of the route table association.",
			Computed:    true,
		},
		"route_table_id": {
			Type:        schema.TypeString,
			Description: "The ID of the route table to attach.",
			Optional:    true,
		},
	}
}
