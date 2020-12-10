package loadbalancer

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const description = "Provide a load_balancer resource"

// New returns the nifcloud_load_balancer resource schema.
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
		"accounting_type": {
			Type:         schema.TypeString,
			Description:  "Accounting type. (1: monthly, 2: pay per use).",
			Optional:     true,
			Default:      "2",
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
		},
		"availability_zones": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			ForceNew:    true,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dns_name": {
			Type:        schema.TypeString,
			Description: "dns_name",
			Computed:    true,
		},
		"ip_version": {
			Type:         schema.TypeString,
			Description:  "The load balancer ip version(v4 or v6).",
			Optional:     true,
			ForceNew:     true,
			Default:      "v4",
			ValidateFunc: validation.StringInSlice([]string{"v4", "v6"}, false),
		},
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
		"network_volume": {
			Type:         schema.TypeInt,
			Description:  "The load balancer max network volume.",
			ValidateFunc: validation.IntInSlice([]int{10, 20, 30, 40, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 1100, 1200, 1300, 1400, 1500, 1600, 1700, 1800, 1900, 2000}),
			Optional:     true,
			Default:      10,
		},
		"policy_type": {
			Type:         schema.TypeString,
			Description:  "policy type. (standard or ats).",
			Optional:     true,
			ForceNew:     true,
			Default:      "standard",
			ValidateFunc: validation.StringInSlice([]string{"standard", "ats"}, false),
		},
	}
}
