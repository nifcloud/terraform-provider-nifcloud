package nassecuritygroup

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a NAS security group resource."

// New returns the nifcloud_nas_security_group resource schema.
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
		"rule": {
			Type:        schema.TypeSet,
			Description: "A list of the NAS security group rule objects.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cidr_ip": {
						Type:        schema.TypeString,
						Description: "The CIDR IP Address that allow access. Cannot be specified with `security_group_name` .",
						Optional:    true,
						ValidateDiagFunc: validator.Any(
							validator.CIDRNetworkAddress,
							validator.IPAddress,
						),
					},
					"security_group_name": {
						Type:        schema.TypeString,
						Description: "The security group name that allow access. Cannot be specified with `cidr_ip` .",
						Optional:    true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(1, 15),
							validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), "Enter the security_group_name within 1-15 characters [0-9a-zA-Z]."),
						),
					},
				},
			},
		},
		"group_name": {
			Type:        schema.TypeString,
			Description: "The name for the NAS security group.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]+[0-9a-zA-Z_-]*$`), "Enter the security_group_name within 1-255 characters [a-zA-Z]+[0-9a-zA-Z_-]."),
			),
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone.",
			Required:    true,
			ForceNew:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The NAS security group description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
	}
}
