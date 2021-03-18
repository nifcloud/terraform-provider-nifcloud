package dbparametergroup

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a DB parameter group resource."

// New returns the nifcloud_dhcp_config resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: deletegroup,

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
			Description: "The name of the DB parameter group.",
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z\-]+$`), ""),
			),
		},
		"family": {
			Type:        schema.TypeString,
			Description: "The DB parameter group family name.",
			Required:    true,
			ForceNew:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The description for the DB parameter group.",
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"parameter": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "The name of the parameter.",
						Required:    true,
					},
					"value": {
						Type:        schema.TypeString,
						Description: "The value of the parameter.",
						Required:    true,
					},
					"apply_method": {
						Type:         schema.TypeString,
						Description:  "Indicates when to apply parameter updates.",
						Optional:     true,
						Default:      "immediate",
						ValidateFunc: validation.StringInSlice([]string{"immediate", "pending-reboot"}, false),
					},
				},
			},
		},
	}
}
