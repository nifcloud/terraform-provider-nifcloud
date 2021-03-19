package webproxy

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a web proxy resource."

// New returns the nifcloud_web_proxy resource schema.
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
		"router_name": {
			Type:          schema.TypeString,
			Description:   "The name for the router. route_id and either is required.",
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"router_id"},
			Computed:      true,
		},
		"router_id": {
			Type:          schema.TypeString,
			Description:   "The id for the router. router_name and either is required.",
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"router_name"},
			Computed:      true,
		},
		"listen_interface_network_name": {
			Type:          schema.TypeString,
			Description:   "The name for the listen network. listen_interface_network_id and either is required.",
			Optional:      true,
			ConflictsWith: []string{"listen_interface_network_id"},
			Computed:      true,
		},
		"listen_interface_network_id": {
			Type:          schema.TypeString,
			Description:   "The id for the listen network. listen_interface_network_name and either is required.",
			Optional:      true,
			ConflictsWith: []string{"listen_interface_network_name"},
			Computed:      true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The web proxy description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
		"listen_port": {
			Type:        schema.TypeString,
			Description: "The port of web proxy.",
			Required:    true,
		},
		"bypass_interface_network_name": {
			Type:          schema.TypeString,
			Description:   "The name for the by pass network",
			Optional:      true,
			ConflictsWith: []string{"bypass_interface_network_id"},
			Computed:      true,
		},
		"bypass_interface_network_id": {
			Type:          schema.TypeString,
			Description:   "The id for the by pass network.",
			Optional:      true,
			ConflictsWith: []string{"bypass_interface_network_name"},
			Computed:      true,
		},
		"name_server": {
			Type:        schema.TypeString,
			Description: "The ip address for dns server.",
			Optional:    true,
			Computed:    true,
		},
	}
}
