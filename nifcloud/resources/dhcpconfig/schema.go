package dhcpconfig

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a dhcp config resource."

// New returns the nifcloud_dhcp_config resource schema.
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
			Delete:  schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"static_mapping": {
			Type:        schema.TypeSet,
			Description: "A list of static mapping ip address.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"static_mapping_ipaddress": {
						Type:             schema.TypeString,
						Description:      "The static mapping IP address.",
						Required:         true,
						ValidateDiagFunc: validator.IPAddress,
					},
					"static_mapping_macaddress": {
						Type:        schema.TypeString,
						Description: "The static mapping MAC address.",
						Required:    true,
					},
					"static_mapping_description": {
						Type:        schema.TypeString,
						Description: "The static mapping IP address description.",
						Optional:    true,
					},
				},
			},
		},
		"ipaddress_pool": {
			Type:        schema.TypeSet,
			Description: "A list of ipaddress pool.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ipaddress_pool_start": {
						Type:             schema.TypeString,
						Description:      "The start IP address of ipAddressPool.",
						Required:         true,
						ValidateDiagFunc: validator.IPAddress,
					},
					"ipaddress_pool_stop": {
						Type:             schema.TypeString,
						Description:      "The stop IP address of ipAddressPool.",
						Required:         true,
						ValidateDiagFunc: validator.IPAddress,
					},
					"ipaddress_pool_description": {
						Type:        schema.TypeString,
						Description: "The ipaddress pool description.",
						Optional:    true,
					},
				},
			},
		},
		"dhcp_config_id": {
			Type:        schema.TypeString,
			Description: "The ID of the dhcp config.",
			Computed:    true,
		},
	}
}
