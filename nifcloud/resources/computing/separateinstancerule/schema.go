package separateinstancerule

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a separate instance rule resource."

// New returns the nifcloud_separate_instance_rule resource schema.
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
			Update:  schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The separate instance rule name.",
			Optional:    true,
			Computed:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"instance_id": {
			Type:        schema.TypeList,
			Description: "The instance name.",
			Optional:    true,
			ForceNew:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 15),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
				),
			},
			ConflictsWith: []string{"instance_unique_id"},
		},
		"instance_unique_id": {
			Type:          schema.TypeList,
			Description:   "The unique ID of instance.",
			Optional:      true,
			ForceNew:      true,
			Elem:          &schema.Schema{Type: schema.TypeString},
			ConflictsWith: []string{"instance_id"},
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The separate instance rule description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
	}
}
