package router

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a router resource."

// New returns the nifcloud_router resource schema.
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

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"accounting_type": {
			Type:         schema.TypeString,
			Description:  "Accounting type. (1: monthly, 2: pay per use).",
			Optional:     true,
			Default:      "2",
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The router description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 500),
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The router name.",
			Optional:    true,
			Computed:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"nat_table_association_id": {
			Type:        schema.TypeString,
			Description: "The ID of the NAT table association.",
			Computed:    true,
		},
		"nat_table_id": {
			Type:        schema.TypeString,
			Description: "The ID of the NAT table to attach.",
			Optional:    true,
		},
		"network_interface": {
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"dhcp": {
						Type:        schema.TypeBool,
						Description: "The flag to enable or disable DHCP.",
						Optional:    true,
					},
					"dhcp_config_id": {
						Type:        schema.TypeString,
						Description: "The ID of the DHCP config to attach.",
						Optional:    true,
					},
					"dhcp_options_id": {
						Type:        schema.TypeString,
						Description: "The ID of the DHCP options to attach.",
						Optional:    true,
					},
					"ip_address": {
						Type:         schema.TypeString,
						Description:  "The IP address of the network interface.",
						Optional:     true,
						ValidateFunc: validation.IsIPAddress,
					},
					"network_id": {
						Type:        schema.TypeString,
						Description: "The ID of the network to attach; 'net-COMMON_GLOBAL' or `net-COMMON_PRIVATE` or private lan network id.",
						Optional:    true,
					},
					"network_name": {
						Type:        schema.TypeString,
						Description: "The private lan name of the network to attach.",
						Optional:    true,
					},
				},
			},
		},
		"router_id": {
			Type:        schema.TypeString,
			Description: "The unique ID of the router.",
			Computed:    true,
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
		"security_group": {
			Type:        schema.TypeString,
			Description: "The security group name to associate with; which can be managed using the nifcloud_security_group resource.",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the router.",
			Optional:    true,
			Default:     "small",
			ValidateFunc: validation.StringInSlice([]string{
				"small",
				"medium",
				"large",
			}, false),
		},
	}
}
