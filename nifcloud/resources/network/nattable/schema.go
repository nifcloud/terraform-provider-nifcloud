package nattable

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a nat table resource."

// New returns the nifcloud_nat_table resource schema.
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
		"snat": {
			Type:        schema.TypeSet,
			Description: "A list of snat objects.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"rule_number": {
						Type:        schema.TypeString,
						Description: "The rule number.",
						Required:    true,
					},
					"description": {
						Type:             schema.TypeString,
						Description:      "The nat table rule description.",
						Optional:         true,
						ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
					},
					"protocol": {
						Type:        schema.TypeString,
						Description: "The protocol.",
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"ALL", "TCP", "UDP", "TCP_UDP", "ICMP",
						}, false),
					},
					"source_address": {
						Type:             schema.TypeString,
						Description:      "The source address.",
						Required:         true,
						ValidateDiagFunc: validator.IPAddress,
					},
					"source_port": {
						Type:         schema.TypeInt,
						Description:  "The source port.",
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"translation_port": {
						Type:         schema.TypeInt,
						Description:  "The translation port.",
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"outbound_interface_network_id": {
						Type:        schema.TypeString,
						Description: "The outbound interface network id; `net-COMMON_GLOBAL` or `net-COMMON_PRIVATE` or private lan network id.",
						Optional:    true,
					},
					"outbound_interface_network_name": {
						Type:        schema.TypeString,
						Description: "The private lan name of target outbound interface network.",
						Optional:    true,
					},
				},
			},
		},
		"dnat": {
			Type:        schema.TypeSet,
			Description: "A list of snat objects.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"rule_number": {
						Type:        schema.TypeString,
						Description: "The rule number.",
						Required:    true,
					},
					"description": {
						Type:             schema.TypeString,
						Description:      "The nat table rule description.",
						Optional:         true,
						ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
					},
					"protocol": {
						Type:        schema.TypeString,
						Description: "The protocol.",
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"ALL", "TCP", "UDP", "TCP_UDP", "ICMP",
						}, false),
					},
					"destination_port": {
						Type:         schema.TypeInt,
						Description:  "The destination port.",
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"translation_address": {
						Type:             schema.TypeString,
						Description:      "The translation address.",
						Required:         true,
						ValidateDiagFunc: validator.IPAddress,
					},
					"translation_port": {
						Type:         schema.TypeInt,
						Description:  "The translation port.",
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"inbound_interface_network_id": {
						Type:        schema.TypeString,
						Description: "The inbound interface network id; `net-COMMON_GLOBAL` or `net-COMMON_PRIVATE` or private lan network id.",
						Optional:    true,
					},
					"inbound_interface_network_name": {
						Type:        schema.TypeString,
						Description: "The private lan name of target inbound interface network.",
						Optional:    true,
					},
				},
			},
		},
		"nat_table_id": {
			Type:        schema.TypeString,
			Description: "The nat table id.",
			Computed:    true,
		},
	}
}
