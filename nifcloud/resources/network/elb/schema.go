package elb

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a multi load balancer resource."

// New returns the nifcloud_elb resource schema.
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
				importParts, err := validateELBImportString(d.Id())
				if err != nil {
					return nil, err
				}
				if err := populateELBFromImport(d, importParts); err != nil {
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
		"elb_name": {
			Type:        schema.TypeString,
			Description: "The name for the multi load balancer.",
			Optional:    true,
			Computed:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), "Enter the elb_name within 1-15 characters [0-9a-zA-Z]."),
			),
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone.",
			Required:    true,
			ForceNew:    true,
		},
		"accounting_type": {
			Type:         schema.TypeString,
			Description:  "Accounting type. (1: monthly, 2: pay per use).",
			Optional:     true,
			Default:      "2",
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
		},
		"network_volume": {
			Type:         schema.TypeInt,
			Description:  "Maximum network volume for the multi load balancer.",
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntInSlice([]int{10, 20, 30, 40, 100, 200, 300, 400, 500}),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The multi load balancer description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 500),
		},
		"balancing_type": {
			Type:         schema.TypeInt,
			Description:  "Balancing type. (1: Round-Robin, 2: Least-Connection).",
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntInSlice([]int{1, 2}),
		},
		"instance_port": {
			Type:         schema.TypeInt,
			Description:  "The port on the instance to route to.",
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"protocol": {
			Type:         schema.TypeString,
			Description:  "The protocol to listen on. Valid values are `HTTP` `HTTPS` `TCP` `UDP`.",
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "HTTP", "HTTPS"}, false),
		},
		"lb_port": {
			Type:         schema.TypeInt,
			Description:  "The port to listen on for the multi load balancer.",
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"ssl_certificate_id": {
			Type:        schema.TypeString,
			Description: "The id of the SSL certificate you have uploaded to NIFCLOUD.",
			Optional:    true,
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
		"health_check_path": {
			Type:        schema.TypeString,
			Description: "The path of the health check.",
			Optional:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(0, 255),
				validation.StringMatch(
					regexp.MustCompile(`^[/][\w/:%&~='<>@\?\(\)\.\,\+\-\*\[\]\^\{\}\|]*$`),
					"Enter the health_check_path within 0-255 characters",
				),
			),
		},
		"health_check_expectation_http_code": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "A list of the expected http code.",
			Optional:    true,
			MaxItems:    10,
		},
		"instances": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "A list of instance names to place in the multi load balancer pool.",
			Optional:    true,
		},
		"network_interface": {
			Type:     schema.TypeSet,
			Required: true,
			ForceNew: true,
			MaxItems: 2,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
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
					"ip_address": {
						Type:        schema.TypeString,
						Description: "The IP address of multi load balancer.",
						Optional:    true,
					},
					"is_vip_network": {
						Type:        schema.TypeBool,
						Description: "The flag of vip network.",
						Default:     true,
						Optional:    true,
					},
				},
			},
		},
		"session_stickiness_policy_enable": {
			Type:        schema.TypeBool,
			Description: "The flag of session stickiness policy.",
			Default:     false,
			Optional:    true,
		},
		"session_stickiness_policy_method": {
			Type:         schema.TypeString,
			Description:  "The session stickiness policy method. (1: Source ip, 2: Cookie)",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
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
		"sorry_page_redirect_url": {
			Type:         schema.TypeString,
			Description:  "The sorry page redirect url.",
			Optional:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(https?):\/\/.+$`), "Invalid format for a sorry_page_redirect_url"),
		},
		"route_table_id": {
			Type:        schema.TypeString,
			Description: "The id of route table to attach.",
			Optional:    true,
		},
		"route_table_association_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The id of route table association.",
		},
		"dns_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The ip address of multi load balancer vip network.",
		},
		"elb_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The id of multi load balancer.",
		},
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The version of multi load balancer.",
		},
	}
}
