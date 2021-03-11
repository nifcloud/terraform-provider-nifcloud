package networkinterface

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides an additional nic resource."

// New returns the nifcloud_network_interface resource schema.
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
			Description: "The availability zone.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "A description for the network interface.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 500),
		},
		"network_id": {
			Type:        schema.TypeString,
			Description: "Private lan ID to create the NIC in.",
			Required:    true,
			ForceNew:    true,
		},
		"ip_address": {
			Type:        schema.TypeString,
			Description: "If DHCP is enabled, specify IP address or `static` or not specified(by DHCP). Otherwise, specify `static`.",
			Optional:    true,
		},
		"private_ip": {
			Type:        schema.TypeString,
			Description: "Private IP address of network interface.",
			Computed:    true,
		},
		"network_interface_id": {
			Type:        schema.TypeString,
			Description: "The ID of network interface.",
			Computed:    true,
		},
	}
}
