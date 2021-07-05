package volume

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a volume resource."

// New returns the nifcloud_volume resource schema.
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
		"size": {
			Type:        schema.TypeInt,
			Description: "The disk size.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.IntBetween(1, 2000),
				validation.IntInSlice([]int{
					100, 200, 300, 400, 500,
					600, 700, 800, 900, 1000,
					1100, 1200, 1300, 1400, 1500,
					1600, 1700, 1800, 1900, 2000,
				}),
			),
		},
		"volume_id": {
			Type:        schema.TypeString,
			Description: "The volume name.",
			Optional:    true,
			Computed:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 32),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), "Enter the volume_id within 1-32 characters [0-9a-zA-Z]."),
			),
		},
		"disk_type": {
			Type:        schema.TypeString,
			Description: "The disk type.",
			Optional:    true,
			ForceNew:    true,
			Default:     "High-Speed Storage A",
		},
		"instance_id": {
			Type:        schema.TypeString,
			Description: "The instance name.",
			Optional:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), "Enter the instance_id within 1-15 characters [0-9a-zA-Z]."),
			),
			ConflictsWith: []string{"instance_unique_id"},
		},
		"instance_unique_id": {
			Type:          schema.TypeString,
			Description:   "The unique ID of instance.",
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"instance_id"},
		},
		"reboot": {
			Type:         schema.TypeString,
			Description:  "The reboot type.",
			Optional:     true,
			Default:      "true",
			ValidateFunc: validation.StringInSlice([]string{"force", "true", "false"}, false),
		},
		"accounting_type": {
			Type:         schema.TypeString,
			Description:  "Accounting type. (1: monthly, 2: pay per use).",
			Optional:     true,
			Default:      "2",
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The elastic ip description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
	}
}
