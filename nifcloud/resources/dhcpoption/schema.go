package dhcpoption

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a dhcp options resource."

// New returns the nifcloud_dhcp_option resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: create,
		ReadContext:   read,
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
		"default_router": {
			Type:        schema.TypeString,
			Description: "The IP address of default gateway.",
			Optional:    true,
			ForceNew:    true,
			ValidateDiagFunc: validator.Any(
				validator.IPAddress,
			),
		},
		"domain_name": {
			Type:        schema.TypeString,
			Description: "The domain name used by the client in host name resolution.",
			Optional:    true,
			ForceNew:    true,
		},
		"domain_name_servers": {
			Type:        schema.TypeList,
			Description: "The IP address of the DNS server.",
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"ntp_servers": {
			Type:        schema.TypeList,
			Description: "The IP address of the NTP server.",
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"netbios_name_servers": {
			Type:        schema.TypeList,
			Description: "The IP address of the NetBIOS server.",
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"netbios_node_type": {
			Type:        schema.TypeString,
			Description: "The NetBIOS node type. (1: Don't use WINS, 2: Don't use broadcast, 4: Priorirtize broadcasting, 8: Prioritize WINS)",
			Optional:    true,
			ForceNew:    true,
		},
		"lease_time": {
			Type:        schema.TypeString,
			Description: "The IP address lease time. (Unit: second)",
			Optional:    true,
			ForceNew:    true,
		},
		"dhcp_option_id": {
			Type:        schema.TypeString,
			Description: "The ID of the dhcp option.",
			Computed:    true,
		},
	}
}
