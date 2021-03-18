package routetable

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a route table resource."

// New returns the nifcloud_route_table resource schema.
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
		"route": {
			Type:        schema.TypeSet,
			Description: "A list of route objects.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cidr_block": {
						Type:        schema.TypeString,
						Description: "The destination IP address or CIDR.",
						Required:    true,
						ValidateDiagFunc: validator.Any(
							validator.CIDRNetworkAddress,
							validator.IPAddress,
						),
					},
					"ip_address": {
						Type:             schema.TypeString,
						Description:      "The target IP address.",
						Optional:         true,
						ValidateDiagFunc: validator.IPAddress,
					},
					"network_id": {
						Type:        schema.TypeString,
						Description: "The id of target network; 'net-COMMON_GLOBAL' or `net-COMMON_PRIVATE` or private lan network id.",
						Optional:    true,
					},
					"network_name": {
						Type:        schema.TypeString,
						Description: "The private lan name of target network.",
						Optional:    true,
					},
				},
			},
		},
		"route_table_id": {
			Type:        schema.TypeString,
			Description: "The id of route table.",
			Computed:    true,
		},
	}
}
