package cluster

import (
	"hash/crc32"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a Kubernetes Service Hatoba cluster resource."

// New returns the nifcloud_hatoba_cluster resource schema.
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
			Create:  schema.DefaultTimeout(120 * time.Minute),
			Update:  schema.DefaultTimeout(120 * time.Minute),
			Delete:  schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The firewall group description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"kubernetes_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"locations": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: true,
			MinItems: 1,
			ForceNew: true,
		},
		"addons_config": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"http_load_balancing": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"disabled": {
									Type:     schema.TypeBool,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"network_config": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network_id": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},
		"firewall_group": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"node_pools": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Set: func(v interface{}) int {
				name := v.(map[string]interface{})["name"].(string)
				return int(crc32.ChecksumIEEE([]byte(name)))
			},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"instance_type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"node_count": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"nodes": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"availability_zone": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"public_ip_address": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_ip_address": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}
