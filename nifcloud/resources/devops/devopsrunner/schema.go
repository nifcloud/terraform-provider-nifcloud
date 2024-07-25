package devopsrunner

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a DevOps Runner resource."

// New returns the nifcloud_devops_runner resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: createRunner,
		ReadContext:   readRunner,
		UpdateContext: updateRunner,
		DeleteContext: deleteRunner,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
			Create:  schema.DefaultTimeout(80 * time.Minute),
			Update:  schema.DefaultTimeout(80 * time.Minute),
			Delete:  schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps Runner.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 30),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z\-]+$`), "Enter a name within 1-30 alphanumeric lowercase characters and hyphens. Hyphens cannot be used at the beginning or end."),
			),
		},
		"instance_type": {
			Type:        schema.TypeString,
			Description: "The instance type of the DevOps Runner.",
			Required:    true,
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone for the DevOps Runner.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"concurrent": {
			Type:        schema.TypeInt,
			Description: "Limits how many jobs can run concurrently, across all registrations.",
			Optional:    true,
			ValidateFunc: validation.All(
				validation.IntBetween(1, 50),
			),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the DevOps Runner.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"network_id": {
			Type:        schema.TypeString,
			Description: "The ID of private lan.",
			Optional:    true,
			ForceNew:    true,
		},
		"private_address": {
			Type:         schema.TypeString,
			Description:  "Private IP address for the DevOps Runner.",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
		},
		"public_ip_address": {
			Type:        schema.TypeString,
			Description: "Public IP address for the DevOps Runner.",
			Computed:    true,
		},
		"system_id": {
			Type:        schema.TypeString,
			Description: "GitLab Runner system ID.",
			Computed:    true,
		},
	}
}
