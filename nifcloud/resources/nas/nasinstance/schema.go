package nasinstance

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a NAS instance resource."

// New returns the nifcloud_nas_instance resource schema.
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
			Create:  schema.DefaultTimeout(60 * time.Minute),
			Update:  schema.DefaultTimeout(60 * time.Minute),
			Delete:  schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allocated_storage": {
			Type:        schema.TypeInt,
			Description: "The allocated storage in gibibytes.",
			Required:    true,
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The AZ for the NAS instance.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"identifier": {
			Type:             schema.TypeString,
			Description:      "The name of the NAS instance.",
			Required:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(1, 63),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The NAS instance description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"nas_security_group_name": {
			Type:        schema.TypeString,
			Description: "The security group name to associate with; which can be managed using the nifcloud_nas_security_group resource.",
			Optional:    true,
			Computed:    true,
		},
		"protocol": {
			Type:         schema.TypeString,
			Description:  "The protocol of the NAS. nfs or cifs",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"nfs", "cifs"}, true),
		},
		"master_username": {
			Type:             schema.TypeString,
			Description:      "The username for the master. (only for cifs protocol)",
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(1, 32),
		},
		"master_user_password": {
			Type:             schema.TypeString,
			Description:      "The password for the master user. (only for cifs protocol)",
			Optional:         true,
			Sensitive:        true,
			ValidateDiagFunc: validator.StringRuneCountBetween(1, 128),
		},
		"authentication_type": {
			Type:         schema.TypeInt,
			Description:  "Type of cifs authentication. (0: local auth, 1: directory service auth)",
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntInSlice([]int{0, 1}),
		},
		"directory_service_domain_name": {
			Type:             schema.TypeString,
			Description:      "The domain name of directory service.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(1, 256),
		},
		"directory_service_administrator_name": {
			Type:             schema.TypeString,
			Description:      "The user name of directory service.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(1, 128),
		},
		"directory_service_administrator_password": {
			Type:             schema.TypeString,
			Description:      "The administrator's password of directory service.",
			Optional:         true,
			Sensitive:        true,
			ValidateDiagFunc: validator.StringRuneCountBetween(1, 128),
		},
		"domain_controllers": {
			Type:        schema.TypeSet,
			Description: "The domain controller info.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"hostname": {
						Type:             schema.TypeString,
						Description:      "The hostname of domain controller.",
						Optional:         true,
						ValidateDiagFunc: validator.StringRuneCountBetween(1, 256),
					},
					"ip_address": {
						Type:        schema.TypeString,
						Description: "The IP address of domain controller.",
						Optional:    true,
					},
				},
			},
		},
		"no_root_squash": {
			Type:        schema.TypeBool,
			Description: "Turn off root squashing.",
			Optional:    true,
			Computed:    true,
		},
		"network_id": {
			Type:        schema.TypeString,
			Description: "The id of private lan.",
			Optional:    true,
		},
		"private_ip_address": {
			Type:             schema.TypeString,
			Description:      "Private IP address for NAS.",
			Optional:         true,
			Computed:         true,
			RequiredWith:     []string{"private_ip_address_subnet_mask"},
			ValidateDiagFunc: validator.IPAddress,
		},
		"private_ip_address_subnet_mask": {
			Type:         schema.TypeString,
			Description:  "The subnet mask of private IP address written in CIDR notation.",
			Optional:     true,
			RequiredWith: []string{"private_ip_address"},
		},
		"public_ip_address": {
			Type:        schema.TypeString,
			Description: "Public IP address for NAS.",
			Computed:    true,
		},
		"type": {
			Type:         schema.TypeInt,
			Description:  "The type of NAS. (0: standard type, 1: high-speed type)",
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntInSlice([]int{0, 1}),
		},
	}
}
