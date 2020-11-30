package instance

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a instance resource."

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
			Create:  schema.DefaultTimeout(30 * time.Minute),
			Update:  schema.DefaultTimeout(20 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"accounting_type": {
			Type:         schema.TypeString,
			Description:  "Accounting type. (1: monthly, 2: pay per use).",
			Optional:     true,
			Default:      "2",
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
		},
		"admin": {
			Type:        schema.TypeString,
			Description: "Admin user for windows os.",
			Optional:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(6, 20),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
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
			Description:      "The instance description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
		"disable_api_termination": {
			Type:        schema.TypeBool,
			Description: "If true, enables instance termination protection.",
			Optional:    true,
			Default:     false,
		},
		"image_id": {
			Type:        schema.TypeString,
			Description: "The os image identifier to use for the instance.",
			Required:    true,
			ForceNew:    true,
		},
		"instance_id": {
			Type:        schema.TypeString,
			Description: "The instance name.",
			Optional:    true,
			Computed:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"instance_type": {
			Type:        schema.TypeString,
			Description: "The type of instance to start. Updates to this field will trigger a stop/start of the instance.",
			Optional:    true,
			Default:     "mini",
		},
		"key_name": {
			Type:          schema.TypeString,
			Description:   "The key name of the Key Pair to use for the instance; which can be managed using the nifcloud_key_pair resource.",
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"admin", "password"},
		},
		"license_name": {
			Type:        schema.TypeString,
			Description: "The license name.",
			Optional:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"RDS",
				"Office(Std)",
				"Office(Pro Plus)",
			}, false),
		},
		"license_num": {
			Type:         schema.TypeInt,
			Description:  "The license count.",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 200),
		},
		"network_interface": {
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network_id": {
						Type:        schema.TypeString,
						Description: "The ID of the network to attach; 'net-COMMON_GLOBAL' or `net-COMMON_PRIVATE` or private lan network id.",
						Optional:    true,
					},
					"network_name": {
						Type:        schema.TypeString,
						Description: "The private lan name of the network to attach.",
						Optional:    true,
					},
					"ip_address": {
						Type:        schema.TypeString,
						Description: "The IP address to select from `static` or `elastic IP address` or `static IP address`; Default(null) is DHCP.",
						Optional:    true,
					},
				},
			},
		},
		"password": {
			Type:        schema.TypeString,
			Description: "Admin password for windows os.",
			Optional:    true,
			ForceNew:    true,
			Sensitive:   true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(6, 32),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), ""),
			),
		},
		"security_group": {
			Type:        schema.TypeString,
			Description: "The security group name to associate with; which can be managed using the nifcloud_security_group resource.",
			Optional:    true,
		},
		"user_data": {
			Type:        schema.TypeString,
			Description: "The user data to provide when launching the instance.",
			Optional:    true,
			ForceNew:    true,
		},
		"instance_state": {
			Type:        schema.TypeString,
			Description: "The state of the instance.",
			Computed:    true,
		},
		"public_ip": {
			Type:        schema.TypeString,
			Description: "The public ip address of instance.",
			Computed:    true,
		},
		"private_ip": {
			Type:        schema.TypeString,
			Description: "The private ip address of instance.",
			Computed:    true,
		},
		"unique_id": {
			Type:        schema.TypeString,
			Description: "The unique ID of instance.",
			Computed:    true,
		},
	}
}
