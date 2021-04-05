package loadbalancerlistener

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const description = "Provide a load_balancer_listener resource"

// New returns the nifcloud_load_balancer_listener resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: delete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				importParts, err := validateLBImportString(d.Id())
				if err != nil {
					return nil, err
				}
				if err := populateLBFromImport(d, importParts); err != nil {
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
		"load_balancer_name": {
			Type:        schema.TypeString,
			Description: "The name for the load_balancer.",
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"filter": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "A list of IP address filter for load balancer",
			Optional:    true,
		},
		"filter_type": {
			Type:         schema.TypeString,
			Description:  "The filter_type of filter (1: Allow, 2: Deny).",
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
		},
		"healthy_threshold": {
			Type:        schema.TypeInt,
			Description: "The number of checks before the instance is declared healthy.",
			Optional:    true,
			Default:     1,
		},
		"unhealthy_threshold": {
			Type:         schema.TypeInt,
			Description:  "The number of checks before the instance is declared unhealthy.",
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 10),
		},
		"health_check_target": {
			Type:        schema.TypeString,
			Description: "The target of the health check. Valid pattern is ${PROTOCOL}:${PORT} or ICMP.",
			Optional:    true,
			Default:     "ICMP",
		},
		"health_check_interval": {
			Type:         schema.TypeInt,
			Description:  "The interval between health checks.",
			Optional:     true,
			Default:      5,
			ValidateFunc: validation.IntBetween(5, 300),
		},
		"instances": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "A list of instance names to place in the load balancer pool.",
			Optional:    true,
		},
		"instance_port": {
			Type:         schema.TypeInt,
			Description:  "The port on the instance to route to.",
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"load_balancer_port": {
			Type:         schema.TypeInt,
			Description:  "The port to listen on for the load balancer.",
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"balancing_type": {
			Type:         schema.TypeInt,
			Description:  "Balancing type. (1: Round-Robin, 2: Least-Connection).",
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntInSlice([]int{1, 2}),
		},
		"ssl_certificate_id": {
			Type:        schema.TypeString,
			Description: "The id of the SSL certificate you have uploaded to NIFCLOUD.",
			Optional:    true,
		},
		"ssl_policy_id": {
			Type:        schema.TypeString,
			Description: "The id of the SSL policy.",
			Optional:    true,
			Computed:    true,
		},
		"ssl_policy_name": {
			Type:        schema.TypeString,
			Description: "The name of the SSL policy.",
			Optional:    true,
			Computed:    true,
		},
		"session_stickiness_policy_enable": {
			Type:        schema.TypeBool,
			Description: "The flag of session stickiness policy.",
			Default:     false,
			Optional:    true,
		},
		"session_stickiness_policy_expiration_period": {
			Type:         schema.TypeInt,
			Description:  "The session stickiness policy expiration period.",
			Optional:     true,
			ValidateFunc: validation.IntBetween(3, 60),
		},
		"sorry_page_enable": {
			Type:        schema.TypeBool,
			Description: "The flag of sorry page.",
			Default:     false,
			Optional:    true,
		},
		"sorry_page_status_code": {
			Type:        schema.TypeInt,
			Description: "The HTTP status code for sorry page.",
			Optional:    true,
		},
		"ip_version": {
			Type:        schema.TypeString,
			Description: "The load balancer ip version(v4 or v6).",
			Computed:    true,
		},
		"policy_type": {
			Type:         schema.TypeString,
			Description:  "policy type (standard or ats).",
			Optional:     true,
			Default:      "standard",
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"standard", "ats"}, false),
		},
	}
}
