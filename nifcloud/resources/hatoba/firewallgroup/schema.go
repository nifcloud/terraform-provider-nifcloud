package firewallgroup

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a Kubernetes Service Hatoba firewall group resource."

// New returns the nifcloud_hatoba_firewall_group resource schema.
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
		"rule": {
			Type:        schema.TypeSet,
			Description: "A list of the firewall group rule objects.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"protocol": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The protocol.",
						ValidateFunc: validation.StringInSlice([]string{
							"ANY", "TCP", "UDP", "ICMP", "GRE", "ESP", "AH", "VRRP", "ICMPv6-all",
						}, false),
					},
					"direction": {
						Type:         schema.TypeString,
						Optional:     true,
						Description:  "The direction of rule being created.",
						ValidateFunc: validation.StringInSlice([]string{"IN", "OUT"}, false),
					},
					"from_port": {
						Type:         schema.TypeInt,
						Optional:     true,
						Description:  "The start port.",
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"to_port": {
						Type:         schema.TypeInt,
						Optional:     true,
						Description:  "The end port.",
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"cidr_ip": {
						Type:        schema.TypeString,
						Description: "The CIDR IP address that allow access.",
						Optional:    true,
						ValidateDiagFunc: validator.Any(
							validator.CIDRNetworkAddress,
							validator.IPAddress,
						),
					},
					"description": {
						Type:             schema.TypeString,
						Optional:         true,
						Description:      "The firewall group rule description.",
						ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
					},
					"id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The unique identifier of rule.",
					},
				},
			},
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name for the firewall group.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 40),
			),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The firewall group description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
	}
}
