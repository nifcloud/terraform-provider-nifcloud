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
				validation.IntBetween(1, 8000),
				validation.IntInSlice([]int{
					100, 200, 300, 400, 500,
					600, 700, 800, 900, 1000,
					1100, 1200, 1300, 1400, 1500,
					1600, 1700, 1800, 1900, 2000,
					2100, 2200, 2300, 2400, 2500,
					2600, 2700, 2800, 2900, 3000,
					3100, 3200, 3300, 3400, 3500,
					3600, 3700, 3800, 3900, 4000,
					4100, 4200, 4300, 4400, 4500,
					4600, 4700, 4800, 4900, 5000,
					5100, 5200, 5300, 5400, 5500,
					5600, 5700, 5800, 5900, 6000,
					6100, 6200, 6300, 6400, 6500,
					6600, 6700, 6800, 6900, 7000,
					7100, 7200, 7300, 7400, 7500,
					7600, 7700, 7800, 7900, 8000,
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
