package devopsfirewallgroup

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
	"golang.org/x/exp/slices"
)

const description = "Provides a DevOps firewall group resource."

// New returns the nifcloud_devops_firewall_group resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: createFirewallGroup,
		ReadContext:   readFirewallGroup,
		UpdateContext: updateFirewallGroup,
		DeleteContext: deleteFirewallGroup,

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
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps firewall group.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 63),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z\-]+$`), "Enter a name within 1-63 alphanumeric lowercase characters and hyphens. Hyphens cannot be used at the beginning or end."),
			),
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone for the DevOps firewall group.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the DevOps firewall group.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"rule": {
			Type:        schema.TypeSet,
			Description: "List of the DevOps firewall rules.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Description: "ID of the rule.",
						Computed:    true,
					},
					"protocol": {
						Type:        schema.TypeString,
						Description: "Protocol. Valid values are `TCP` or `ICMP`.",
						Optional:    true,
						ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
							value := v.(string)
							var diags diag.Diagnostics
							if !slices.Contains([]string{"TCP", "ICMP"}, value) {
								diag := diag.Diagnostic{
									Severity: diag.Error,
									Summary:  "wrong value",
									Detail:   fmt.Sprintf("%q is neither \"TCP\" nor \"ICMP\".", value),
								}
								diags = append(diags, diag)
							}
							return diags
						},
					},
					"port": {
						Type:        schema.TypeInt,
						Description: "Port. Valid values are `22` or `443`.",
						Optional:    true,
						ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
							value := v.(int)
							var diags diag.Diagnostics
							if value != 22 && value != 443 {
								diag := diag.Diagnostic{
									Severity: diag.Error,
									Summary:  "wrong value",
									Detail:   fmt.Sprintf("%d is neither 22 nor 443.", value),
								}
								diags = append(diags, diag)
							}
							return diags
						},
					},
					"cidr_ip": {
						Type:        schema.TypeString,
						Description: "CIDR block or IPv4 address.",
						Optional:    true,
						ValidateDiagFunc: validator.Any(
							validator.CIDRNetworkAddress,
							validator.IPAddress,
						),
					},
					"description": {
						Type:             schema.TypeString,
						Description:      "Description of the rule.",
						Optional:         true,
						ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
					},
				},
			},
		},
	}
}
