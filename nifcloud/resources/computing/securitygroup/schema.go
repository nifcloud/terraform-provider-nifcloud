package securitygroup

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a security group resource."

// New returns the nifcloud_security_group resource schema.
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
		"group_name": {
			Type:        schema.TypeString,
			Description: "The name for the security group.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The security group description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone.",
			Required:    true,
			ForceNew:    true,
		},
		"log_limit": {
			Type:         schema.TypeInt,
			Description:  "The number of log data for security group.",
			Optional:     true,
			Default:      1000,
			ValidateFunc: validation.IntInSlice([]int{1000, 100000}),
		},
		"revoke_rules_on_delete": {
			Type:        schema.TypeBool,
			Description: "Instruct Terraform to revoke all of the Security Groups attached In and Out rules before deleting the rule itself. ",
			Default:     false,
			Optional:    true,
		},
	}
}
