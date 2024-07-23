package devopsbackuprule

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a DevOps backup rule resource."

// New returns the nifcloud_devops_backup_rule resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: createBackupRule,
		ReadContext:   readBackupRule,
		UpdateContext: updateBackupRule,
		DeleteContext: deleteBackupRule,

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
			Description: "The name of the DevOps backup rule.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 63),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z\-]+$`), "Enter a name within 1-63 alphanumeric lowercase characters and hyphens. Hyphens cannot be used at the beginning or end."),
			),
		},
		"instance_id": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps instance.",
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 30),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z\-]+$`), "Enter a name within 1-30 alphanumeric lowercase characters and hyphens. Hyphens cannot be used at the beginning or end."),
			),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the DevOps backup rule.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"backup_time": {
			Type:        schema.TypeString,
			Description: "Cron expression for backup time.",
			Computed:    true,
		},
	}
}
