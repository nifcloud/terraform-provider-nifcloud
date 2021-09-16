package cluster

import (
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
			Type:        schema.TypeString,
			Description: "The name of the cluster.",
			Required:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The cluster description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Description: "The version of Kubernetes.",
			Optional:    true,
			Computed:    true,
		},
		"kube_config_raw": {
			Type:        schema.TypeString,
			Description: "The raw Kubernetes config to be used by kubectl and other compatible tools.",
			Computed:    true,
			Sensitive:   true,
		},
		"locations": {
			Type:        schema.TypeList,
			Description: "The cluster location. availability zone can be specified.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: true,
			MinItems: 1,
			ForceNew: true,
		},
		"addons_config": {
			Type:        schema.TypeList,
			Description: "The configs for Kubernetes addons.",
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"http_load_balancing": {
						Type:        schema.TypeList,
						Description: "The configs for HTTP load balancer.",
						Optional:    true,
						Computed:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"disabled": {
									Type:        schema.TypeBool,
									Description: "Disable the HTTP load balancing addon.",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
		"network_config": {
			Type:        schema.TypeList,
			Description: "The configs for cluster network.",
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network_id": {
						Type:        schema.TypeString,
						Description: "The ID of private LAN.",
						Required:    true,
						ForceNew:    true,
					},
				},
			},
		},
		"firewall_group": {
			Type:        schema.TypeString,
			Description: "The firewall group name to associate with.",
			Required:    true,
			ForceNew:    true,
		},
		"node_pools": {
			Type:        schema.TypeSet,
			Description: "The node pool config",
			Required:    true,
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "The name of the node pool.",
						Required:    true,
					},
					"instance_type": {
						Type:        schema.TypeString,
						Description: "The instance type for node pool.",
						Required:    true,
					},
					"node_count": {
						Type:        schema.TypeInt,
						Description: "The desired node count in this node pool.",
						Required:    true,
					},
					"nodes": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Description: "The name of the node.",
									Computed:    true,
								},
								"availability_zone": {
									Type:        schema.TypeString,
									Description: "The availability zone where the node located.",
									Computed:    true,
								},
								"public_ip_address": {
									Type:        schema.TypeString,
									Description: "The public IP address of the node.",
									Computed:    true,
								},
								"private_ip_address": {
									Type:        schema.TypeString,
									Description: "The private IP address of the node.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}
