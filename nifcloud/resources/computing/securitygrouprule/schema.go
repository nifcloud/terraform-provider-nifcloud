package securitygrouprule

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a security group rule resource. Represents a single in or out group rule, which can be added to external Security Groups."

// New returns the nifcloud_security_group_rule resource schema.
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
				importParts, err := validateSecurityGroupRuleImportString(d.Id())
				if err != nil {
					return nil, err
				}
				if err := populateSecurityGroupRuleFromImport(d, importParts); err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Description:  "The type of rule being created. Valid options are IN (Incoming) or OUT (Outgoing).",
			Optional:     true,
			Default:      "IN",
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"IN", "OUT"}, false),
		},
		"cidr_ip": {
			Type:        schema.TypeString,
			Description: "The CIDR IP Address. Cannot be specified with `source_security_group_name` .",
			Optional:    true,
			ForceNew:    true,
			ValidateDiagFunc: validator.Any(
				validator.CIDRNetworkAddress,
				validator.IPAddress,
			),
		},
		"from_port": {
			Type:         schema.TypeInt,
			Description:  "The start port",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 65535),
		},
		"protocol": {
			Type:        schema.TypeString,
			Description: "The protocol.",
			Optional:    true,
			Default:     "TCP",
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"ANY", "TCP", "UDP", "ICMP", "GRE", "ESP", "AH", "VRRP", "ICMPv6-all",
			}, false),
		},
		"security_group_names": {
			Type:        schema.TypeList,
			Description: "The security group name list to apply this rule.",
			Required:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 15),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
				),
			},
		},
		"source_security_group_name": {
			Type:          schema.TypeString,
			Description:   "The security group name that allow access. Cannot be specified with `cidr_ip` .",
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"cidr_ip"},
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"to_port": {
			Type:         schema.TypeInt,
			Description:  "The end port",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 65535),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The security group rule description.",
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
	}
}
