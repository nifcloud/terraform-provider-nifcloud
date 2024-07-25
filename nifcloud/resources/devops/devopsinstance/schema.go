package devopsinstance

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a DevOps instance resource."

// New returns the nifcloud_devops_instance resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: createInstance,
		ReadContext:   readInstance,
		UpdateContext: updateInstance,
		DeleteContext: deleteInstance,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
			Create:  schema.DefaultTimeout(80 * time.Minute),
			Update:  schema.DefaultTimeout(80 * time.Minute),
			Delete:  schema.DefaultTimeout(30 * time.Minute),
		},
		CustomizeDiff: customdiff.ValidateChange(
			"disk_size",
			func(ctx context.Context, o, n, meta interface{}) error {
				if n.(int) < o.(int) {
					return fmt.Errorf("new disk size value must be greater than or equal to old value %d", o.(int))
				}
				return nil
			},
		),
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"instance_type": {
			Type:        schema.TypeString,
			Description: "The instance type of the DevOps instance.",
			Required:    true,
		},
		"firewall_group_name": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps firewall group to associate.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 63),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z\-]+$`), "Enter a name within 1-63 alphanumeric lowercase characters and hyphens. Hyphens cannot be used at the beginning or end."),
			),
		},
		"parameter_group_name": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps parameter group to associate.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 63),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z\-]+$`), "Enter a name within 1-63 alphanumeric lowercase characters and hyphens. Hyphens cannot be used at the beginning or end."),
			),
		},
		"disk_size": {
			Type:        schema.TypeInt,
			Description: "The allocated storage in gigabytes.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.IntBetween(100, 400),
				validation.IntDivisibleBy(100),
			),
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone for the DevOps instance.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the DevOps instance.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"initial_root_password": {
			Type:        schema.TypeString,
			Description: "Initial password for the root user.",
			Required:    true,
			Sensitive:   true,
			ForceNew:    true,
		},
		"network_id": {
			Type:        schema.TypeString,
			Description: "The ID of private lan.",
			Optional:    true,
		},
		"private_address": {
			Type:         schema.TypeString,
			Description:  "Private IP address for the DevOps instance.",
			Optional:     true,
			ValidateFunc: validation.IsCIDR,
		},
		"object_storage_account": {
			Type:        schema.TypeString,
			Description: "The account name of the object storage service.",
			Optional:    true,
			ForceNew:    true,
		},
		"object_storage_region": {
			Type:        schema.TypeString,
			Description: "The region where the bucket exists.",
			Optional:    true,
			ForceNew:    true,
		},
		"lfs_bucket_name": {
			Type:        schema.TypeString,
			Description: "The name of the bucket to put LFS objects.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"packages_bucket_name": {
			Type:        schema.TypeString,
			Description: "The name of the bucket to put packages.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"container_registry_bucket_name": {
			Type:        schema.TypeString,
			Description: "The name of the bucket to put container registry objects.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"to": {
			Type:        schema.TypeString,
			Description: "Mail address where alerts are sent.",
			Optional:    true,
		},
		"gitlab_url": {
			Type:        schema.TypeString,
			Description: "URL for GitLab.",
			Computed:    true,
		},
		"registry_url": {
			Type:        schema.TypeString,
			Description: "URL for GitLab container registry.",
			Computed:    true,
		},
		"public_ip_address": {
			Type:        schema.TypeString,
			Description: "Public IP address for the DevOps instance.",
			Computed:    true,
		},
	}
}
