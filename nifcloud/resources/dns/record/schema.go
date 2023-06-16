package record

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a dns record resource."

// New returns the nifcloud_dns_record resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: create,
		ReadContext:   read,
		DeleteContext: delete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				importParts, err := validateDnsRecordImportString(d.Id())
				if err != nil {
					return nil, err
				}
				if err := populateDnsRecordFromImport(d, importParts); err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"zone_id": {
			Type:        schema.TypeString,
			Description: "The ID of the hosted zone to contain this record.",
			Required:    true,
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the record.",
			Required:    true,
			ForceNew:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the record.",
			Required:    true,
			ForceNew:    true,
		},
		"record": {
			Type:        schema.TypeString,
			Description: "The value of the record.",
			Required:    true,
			ForceNew:    true,
		},
		"ttl": {
			Type:         schema.TypeInt,
			Description:  "The TTL of the record.",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(60, 86400),
			Default:      3600,
		},
		"weighted_routing_policy": {
			Type:          schema.TypeList,
			Description:   "The configs for weighted routing policy. Conflicts with failover_routing_policy.",
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"failover_routing_policy"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"weight": {
						Type:        schema.TypeInt,
						Description: "The record weighted value.",
						Optional:    true,
						ForceNew:    true,
					},
				},
			},
		},
		"failover_routing_policy": {
			Type:          schema.TypeList,
			Description:   "The configs for failover routing policy. Conflicts with weighted_routing_policy.",
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"weighted_routing_policy"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:        schema.TypeString,
						Description: "The record failover type.",
						Optional:    true,
						ForceNew:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"PRIMARY", "SECONDARY"}, false),
					},
					"health_check": {
						Type:        schema.TypeList,
						Description: "The configs for health check if using failover.",
						Optional:    true,
						ForceNew:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"protocol": {
									Type:        schema.TypeString,
									Description: "The health check protocol.",
									Optional:    true,
									ForceNew:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"HTTP", "HTTPS", "TCP"}, false),
								},
								"ip_address": {
									Type:             schema.TypeString,
									Description:      "The health check IP address.",
									Optional:         true,
									ForceNew:         true,
									ValidateDiagFunc: validator.StringRuneCountBetween(1, 32),
								},
								"port": {
									Type:         schema.TypeInt,
									Description:  "The health check port.",
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},
								"resource_path": {
									Type:             schema.TypeString,
									Description:      "The health check resource path if using HTTP or HTTPS protocol.",
									Optional:         true,
									ForceNew:         true,
									ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
								},
								"resource_domain": {
									Type:             schema.TypeString,
									Description:      "The health check resource domain if using HTTP or HTTPS protocol.",
									Optional:         true,
									ForceNew:         true,
									ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
								},
							},
						},
					},
				},
			},
		},
		"comment": {
			Type:             schema.TypeString,
			Description:      "The comment of the record.",
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"set_identifier": {
			Type:        schema.TypeString,
			Description: "The unique identifier to differentiate records with routing policies from one another.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
	}
}
